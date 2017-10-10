package service

import (
	"dbclient"
	// "encoding/json"
	"database/sql"
	"fmt"
	// "github.com/gorilla/mux"
	// "io"
	// "log"
	// "net/http"
	// "os"
	// "strconv"
	// "html/template"
	"reflect"
	"strings"
	// "net"
)

type IAdsClient interface {
	dbclient.IMysqlClient
	GetAdById(id int64) []interface{}
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
	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}
	ret = dbclient.Query(tx, "select * from users_test", nil)
	tx.Commit()
	return
}
