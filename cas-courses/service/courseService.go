package service

import (
	"dbclient"
	// "encoding/json"
	"database/sql"
	// "fmt"
	// cfg "cas-dicts/config"
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
	"reflect"
	"strings"
	"time"
	// "net"
)

var currentTable string = "course"

type ICourseClient interface { 
	dbclient.IMysqlClient

	GetCourseById(id, schoolId int64) (ret interface{})
	GetAllCourses(page, items int, dtype string, schoolId int64, isVisible string) (ret []interface{})

	InsertCourse(name string, schoolId int64, desc, ctype string,
		isActive bool, start, end time.Time, event string, isVisible bool) (sql.Result, bool)
	UpdateCourse(id, uschoolId int64, name string, schoolId int64, desc, ctype string,
		isActive bool, start, end time.Time, event string, isVisible bool) (sql.Result, bool)
	DelCourseById(id int64, schoolId int64) (sql.Result, bool)
	DelCourseByIdReal(id int64, schoolId int64) (sql.Result, bool)
	DelCourses(ids []int64, schoolId int64) (sql.Result, bool)

	GetBaseInfo(username string) (int64, string, string)
}

type CourseClient struct {
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

func (client *CourseClient) GetBaseInfo(username string) (int64, string, string) {
	return baseinfo.GetSchoolInfoFromUser(client.Db, username)
}

func formatResultSet(m map[string]string) interface{} {
	ret := map[string]interface{}{}
	log.Println("query course return...", m)

	ret["id"], _ = strconv.ParseInt(m["id"], 10, 64)
	ret["name"] = m["name"]
	ret["description"] = m["desc"]
	ret["type"] = m["type"]

	ret["school_id"], _ = strconv.ParseInt(m["school_id"], 10, 64)
	if m["is_active"] == "0" {
		ret["is_active"] = false
	} else {
		ret["is_active"] = true
	}

	if m["is_visible"] == "0" {
		ret["is_visible"] = false
	} else {
		ret["is_visible"] = true
	}

	ret["start"] = m["start"]
	ret["end"] = m["end"]
	ret["events"] = m["events"]

	return ret
}

func fmtDate(tm time.Time) interface{} {
	if tm == (time.Time{}) {
		return nil
	} else {
		return tm
	}
}

func (client *CourseClient) GetCourseById(id, schoolId int64) (ret interface{}) {
	dbret := dbclient.Query(client.Db, "select * from "+currentTable+" s where s.id = ? and deleted = false and school_id = ? ", formatResultSet, id, schoolId)
	if len(dbret) >= 1 {
		// ret = ret[:1]
		ret = dbret[0]
	} else {
		ret = map[string]string{}
	}
	return
}

func (client *CourseClient) GetAllCourses(page, items int, ctype string, schoolId int64, isVisible string) (ret []interface{}) {

	// sid, sname, utype := baseinfo.GetSchoolInfoFromUser(client.Db, "admin002")
	// log.Println("school info...", sid, sname, utype)

	clauses := " deleted = false "
	conds := []interface{}{}

	if ctype != "" {
		clauses = clauses + "and type = ? "
		conds = append(conds, ctype)
	}

	visible, err := strconv.ParseBool(isVisible)
	if err == nil {
		clauses = clauses + "and is_visible = ? "
		conds = append(conds, visible)
	} else {
		log.Println(err)
	}

	clauses = clauses + "and school_id = ? "
	conds = append(conds, schoolId)

	if page == 0 {
		ret = dbclient.Query(client.Db, "select * from "+currentTable+" where "+clauses, formatResultSet, conds...)
	} else {
		offset := (page - 1) * items
		conds = append(conds, strconv.Itoa(items), strconv.Itoa(offset))
		ret = dbclient.Query(client.Db, "select * from "+currentTable+" where "+clauses+" limit ? offset ? ", formatResultSet, conds...)
	}
	return
}

func (client *CourseClient) InsertCourse(name string, schoolId int64, desc, ctype string,
	isActive bool, start, end time.Time, event string, isVisible bool) (sql.Result, bool) {

	rstart, rend := fmtDate(start), fmtDate(end)

	sql, vals := dbclient.BuildInsert(currentTable, dbclient.ParamsPairs(
		"name", name,
		"desc", desc,
		"type", ctype,
		"school_id", schoolId,
		"is_active", isActive,
		"deleted", false,
		"start", rstart,
		"end", rend,
		"events", event,
		"is_visible", isVisible,
		"create_time", time.Now(),
	),
	)

	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}
	ret := dbclient.Exec(tx, sql, vals...)
	log.Println(ret)
	tx.Commit()
	return ret, true
}

func (client *CourseClient) UpdateCourse(id, uschoolId int64, name string, schoolId int64, desc, ctype string,
	isActive bool, start, end time.Time, event string, isVisible bool) (sql.Result, bool) {

	rstart, rend := fmtDate(start), fmtDate(end)

	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}

	sql, vals := dbclient.BuildUpdate(currentTable, dbclient.ParamsPairs(
		"name", name,
		"desc", desc,
		"type", ctype,
		"school_id", schoolId,
		"is_active", isActive,
		"start", rstart,
		"end", rend,
		"events", event,
		"is_visible", isVisible,
		"update_time", time.Now(),
	), dbclient.ParamsPairs(
		"id", id,
		"school_id", schoolId,
	),
	)

	ret := dbclient.Exec(tx, sql, vals...)
	log.Println(ret)
	tx.Commit()
	return ret, true
}

func (client *CourseClient) DelCourseById(id int64, schoolId int64) (sql.Result, bool) {
	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}

	sql, vals := dbclient.BuildUpdate(currentTable, dbclient.ParamsPairs(
		"deleted", true,
	), dbclient.ParamsPairs(
		"id", id,
		"school_id", schoolId,
	),
	)

	log.Println(sql, id)
	ret := dbclient.Exec(tx, sql, vals...)
	log.Println(ret)
	tx.Commit()
	return ret, true
}

func (client *CourseClient) DelCourseByIdReal(id int64, schoolId int64) (sql.Result, bool) {
	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}

	sql, vals := dbclient.BuildDelete(currentTable, dbclient.ParamsPairs(
		"id", id,
		"school_id", schoolId,
	),
	)

	ret := dbclient.Exec(tx, sql, vals...)
	log.Println(ret)
	tx.Commit()
	return ret, true
}

func (client *CourseClient) DelCourses(ids []int64, schoolId int64) (sql.Result, bool) {
	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}

	ids2str := make([]string, len(ids))
	for i, v := range ids {
		ids2str[i] = strconv.FormatInt(v, 10)
	}

	sql, vals := dbclient.BuildUpdateWithOpts(currentTable, dbclient.ParamsPairs(
		"is_visible", false,
	), dbclient.ParamsPairs(
		"school_id", schoolId,
	), nil,
		"id in "+"("+strings.Join(ids2str, ",")+")",
	)

	ret := dbclient.Exec(tx, sql, vals...)
	log.Println(ret)
	tx.Commit()
	return ret, true
}
