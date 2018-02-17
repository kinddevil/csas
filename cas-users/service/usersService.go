package service

import (
	"dbclient"
	// "encoding/json"
	"database/sql"
	"fmt"
	// cfg "cas-users/config"
	// "github.com/gorilla/mux"
	// "io"
	"baseinfo"
	// "io/ioutil"
	"log"
	// "net/http"
	// "os"
	"strconv"
	// "html/template"
	// "github.com/satori/go.uuid"
	// return uuid.NewV4().String()
	"errors"
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"strings"
	"time"
	// "net"
)

var currentTable string = "user"

type IUsersClient interface {
	dbclient.IMysqlClient

	GetUserByUsername(username string) (ret interface{})
	GetAllUsers(page, items int, schoolId int64, uType, targetType string) (ret []interface{})
	InsertUser(username, name, email, roleIds,
		password, schoolName, uType, phone string, schoolId int64, activated, isLock bool) (sql.Result, bool)
	UpdateUser(username, name, email, roleIds,
		schoolName, uType, phone string, schoolId int64, activated, isLock bool) (sql.Result, bool)
	UpdatePwd(username, password string) (sql.Result, bool)
	CheckPwd(username, pwd string) error
	DelUserById(username string) (sql.Result, bool)
	DelUsers(usernames []string, schoolId int64, uType string) (sql.Result, bool)

	GetBaseInfo(username string) (int64, string, string)
}

type UsersClient struct {
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

func (client *UsersClient) GetBaseInfo(username string) (int64, string, string) {
	return baseinfo.GetSchoolInfoFromUser(client.Db, username)
}

func formatResultSet(m map[string]string) interface{} {
	ret := map[string]interface{}{}

	ret["id"] = m["username"]
	ret["username"] = m["username"]
	ret["name"] = m["name"]
	ret["email"] = m["email"]
	ret["activated"] = m["activated"]
	ret["activationkey"] = m["activationkey"]
	ret["resetpasswordkey"] = m["resetpasswordkey"]

	ret["school_id"] = m["school_id"]
	ret["school_name"] = m["school_name"]
	ret["type"] = m["type"]
	ret["phone"] = m["phone"]
	ret["last_load"] = m["last_load"]

	ret["modify_time"] = m["modify_time"]
	if m["is_lock"] == "1" {
		ret["is_lock"] = true
	} else {
		ret["is_lock"] = false
	}

	ret["load_counter"] = m["load_counter"]
	ret["pass_version"] = m["pass_version"]

	ret["prname"] = m["pname"]
	ret["prid"] = m["pid"]
	return ret
}

func formatUserWithRole(userroles []interface{}) []interface{} {
	users := map[string]interface{}{}
	ret := make([]interface{}, 0, 10)
	for _, v := range userroles {
		userrole := v.(map[string]interface{})

		if _, ok := users[userrole["username"].(string)]; !ok {
			users[userrole["username"].(string)] = v
			users[userrole["username"].(string)].(map[string]interface{})["role"] = make([]map[string]interface{}, 0, 5)
		}

		log.Println("users...", users, userrole["prname"], userrole["prid"])
		if userrole["prname"] != "" && userrole["prname"] != nil {
			privilegeId, _ := strconv.Atoi(userrole["prid"].(string))
			role := map[string]interface{}{"id": privilegeId, "name": userrole["prname"]}
			users[userrole["username"].(string)].(map[string]interface{})["role"] = append(users[userrole["username"].(string)].(map[string]interface{})["role"].([]map[string]interface{}), role)
		}
	}
	for _, v := range users {
		ret = append(ret, v)
	}
	return ret
}

func (client *UsersClient) GetUserByUsername(username string) (ret interface{}) {
	query := fmt.Sprintf(`select u.*, p.id as pid, p.name as pname from %v u left join user_privilege up on u.username = up.username
			left join privilege p on up.privilegeid = p.id where u.username = ? and u.is_deleted = false order by u.username
				`, currentTable)
	dbret := dbclient.Query(client.Db, query, formatResultSet, username)
	dataret := formatUserWithRole(dbret)
	if len(dataret) >= 1 {
		// ret = ret[:1]
		ret = dataret[0]
	} else {
		ret = map[string]interface{}{}
	}
	return
}

func (client *UsersClient) GetAllUsers(page, items int, schoolId int64, uType, targetType string) (ret []interface{}) {

	// sid, sname, utype := baseinfo.GetSchoolInfoFromUser(client.Db, "admin002")
	// log.Println("school info...", sid, sname, utype)

	query := fmt.Sprintf(`select u.*, p.id as pid, p.name as pname from %v u left join user_privilege up on u.username = up.username
			left join privilege p on up.privilegeid = p.id
				`, currentTable)

	clauses := " u.is_deleted = false "
	conds := []interface{}{}

	if targetType != "" {
		clauses = clauses + "and u.type = ? "
		conds = append(conds, targetType)
	}

	if uType == "admin" {
		clauses = clauses + "and u.type like ? "
		conds = append(conds, "%admin")
	} else {
		clauses = clauses + "and u.school_id = ? "
		conds = append(conds, schoolId)
	}

	if page == 0 {
		ret = dbclient.Query(client.Db, query+" where "+clauses, formatResultSet, conds...)
	} else {
		offset := (page - 1) * items
		conds = append(conds, strconv.Itoa(items), strconv.Itoa(offset))
		ret = dbclient.Query(client.Db, query+" where "+clauses+" limit ? offset ? ", formatResultSet, conds...)
	}
	ret = formatUserWithRole(ret)
	return
}

func (client *UsersClient) InsertUser(username, name, email, roleIds,
	password, schoolName, uType, phone string, schoolId int64, activated, isLock bool) (sql.Result, bool) {
	sql, vals := dbclient.BuildInsert(currentTable, dbclient.ParamsPairs(
		"username", username,
		"name", name,
		"email", email,
		"is_deleted", false,

		"password", generatePwd(password),
		"school_id", schoolId,
		"school_name", schoolName,
		"type", uType,
		"phone", phone,
		"activated", !isLock,
		"is_lock", isLock,

		"modify_time", time.Now(),
		"create_time", time.Now(),
	),
	)

	delQuery := "delete from user_privilege where username = ?"
	insQuery := "insert into user_privilege(username, privilegeid) values(?, ?)"

	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}
	ret := dbclient.Exec(tx, sql, vals...)
	log.Println(ret)

	dbclient.Exec(tx, delQuery, username)
	for _, v := range strings.Split(roleIds, ",") {
		pid, _ := strconv.Atoi(v)
		dbclient.Exec(tx, insQuery, username, pid)
	}

	tx.Commit()
	return ret, true
}

func (client *UsersClient) UpdateUser(username, name, email, roleIds,
	schoolName, uType, phone string, schoolId int64, activated, isLock bool) (sql.Result, bool) {
	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}

	sql, vals := dbclient.BuildUpdate(currentTable, dbclient.ParamsPairs(
		"name", name,
		"email", email,

		"school_id", schoolId,
		"school_name", schoolName,
		"type", uType,
		"phone", phone,
		"activated", !isLock,
		"is_lock", isLock,

		"modify_time", time.Now(),
		"update_time", time.Now(),
	), dbclient.ParamsPairs(
		"username", username,
	),
	)

	delQuery := "delete from user_privilege where username = ?"
	insQuery := "insert into user_privilege(username, privilegeid) values(?, ?)"

	ret := dbclient.Exec(tx, sql, vals...)
	log.Println(ret)

	dbclient.Exec(tx, delQuery, username)
	for _, v := range strings.Split(roleIds, ",") {
		pid, _ := strconv.Atoi(v)
		dbclient.Exec(tx, insQuery, username, pid)
	}

	tx.Commit()
	return ret, true
}

func (client *UsersClient) UpdatePwd(username, pwd string) (sql.Result, bool) {
	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}

	sql, vals := dbclient.BuildUpdate(currentTable, dbclient.ParamsPairs(
		"password", generatePwd(pwd),

		"modify_time", time.Now(),
		"update_time", time.Now(),
	), dbclient.ParamsPairs(
		"username", username,
	),
	)

	ret := dbclient.Exec(tx, sql, vals...)
	log.Println(ret)
	tx.Commit()
	return ret, true
}

func (client *UsersClient) CheckPwd(username, pwd string) error {
	query := fmt.Sprintf(`select password from %v where username = ?
				`, currentTable)
	upwds := dbclient.Query(client.Db, query, nil, username)
	if len(upwds) < 1 || upwds[0] == nil {
		return errors.New("Cannot find available user...")
	}
	dbPwd := upwds[0].(map[string]string)["password"]
	if err := bcrypt.CompareHashAndPassword([]byte(dbPwd), []byte(pwd)); err != nil {
		return errors.New("password not match")
	}
	return nil
}

func (client *UsersClient) DelUserById(username string) (sql.Result, bool) {
	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}

	sql, vals := dbclient.BuildUpdate(currentTable, dbclient.ParamsPairs(
		"is_deleted", true,
		"activated", false,
	), dbclient.ParamsPairs(
		"username", username,
	),
	)

	ret := dbclient.Exec(tx, sql, vals...)
	log.Println(ret)
	tx.Commit()
	return ret, true
}

func (client *UsersClient) DelUserByUsernameReal(username string) (sql.Result, bool) {
	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}

	sql, vals := dbclient.BuildDelete(currentTable, dbclient.ParamsPairs(
		"username", username,
	),
	)

	ret := dbclient.Exec(tx, sql, vals...)
	log.Println(ret)
	tx.Commit()
	return ret, true
}

func (client *UsersClient) DelUsers(usernames []string, schoolId int64, uType string) (sql.Result, bool) {
	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}

	cond := ""
	if uType == "admin" {
		cond = " and type like '%admin' "
	} else {
		cond = fmt.Sprintf(" and school_id = %v ", schoolId)
	}

	pUsernames := make([]string, len(usernames))
	for i, v := range usernames {
		pUsernames[i] = "'" + v + "'"
	}

	sql, vals := dbclient.BuildUpdateWithOpts(currentTable, dbclient.ParamsPairs(
		"is_deleted", true,
		"activated", false,
	), nil, nil,
		"username in "+"("+strings.Join(pUsernames, ",")+")"+cond,
	)

	log.Println("delete sql...", sql, vals)
	ret := dbclient.Exec(tx, sql, vals...)
	log.Println(ret)
	tx.Commit()
	return ret, true
}

func generatePwd(pass string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}
