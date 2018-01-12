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
	// "strconv"
	// "html/template"
	// "github.com/satori/go.uuid"
	// return uuid.NewV4().String()
	"reflect"
	"strings"
	// "time"
	// "net"
)

var currentTable string = "assets"

type IAssetsClient interface {
	dbclient.IMysqlClient

	InsertOrUpdateAssets(id, filename, region, path string, size int) (sql.Result, bool)

	GetBaseInfo(username string) (int64, string, string)
}

type AssetsClient struct {
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

func (client *AssetsClient) GetBaseInfo(username string) (int64, string, string) {
	return baseinfo.GetSchoolInfoFromUser(client.Db, username)
}

func (client *AssetsClient) InsertOrUpdateAssets(id, filename, region, path string, size int) (sql.Result, bool) {
	sqlStr := "INSERT INTO " + currentTable + " (id, filename, region, path, size) values(?, ?, ?, ?, ?) " + " ON DUPLICATE KEY UPDATE filename=?, region=?, path=?, size=?"

	tx, err := client.Db.Begin()
	if err != nil {
		panic(err)
	}
	ret := dbclient.Exec(tx, sqlStr, id, filename, region, path, size, filename, region, path, size)
	log.Println(ret)
	tx.Commit()
	return ret, true
}
