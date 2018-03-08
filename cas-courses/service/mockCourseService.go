package service

import (
	// "encoding/json"
	// "fmt"
	"database/sql"
	// cfg "cas-dicts/config"
	// "github.com/gorilla/mux"
	// "io"
	// "baseinfo"
	// "io/ioutil"
	"log"
	// "net/http"
	// "os"
	// "strconv"
	// "html/template"
	// "github.com/satori/go.uuid"
	// return uuid.NewV4().String()
	// "reflect"
	// "strings"
	"time"
	// "net"
  "github.com/stretchr/testify/mock"
)

// type IMockCalendarClient interface {
// 	dbclient.IMysqlClient

// 	GetCalendarById(id, schoolId int64) (ret interface{})
// 	GetAllCalendars(page, items int, dtype string, schoolId int64, isVisible string) (ret []interface{})

// 	InsertCalendar(name string, schoolId int64, desc, ctype string,
// 		isActive bool, start, end time.Time, event string, isVisible bool) (sql.Result, bool)
// 	UpdateCalendar(id, uschoolId int64, name string, schoolId int64, desc, ctype string,
// 		isActive bool, start, end time.Time, event string, isVisible bool) (sql.Result, bool)
// 	DelCalendarById(id int64, schoolId int64) (sql.Result, bool)
// 	DelCalendarByIdReal(id int64, schoolId int64) (sql.Result, bool)
// 	DelCalendars(ids []int64, schoolId int64) (sql.Result, bool)

// 	GetBaseInfo(username string) (int64, string, string)
// }

// type IMysqlClient interface {
// 	Init(dbUrl string)
// 	Open(dbUrl string)
// 	Close()
// 	Ping()
// 	MaxConns(int)
// 	MaxIdleConns(int)
// 	Seed()
// }

type MockMysqkClient struct {

}

type MockCalendarClient struct {
	MockMysqkClient
  mock.Mock
}

func (client *MockCalendarClient) GetAllCalendars(page, items int, ctype string, schoolId int64, isVisible string) []interface{} {
	args := client.Mock.Called(page, items, ctype, schoolId, isVisible)
	log.Println("GetAllCalendars mock return...", args.Get(0))
  return args.Get(0).([]interface{})
}

func (client *MockCalendarClient) GetBaseInfo(username string) (int64, string, string) {
	args := client.Mock.Called(username)
	return args.Get(0).(int64), args.Get(1).(string), args.Get(2).(string)
}

func (client *MockCalendarClient) GetCalendarById(id, schoolId int64) interface{} {
	return nil
}

func (client *MockCalendarClient) InsertCalendar(name string, schoolId int64, desc, ctype string,
	isActive bool, start, end time.Time, event string, isVisible bool) (sql.Result, bool) {

	return nil, false
}

func (client *MockCalendarClient) UpdateCalendar(id, uschoolId int64, name string, schoolId int64, desc, ctype string,
	isActive bool, start, end time.Time, event string, isVisible bool) (sql.Result, bool) {

	return nil, false
}

func (client *MockCalendarClient) DelCalendarById(id int64, schoolId int64) (sql.Result, bool) {
	return nil, false
}

func (client *MockCalendarClient) DelCalendarByIdReal(id int64, schoolId int64) (sql.Result, bool) {
	return nil, false
}

func (client *MockCalendarClient) DelCalendars(ids []int64, schoolId int64) (sql.Result, bool) {
	return nil, false
}

func (client *MockMysqkClient) Close() {
}

func (client *MockMysqkClient) Init(dbUrl string) {
}

func (client *MockMysqkClient) Open(dbUrl string) {
}

func (client *MockMysqkClient) Ping() {
}

func (client *MockMysqkClient) MaxConns(cons int) {
}

func (client *MockMysqkClient) MaxIdleConns(idles int) {
}

func (client *MockMysqkClient) Seed() {
}


