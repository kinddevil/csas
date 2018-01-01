package service

import (
	"dbclient"
	// "encoding/json"
	cfg "ads/config"
	"database/sql"
	"fmt"
	// "github.com/gorilla/mux"
	// "io"
	"baseinfo"
	// "io/ioutil"
	"log"
	// "net/http"
	// "os"
	"strconv"
	// "html/template"
	"github.com/satori/go.uuid"
	// return uuid.NewV4().String()
	"reflect"
	"strings"
	// "time"
	// "net"
)

var tableAd string = "advertising"

type IAdsClient interface {
	dbclient.IMysqlClient
	InsertAd(title, province, city, startTime, expireTime, schoolIds string, isAnonymous, isSchool, isTeacher, isStudent bool) (sql.Result, bool)
	UpdateAd(id, pending int, title, province, city, startTime, expireTime, schoolIds string, isAnonymous, isSchool, isTeacher, isStudent bool) (sql.Result, bool)
	GetAdById(id int64) interface{}
	GetAllAds(page, items int) []interface{}
	InsertEmptyAd() (sql.Result, bool)
	SaveUploadFiles(adId int, filename string) (sql.Result, bool)
	GetBaseInfo(username string) (int64, string, string)

	SetupDb()
}

type AdsClient struct {
	dbclient.MysqlClient
}

type UserInfo struct {
	Ids   int
	Names string
	pri   string
}

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

func (client *AdsClient) Query(sqlStr string, args ...interface{}) {
	tx, err := client.Db.Begin()
	stmtOut, err := tx.Prepare(sqlStr)
	defer stmtOut.Close()
	rows, err := stmtOut.Query(args...)
	columns, err := rows.Columns()
	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			fmt.Println(columns[i], ": ", value)
		}
		fmt.Println("-----------------------------------")
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	tx.Commit()
}

// isAnonymous for login
func (client *AdsClient) InsertAd(title, province, city, startTime, expireTime, schoolIds string, isAnonymous, isSchool, isTeacher, isStudent bool) (sql.Result, bool) {
	// bitwise: 0: Anonymous 1 School 2 teacher 3 student
	log.Println("status...", isAnonymous, isSchool, isTeacher, isStudent)
	displayPages := 0
	if isAnonymous {
		displayPages = displayPages | 1
	}
	if isSchool {
		displayPages = displayPages | 1<<1
	}
	if isTeacher {
		displayPages = displayPages | 1<<2
	}
	if isStudent {
		displayPages = displayPages | 1<<3
	}

	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}
	log.Println("time...", startTime, expireTime)
	sql, vals := dbclient.BuildInsert(tableAd, dbclient.ParamsPairs(
		"title", title,
		"province", province,
		"city", city,
		"start_time", startTime,
		"expire_time", expireTime,
		"school_ids", schoolIds,
		"display_pages", displayPages,
	),
	)
	ret := dbclient.Exec(tx, sql, vals...)
	fmt.Println(ret)
	tx.Commit()
	return ret, true
}

func (client *AdsClient) InsertEmptyAd() (sql.Result, bool) {
	sql, vals := dbclient.BuildInsert(tableAd, dbclient.ParamsPairs(
		"title", "",
	),
	)
	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}
	ret := dbclient.Exec(tx, sql, vals...)
	fmt.Println(ret)
	tx.Commit()
	return ret, true
}

func (client *AdsClient) SaveUploadFiles(adId int, filename string) (sql.Result, bool) {
	sql, vals := dbclient.BuildInsert("ads_files", dbclient.ParamsPairs(
		"id", uuid.NewV4().String(),
		"advertising_id", adId,
		"name", filename,
	),
	)
	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}
	ret := dbclient.Exec(tx, sql, vals...)
	fmt.Println(ret)
	tx.Commit()
	return ret, true
}

func (client *AdsClient) UpdateAd(id, pending int, title, province, city, startTime, expireTime, schoolIds string, isAnonymous, isSchool, isTeacher, isStudent bool) (sql.Result, bool) {
	// bitwise: 0: Anonymous 1 School 2 teacher 3 student
	displayPages := 0
	if isAnonymous {
		displayPages = displayPages | 1
	}
	if isSchool {
		displayPages = displayPages | (1 << 1)
	}
	if isTeacher {
		displayPages = displayPages | (1 << 2)
	}
	if isStudent {
		displayPages = displayPages | (1 << 3)
	}

	log.Println("display pages...", displayPages, isAnonymous, isSchool, isTeacher, isStudent)

	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}
	sql, vals := dbclient.BuildUpdate(tableAd, dbclient.ParamsPairs(
		"title", title,
		"province", province,
		"city", city,
		"start_time", startTime,
		"expire_time", expireTime,
		"school_ids", schoolIds,
		"display_pages", displayPages,
		"pending", pending,
	), dbclient.ParamsPairs(
		"id", id,
	),
	)
	ret := dbclient.Exec(tx, sql, vals...)
	fmt.Println(ret)
	tx.Commit()
	return ret, true
}

func formatAdsResultSet(m map[string]string) interface{} {
	ret := map[string]interface{}{}
	log.Println("query ad return...", m)

	ret["id"] = m["id"]
	ret["title"] = m["title"]
	ret["from"] = m["start_time"]
	ret["to"] = m["expire_time"]
	ret["school_ids"] = m["school_ids"]

	ret["on_login_page"] = false
	ret["on_school_page"] = false
	ret["on_teacher_page"] = false
	ret["on_student_page"] = false

	display, _ := strconv.Atoi(m["display_pages"])
	log.Println("display pages...", display)
	if display&1 == 1 {
		ret["on_login_page"] = true
	}

	if display&(1<<1) == 2 {
		ret["on_school_page"] = true
	}

	if display&(1<<2) == 4 {
		ret["on_teacher_page"] = true
	}

	if display&(1<<3) == 8 {
		ret["on_student_page"] = true
	}

	ret["view_count"] = m["view"]
	ret["click_count"] = m["click"]

	// ads do not have preview_url
	if m["name"] != "" {
		ret["preview_url"] = cfg.Prefix + "/assets/adimg/" + m["name"]
	}

	if m["pending"] == "1" {
		ret["is_lock"] = true
	} else {
		ret["is_lock"] = false
	}

	return ret
}

func (client *AdsClient) GetAdById(id int64) (ret interface{}) {
	// client.Query(sqlStr, 0)

	GetFieldMap(&UserInfo{1, "name", "other"})
	dbret := dbclient.Query(client.Db, "select ad.*, file.name from "+tableAd+" ad left join ads_files file on ad.id = file.advertising_id where ad.id = ?", formatAdsResultSet, id)
	if len(dbret) >= 1 {
		// ret = ret[:1]
		ret = dbret[0]
	} else {
		ret = map[string]string{}
	}
	return
}

func (client *AdsClient) GetAllAds(page, items int) (ret []interface{}) {

	// sid, sname, utype := baseinfo.GetSchoolInfoFromUser(client.Db, "admin002")
	// log.Println("school info...", sid, sname, utype)

	if page == 0 {
		ret = dbclient.Query(client.Db, "select * from "+tableAd, formatAdsResultSet)
	} else {
		offset := page * items
		ret = dbclient.Query(client.Db, "select * from "+tableAd+" limit ? offset ? ", formatAdsResultSet, strconv.Itoa(items), strconv.Itoa(offset))
	}
	return
}

func (client *AdsClient) GetBaseInfo(username string) (int64, string, string) {
	return baseinfo.GetSchoolInfoFromUser(client.Db, username)
}

func (client *AdsClient) UploadFile() {

}

func (client *AdsClient) SetupDb() {
	// fmt.Println("setupdb...")
	// dir, _ := os.Getwd()
	// data, _ := ioutil.ReadFile(dir + "/config/sql")

	// sql := string(data)

	// stmtOut, _ := client.Db.Prepare(sql)
	// ret, err := stmtOut.Exec()
	// fmt.Println(ret, err)
}
