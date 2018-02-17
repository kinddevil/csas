package service

import (
	"dbclient"
	// "encoding/json"
	"database/sql"
	"fmt"
	// cfg "cas-dicts/config"
	// "github.com/gorilla/mux"
	// "io"
	"baseinfo"
	// "io/ioutil"
	// "errors"
	"log"
	"sync"
	// "net/http"
	// "os"
	"strconv"
	// "html/template"
	// "github.com/satori/go.uuid"
	// return uuid.NewV4().String()
	"reflect"
	"strings"
	"time"
	// "net"
)

var currentTable string = "privilege"

type IRolesClient interface {
	dbclient.IMysqlClient

	GetRoleById(id int64, username string) (ret interface{})
	GetAllRoles(page, items int, username string) (ret []interface{})
	InsertRole(pname, username, permissions string, usernames []string) (sql.Result, bool)
	UpdateRole(id int64, pname, username, permissions string, usernames []string) (sql.Result, bool)
	DelRoles(ids []int64) (sql.Result, bool)

	GetBaseInfo(username string) (int64, string, string)
}

type RolesClient struct {
	dbclient.MysqlClient
}

// Reflect all fields to map
func GetFieldMap(obj interface{}) (ret map[string]string) {
	val := reflect.ValueOf(obj).Elem()
	ret = make(map[string]string)
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		key := strings.ToLower(typeField.Name)
		if typeField.PkgPath != "" {
			// Private method
			continue
		} else {
			ret[key] = typeField.Name
		}
	}
	return
}

func (client *RolesClient) GetBaseInfo(username string) (int64, string, string) {
	return baseinfo.GetSchoolInfoFromUser(client.Db, username)
}

func formatResultSet(m map[string]string) interface{} {
	ret := map[string]interface{}{}
	log.Println("query dict return...", m)

	ret["id"], _ = strconv.ParseInt(m["id"], 10, 64)
	ret["name"] = m["name"]
	ret["description"] = m["description"]
	ret["type"] = m["type"]
	ret["permission_ids"] = m["permission_ids"]
	if m["school_id"] == "0" {
		ret["school_id"] = ""
	} else {
		ret["school_id"] = m["school_id"]
	}

	return ret
}

func getUsersByRoles(client *RolesClient, roleIds []int64) (map[int64][]map[string]string, error) {
	roleIdsStr := func(roleIds []int64) []string {
		ret := make([]string, len(roleIds))
		for i, val := range roleIds {
			ret[i] = fmt.Sprintf("%v", val)
		}
		return ret
	}(roleIds)
	inQuery := "(" + strings.Join(roleIdsStr, ",") + ")"
	sqlStr := `select up.privilegeid, u.username, u.name from user_privilege up
						left join user u on up.username = u.username
							where up.privilegeid in ` + inQuery
	ret := make(map[int64][]map[string]string)
	var pid int
	var uname sql.NullString
	var name sql.NullString
	rows, err := client.Db.Query(sqlStr)
	defer rows.Close()

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&pid, &uname, &name)
		if err != nil {
			return nil, err
		}

		if !uname.Valid {
			continue
		}

		pid64 := int64(pid)
		if _, ok := ret[pid64]; !ok {
			ret[pid64] = []map[string]string{{"id": uname.String, "name": name.String}}
		} else {
			ret[pid64] = append(ret[pid64], map[string]string{"id": uname.String, "name": name.String})
		}
	}

	log.Println("ret users...", ret)

	return ret, nil

}

func (client *RolesClient) GetRoleById(id int64, username string) (ret interface{}) {
	sid, _, stype := baseinfo.GetSchoolInfoFromUser(client.Db, username)
	if stype == "" {
		log.Println("there is no user or invalid user for given token")
		return nil
	}

	query := " p where p.id = ? and is_deleted = false"
	params := []interface{}{id}

	if sid != 0 {
		query = query + " and school_id = ? "
		params = append(params, sid)
	}

	dbret := dbclient.Query(client.Db, "select * from "+currentTable+query, formatResultSet, params...)
	if len(dbret) >= 1 {
		// ret = ret[:1]
		ret = dbret[0]

		pid := ret.(map[string]interface{})["id"].(int64)
		users, err := getUsersByRoles(client, []int64{pid})
		if err != nil {
			log.Println("find user by privilege error", err)
		}

		ret.(map[string]interface{})["user_ids"] = users[pid]
	} else {
		ret = map[string]string{}
	}
	return
}

func (client *RolesClient) GetAllRoles(page, items int, username string) (ret []interface{}) {

	sid, _, stype := baseinfo.GetSchoolInfoFromUser(client.Db, username)
	if stype == "" {
		log.Println("there is no user or invalid user for given token")
		return nil
	}

	clauses := " is_deleted = false "
	conds := []interface{}{}

	if sid != 0 {
		clauses = clauses + " and school_id = ? "
		conds = append(conds, sid)
	}

	if page == 0 {
		ret = dbclient.Query(client.Db, "select * from "+currentTable+" where "+clauses, formatResultSet, conds...)
	} else {
		offset := (page - 1) * items
		conds = append(conds, strconv.Itoa(items), strconv.Itoa(offset))
		ret = dbclient.Query(client.Db, "select * from "+currentTable+" where "+clauses+" limit ? offset ? ", formatResultSet, conds...)
	}

	var wg sync.WaitGroup
	for i, role := range ret {
		pid := role.(map[string]interface{})["id"].(int64)
		wg.Add(1)
		go func() {
			users, err := getUsersByRoles(client, []int64{role.(map[string]interface{})["id"].(int64)})
			if err != nil {
				log.Println("find user by privilege error", err)
			}
			ret[i].(map[string]interface{})["user_ids"] = users[pid]
			wg.Done()
		}()
		wg.Wait()
	}

	return
}

func (client *RolesClient) InsertRole(pname, username, permissions string, usernames []string) (sql.Result, bool) {
	sid, sname, stype := baseinfo.GetSchoolInfoFromUser(client.Db, username)

	if stype == "" {
		log.Println("there is no user or invalid user with no type")
		return nil, false
	}

	log.Println("baseinfo...", sid, sname, stype)

	sql, vals := dbclient.BuildInsert(currentTable, dbclient.ParamsPairs(
		"name", pname,
		"school_id", sid,
		"permission_ids", permissions,
		"is_deleted", false,
		"create_time", time.Now(),
	),
	)

	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}
	ret := dbclient.Exec(tx, sql, vals...)
	pid, ierr := ret.LastInsertId()
	if ierr != nil {
		log.Println(ierr)
		tx.Rollback()
		return ret, false
	}

	if len(usernames) > 0 && usernames[0] != "" {
		userChecks := make([]string, len(usernames))
		for i, val := range usernames {
			userChecks[i] = "'" + val + "'"
		}
		userCheckCond := " and username in " + "(" + strings.Join(userChecks, ",") + ")"
		if sid == 0 {
			// if it is super admin, must be admin or school-admin
			// userCheckCond = userCheckCond + " and school_id is not null "
			userCheckCond = userCheckCond + " and type not like '%admin' "
		} else {
			// if it is not super admin, school id must be the same
			userCheckCond = userCheckCond + fmt.Sprintf(" and school_id!=%v ", sid)
		}

		// Check school...
		userChcSql := "select * from user where 1=1 "
		fmt.Println(userChcSql+userCheckCond, "check sql...")
		count := dbclient.QueryWithTran(tx, userChcSql+userCheckCond, nil) // TODO: add error handler
		if len(count) > 0 {
			log.Println("contians other school user...")
			tx.Rollback()
			return nil, false
		}
	}

	for _, uname := range usernames {
		sqlUserPrivilege, vals := dbclient.BuildInsert("user_privilege", dbclient.ParamsPairs(
			"username", uname,
			"privilegeid", pid,
			"create_time", time.Now(),
		),
		)
		dbclient.Exec(tx, sqlUserPrivilege, vals...) // TODO: add error handler
	}

	log.Println(ret)
	tx.Commit()
	return ret, true
}

func (client *RolesClient) UpdateRole(id int64, pname, username, permissions string, usernames []string) (sql.Result, bool) {
	sid, _, stype := baseinfo.GetSchoolInfoFromUser(client.Db, username)

	log.Println(sid, stype, "get school...")
	if stype == "" {
		log.Println("there is no user or invalid user with no type")
		return nil, false
	}

	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}

	sql, vals := dbclient.BuildUpdate(currentTable, dbclient.ParamsPairs(
		"name", pname,
		"school_id", sid,
		"permission_ids", permissions,
	), dbclient.ParamsPairs(
		"id", id,
	),
	)

	ret := dbclient.Exec(tx, sql, vals...)
	log.Println(ret)

	// Check school...
	if len(usernames) > 0 && usernames[0] != "" {
		userChecks := make([]string, len(usernames))
		for i, val := range usernames {
			userChecks[i] = "'" + val + "'"
		}
		userCheckCond := " and username in " + "(" + strings.Join(userChecks, ",") + ")"
		if sid == 0 {
			// if it is super admin, must be admin or school-admin
			// userCheckCond = userCheckCond + " and school_id is not null "
			userCheckCond = userCheckCond + " and type not like '%admin' "
		} else {
			// if it is not super admin, school id must be the same
			userCheckCond = userCheckCond + fmt.Sprintf(" and school_id!=%v ", sid)
		}

		// Check school...
		userChcSql := "select * from user where 1=1 "
		fmt.Println(userChcSql+userCheckCond, "check sql...")
		count := dbclient.QueryWithTran(tx, userChcSql+userCheckCond, nil) // TODO: add error handler
		if len(count) > 0 {
			log.Println("contians other school user...")
			tx.Rollback()
			return nil, false
		}
	}

	sqlDel, vals := dbclient.BuildDelete("user_privilege", dbclient.ParamsPairs(
		"privilegeid", id,
	),
	)
	log.Println(sqlDel, "delsql...")
	dbclient.Exec(tx, sqlDel, vals...) // TODO: add error handler

	for _, uname := range usernames {
		sqlUserPrivilege, vals := dbclient.BuildInsert("user_privilege", dbclient.ParamsPairs(
			"username", uname,
			"privilegeid", id,
			"create_time", time.Now(),
		),
		)
		dbclient.Exec(tx, sqlUserPrivilege, vals...) // TODO: add error handler
	}

	tx.Commit()
	return ret, true
}

func (client *RolesClient) DelRoles(ids []int64) (sql.Result, bool) {
	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}

	ids2str := make([]string, len(ids))
	for i, v := range ids {
		ids2str[i] = strconv.FormatInt(v, 10)
	}

	sql, vals := dbclient.BuildUpdateWithOpts(currentTable, dbclient.ParamsPairs(
		"is_deleted", true,
	), nil, nil,
		"id in "+"("+strings.Join(ids2str, ",")+")",
	)

	ret := dbclient.Exec(tx, sql, vals...)
	log.Println(ret)
	tx.Commit()
	return ret, true
}
