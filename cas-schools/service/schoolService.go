package service

import (
	"dbclient"
	// "encoding/json"
	"database/sql"
	// "fmt"
	// cfg "schools/config"
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

var tableSchool string = "schools"

type ISchoolsClient interface {
	dbclient.IMysqlClient

	GetSchoolById(id int64) (ret []interface{})
	GetAllSchools(page, items int) (ret []interface{})
	InsertSchool(name, contact, phone, province, city, county, provinceCode, cityCode, countyCode, addr, fax, email, web, post, from, to, contractId, contract string, isPayment, isLock bool, teacherNo, studentNo int) (sql.Result, bool)
	UpdateSchool(id int64, name, contact, phone, province, city, county, provinceCode, cityCode, countyCode, addr, fax, email, web, post, from, to, contractId, contract string, isPayment, isLock bool, teacherNo, studentNo int) (sql.Result, bool)
	DelSchoolById(id int64) (sql.Result, bool)

	GetBaseInfo(username string) (int64, string, string)
}

type SchoolsClient struct {
	dbclient.MysqlClient
}

type UserInfo struct {
	Ids   int
	Names string
	pri   string
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

func (client *SchoolsClient) GetBaseInfo(username string) (int64, string, string) {
	return baseinfo.GetSchoolInfoFromUser(client.Db, username)
}

func formatSchoolsResultSet(m map[string]string) interface{} {
	ret := map[string]interface{}{}
	log.Println("query school return...", m)

	ret["id"] = m["id"]
	ret["name"] = m["name"]
	ret["contact"] = m["contact"]
	ret["phone"] = m["phone"]
	ret["province"] = m["province"]
	ret["city"] = m["city"]
	ret["county"] = m["county"]
	ret["province_code"] = m["province_code"]
	ret["city_code"] = m["city_code"]
	ret["address"] = m["address"]
	ret["fax"] = m["fax"]
	ret["email"] = m["email"]
	ret["web"] = m["web"]
	ret["post"] = m["post"]
	ret["from"] = m["start_time"]
	ret["to"] = m["expire_time"]
	ret["is_payment"], _ = strconv.ParseBool(m["is_payment"])
	ret["teacher"], _ = strconv.Atoi(m["teacher_no"])
	ret["student"], _ = strconv.Atoi(m["student_no"])
	ret["contract_id"] = m["contract_id"]
	ret["contract"] = m["contract"]
	ret["is_lock"], _ = strconv.ParseBool(m["is_lock"])

	return ret
}

func (client *SchoolsClient) GetSchoolById(id int64) (ret []interface{}) {
	ret = dbclient.Query(client.Db, "select * from "+tableSchool+" s where s.id = ? and deleted = false", formatSchoolsResultSet, id)
	if len(ret) > 1 {
		ret = ret[:1]
	}
	return
}

func (client *SchoolsClient) GetAllSchools(page, items int) (ret []interface{}) {

	// sid, sname, utype := baseinfo.GetSchoolInfoFromUser(client.Db, "admin002")
	// log.Println("school info...", sid, sname, utype)

	if page == 0 {
		ret = dbclient.Query(client.Db, "select * from "+tableSchool+" where deleted = false", formatSchoolsResultSet)
	} else {
		offset := page * items
		ret = dbclient.Query(client.Db, "select * from "+tableSchool+" where deleted = false limit ? offset ? ", formatSchoolsResultSet, strconv.Itoa(items), strconv.Itoa(offset))
	}
	return
}

func (client *SchoolsClient) InsertSchool(name, contact, phone, province, city, county, provinceCode, cityCode, countyCode, addr, fax, email, web, post, from, to, contractId, contract string, isPayment, isLock bool, teacherNo, studentNo int) (sql.Result, bool) {
	sql, vals := dbclient.BuildInsert(tableSchool, dbclient.ParamsPairs(
		"name", name,
		"contact", contact,
		"phone", phone,
		"province", province,
		"city", city,
		"county", county,
		"province_code", provinceCode,
		"city_code", cityCode,
		"county_code", countyCode,
		"address", addr,
		"fax", fax,
		"email", email,
		"web", web,
		"post", post,
		"start_time", from,
		"expire_time", to,
		"contract_id", contractId,
		"contract", contract,
		"is_payment", isPayment,
		"is_lock", isLock,
		"teacher_no", teacherNo,
		"student_no", studentNo,

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

func (client *SchoolsClient) UpdateSchool(id int64, name, contact, phone, province, city, county, provinceCode, cityCode, countyCode, addr, fax, email, web, post, from, to, contractId, contract string, isPayment, isLock bool, teacherNo, studentNo int) (sql.Result, bool) {
	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}

	sql, vals := dbclient.BuildUpdate(tableSchool, dbclient.ParamsPairs(
		"name", name,
		"contact", contact,
		"phone", phone,
		"province", province,
		"city", city,
		"county", county,
		"province_code", provinceCode,
		"city_code", cityCode,
		"county_code", countyCode,
		"address", addr,
		"fax", fax,
		"email", email,
		"web", web,
		"post", post,
		"start_time", from,
		"expire_time", to,
		"contract_id", contractId,
		"contract", contract,
		"is_payment", isPayment,
		"is_lock", isLock,
		"teacher_no", teacherNo,
		"student_no", studentNo,
	), dbclient.ParamsPairs(
		"id", id,
	),
	)

	ret := dbclient.Exec(tx, sql, vals...)
	log.Println(ret)
	tx.Commit()
	return ret, true
}

func (client *SchoolsClient) DelSchoolById(id int64) (sql.Result, bool) {
	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}

	sql, vals := dbclient.BuildUpdate(tableSchool, dbclient.ParamsPairs(
		"deleted", true,
	), dbclient.ParamsPairs(
		"id", id,
	),
	)

	ret := dbclient.Exec(tx, sql, vals...)
	log.Println(ret)
	tx.Commit()
	return ret, true
}
