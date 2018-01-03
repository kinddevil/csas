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

var currentTable string = "dict"

type IDictsClient interface {
	dbclient.IMysqlClient

	GetDictById(id int64) (ret interface{})
	GetAllDicts(page, items int) (ret []interface{})
	InsertDict(name, desc, dtype string) (sql.Result, bool)
	UpdateDict(id int64, name, desc, dtype string) (sql.Result, bool)
	DelDictById(id int64) (sql.Result, bool)
	DelDictByIdReal(id int64) (sql.Result, bool)

	GetBaseInfo(username string) (int64, string, string)
}

type DictsClient struct {
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

func (client *DictsClient) GetBaseInfo(username string) (int64, string, string) {
	return baseinfo.GetSchoolInfoFromUser(client.Db, username)
}

func formatResultSet(m map[string]string) interface{} {
	ret := map[string]interface{}{}
	log.Println("query dict return...", m)

	ret["id"] = m["id"]
	ret["name"] = m["key"]
	ret["description"] = m["desc"]
	ret["type"] = m["type"]
	ret["value"] = m["value"]
	ret["status"] = m["status"]

	return ret
}

func (client *DictsClient) GetDictById(id int64) (ret interface{}) {
	dbret := dbclient.Query(client.Db, "select * from "+currentTable+" s where s.id = ? and is_deleted = false", formatResultSet, id)
	if len(dbret) >= 1 {
		// ret = ret[:1]
		ret = dbret[0]
	} else {
		ret = map[string]string{}
	}
	return
}

func (client *DictsClient) GetAllDicts(page, items int) (ret []interface{}) {

	// sid, sname, utype := baseinfo.GetSchoolInfoFromUser(client.Db, "admin002")
	// log.Println("school info...", sid, sname, utype)

	if page == 0 {
		ret = dbclient.Query(client.Db, "select * from "+currentTable+" where is_deleted = false", formatResultSet)
	} else {
		offset := page * items
		ret = dbclient.Query(client.Db, "select * from "+currentTable+" where is_deleted = false limit ? offset ? ", formatResultSet, strconv.Itoa(items), strconv.Itoa(offset))
	}
	return
}

func (client *DictsClient) InsertDict(name, desc, dtype string) (sql.Result, bool) {
	sql, vals := dbclient.BuildInsert(currentTable, dbclient.ParamsPairs(
		"key", name,
		"desc", desc,
		"type", dtype,
		"is_deleted", false,
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

func (client *DictsClient) UpdateDict(id int64, name, desc, dtype string) (sql.Result, bool) {
	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}

	sql, vals := dbclient.BuildUpdate(currentTable, dbclient.ParamsPairs(
		"key", name,
		"desc", desc,
		"type", dtype,
	), dbclient.ParamsPairs(
		"id", id,
	),
	)

	ret := dbclient.Exec(tx, sql, vals...)
	log.Println(ret)
	tx.Commit()
	return ret, true
}

func (client *DictsClient) DelDictById(id int64) (sql.Result, bool) {
	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}

	sql, vals := dbclient.BuildUpdate(currentTable, dbclient.ParamsPairs(
		"is_deleted", true,
	), dbclient.ParamsPairs(
		"id", id,
	),
	)

	ret := dbclient.Exec(tx, sql, vals...)
	log.Println(ret)
	tx.Commit()
	return ret, true
}

func (client *DictsClient) DelDictByIdReal(id int64) (sql.Result, bool) {
	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}

	sql, vals := dbclient.BuildDelete(currentTable, dbclient.ParamsPairs(
		"id", id,
	),
	)

	ret := dbclient.Exec(tx, sql, vals...)
	log.Println(ret)
	tx.Commit()
	return ret, true
}
