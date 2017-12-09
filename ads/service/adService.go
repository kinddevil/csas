package service

import (
	"dbclient"
	// "encoding/json"
	"database/sql"
	"fmt"
	// "github.com/gorilla/mux"
	// "io"
	// "io/ioutil"
	// "log"
	// "net/http"
	// "os"
	"strconv"
	// "html/template"
	"reflect"
	"strings"
	// "net"
)

var tableAd string = "advertising"

type IAdsClient interface {
	dbclient.IMysqlClient
	InsertAd(imageName, imageLink, schoolIds, province, city, title string, displayPages int) bool
	UpdateAd(id int, imageName, imageLink, schoolIds, province, city, title string, displayPages int) bool
	GetAdById(id int64) []interface{}
	GetAllAds(page, items int) []interface{}
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

func (client *AdsClient) GetAdById(id int64) (ret []interface{}) {
	// client.Query(sqlStr, 0)
	GetFieldMap(&UserInfo{1, "name", "other"})
	ret = dbclient.Query(client.Db, "select * from users_test", nil)
	return
}

func (client *AdsClient) InsertAd(imageName, imageLink, schoolIds, province, city, title string, displayPages int) bool {
	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}
	sql, vals := dbclient.BuildInsert(tableAd, dbclient.ParamsPairs(
		"img", imageName,
		"link", imageLink,
		"province", province,
		"city", city,
		"display_pages", displayPages,
		"title", title,
	),
	)
	ret := dbclient.Exec(tx, sql, vals...)
	fmt.Println(ret)
	tx.Commit()
	return true
}

func (client *AdsClient) UpdateAd(id int, imageName, imageLink, schoolIds, province, city, title string, displayPages int) bool {
	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}
	sql, vals := dbclient.BuildUpdate(tableAd, dbclient.ParamsPairs(
		"img", imageName,
		"link", imageLink,
		"province", province,
		"city", city,
		"display_pages", displayPages,
		"title", title,
	), dbclient.ParamsPairs(
		"id", id,
	),
	)
	ret := dbclient.Exec(tx, sql, vals...)
	fmt.Println(ret)
	tx.Commit()
	return true
}

func (client *AdsClient) GetAllAds(page, items int) (ret []interface{}) {
	if page == 0 {
		ret = dbclient.Query(client.Db, "select * from "+tableAd, nil)
	} else {
		offset := page * items
		ret = dbclient.Query(client.Db, "select * from "+tableAd+" limit ? offset ? ", nil, strconv.Itoa(items), strconv.Itoa(offset))
	}
	return
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
