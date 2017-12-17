package dbclient

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strconv"
	"strings"
	// "unsafe"
	// "sync"
	// "time"
)

// var (
// 	mysqlClient *MysqlClient
// 	monce       sync.Once
// )

type IMysqlClient interface {
	Init(dbUrl string)
	Open(dbUrl string)
	Close()
	Ping()
	MaxConns(int)
	MaxIdleConns(int)
	Seed()
}

type MysqlClient struct {
	Db *sql.DB
}

func (client *MysqlClient) Init(dbUrl string) {
	client.Open(dbUrl)
	client.MaxConns(5)
	client.MaxIdleConns(5)
	client.Ping()
}

func (client *MysqlClient) Open(dbUrl string) {
	var err error
	client.Db, err = sql.Open("mysql", dbUrl)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
}

func (client *MysqlClient) MaxConns(num int) {
	client.Db.SetMaxOpenConns(num)
}

func (client *MysqlClient) MaxIdleConns(num int) {
	client.Db.SetMaxIdleConns(num)
}

func (client *MysqlClient) Close() {
	if client.Db != nil {
		client.Db.Close()
	}
}

func (client *MysqlClient) Ping() {
	err := client.Db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}

func (client *MysqlClient) Seed() {
	// RowsAffected changes, err := ret.RowsAffected()
	stmtOut0, _ := client.Db.Prepare("CREATE TABLE IF NOT EXISTS users_test (id bigint primary key, name text)")
	stmtOut1, _ := client.Db.Prepare("truncate table users_test")
	stmtOut0.Exec()
	stmtOut1.Exec()
	tx, _ := client.Db.Begin()
	for i := 0; i < 100; i++ {
		sql := "insert into users_test value(" + strconv.Itoa(i) + ", 'user00" + strconv.Itoa(i) + "')"
		fmt.Println(sql)
		stmtOut2, _ := tx.Prepare(sql)
		stmtOut2.Exec()
	}
	stmtOut3, _ := tx.Prepare("insert into users_test(id) value(1001)")
	stmtOut3.Exec()
	tx.Commit()
}

func Exec(tx *sql.Tx, sqlStr string, args ...interface{}) sql.Result {
	stmtOut, err := tx.Prepare(sqlStr)
	checkErr(err)
	defer stmtOut.Close()

	ret, err := stmtOut.Exec(args...)
	checkErr(err)
	return ret
}

func Query(db *sql.DB, sqlStr string, parseFun func(map[string]string) interface{}, args ...interface{}) []interface{} {
	//Return slice of interface
	fmt.Println(db, sqlStr, parseFun)
	ret := []interface{}{}

	// fmt.Println("ret...", ret, unsafe.Sizeof(ret[0]), unsafe.Sizeof(1), unsafe.Sizeof(true))
	stmtOut, err := db.Prepare(sqlStr)
	checkErr(err)
	defer stmtOut.Close()

	rows, err := stmtOut.Query(args...)
	checkErr(err)
	columns, err := rows.Columns()
	checkErr(err)
	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	warning := false
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		checkErr(err)

		mrow := make(map[string]string)
		for i, col := range values {
			colName := columns[i]
			if !warning {
				if _, ok := mrow[colName]; ok {
					fmt.Println("Warning: query field has same name of", colName)
					warning = true
				}
			}
			mrow[colName] = string(col)

			// TODO: Null Handler
			// if col == nil {
			// 	value = ""
			// } else {
			// 	value = string(col)
			// }
		}
		if parseFun == nil {
			ret = append(ret, mrow)
		} else {
			obj := parseFun(mrow)
			ret = append(ret, obj)
		}
	}
	checkErr(rows.Err())
	// fmt.Println("ret...", ret, len(ret), cap(ret), unsafe.Sizeof(ret[0]), unsafe.Sizeof(1), unsafe.Sizeof(true))
	return ret
}

func QueryWithTran(tx *sql.Tx, sqlStr string, parseFun func(map[string]string) interface{}, args ...interface{}) []interface{} {
	//Return slice of interface
	fmt.Println(tx, sqlStr, parseFun)
	ret := []interface{}{}

	// fmt.Println("ret...", ret, unsafe.Sizeof(ret[0]), unsafe.Sizeof(1), unsafe.Sizeof(true))
	stmtOut, err := tx.Prepare(sqlStr)
	checkErr(err)
	defer stmtOut.Close()

	rows, err := stmtOut.Query(args...)
	checkErr(err)
	columns, err := rows.Columns()
	checkErr(err)
	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	warning := false
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		checkErr(err)

		mrow := make(map[string]string)
		for i, col := range values {
			colName := columns[i]
			if !warning {
				if _, ok := mrow[colName]; ok {
					fmt.Println("Warning: query field has same name of", colName)
					warning = true
				}
			}
			mrow[colName] = string(col)

			// TODO: Null Handler
			// if col == nil {
			// 	value = ""
			// } else {
			// 	value = string(col)
			// }
		}
		if parseFun == nil {
			ret = append(ret, mrow)
		} else {
			obj := parseFun(mrow)
			ret = append(ret, obj)
		}
	}
	checkErr(rows.Err())
	// fmt.Println("ret...", ret, len(ret), cap(ret), unsafe.Sizeof(ret[0]), unsafe.Sizeof(1), unsafe.Sizeof(true))
	return ret
}

func ParamsPairs(args ...interface{}) map[string]interface{} {
	length := len(args)
	ret := map[string]interface{}{}
	for i := 0; i <= length/2; i++ {
		keyInd := i * 2
		valInd := keyInd + 1
		if valInd < length {
			ret[args[keyInd].(string)] = args[valInd]
		}
	}
	return ret
}

func BuildInsert(tableName string, pairs map[string]interface{}) (string, []interface{}) {
	length := len(pairs)
	vals := make([]interface{}, 0, length)
	fields := make([]string, 0, length)
	placeHolders := make([]string, 0, length)
	for k, v := range pairs {
		vals = append(vals, v)
		fields = append(fields, k)
		placeHolders = append(placeHolders, "?")
	}
	sql := fmt.Sprintf("INSERT INTO %s(%s) values(%s)", tableName, strings.Join(fields, ","), strings.Join(placeHolders, ","))
	return sql, vals
}

func BuildUpdate(tableName string, pairsVal map[string]interface{}, pairsCond map[string]interface{}) (string, []interface{}) {

	lenVal := len(pairsVal)
	lenCond := len(pairsCond)
	length := lenVal + lenCond

	if lenCond < 1 {
		panic(fmt.Sprintf("No conditions for update for table %s !!!", tableName))
	}

	vals := make([]interface{}, 0, length)
	targets := make([]string, 0, lenVal)
	conds := make([]string, 0, lenCond)

	// TODO:eliminate null values
	for k, v := range pairsVal {
		// targets = append(targets, fmt.Sprintf("%s=?", k))
		targets = append(targets, k+"=?")
		vals = append(vals, v)
	}

	for k, v := range pairsCond {
		// conds = append(conds, fmt.Sprintf("%s=?", k))
		conds = append(conds, k+"=?")
		vals = append(vals, v)
	}

	sql := fmt.Sprintf("UPDATE %s SET %s where %s", tableName, strings.Join(targets, ","), strings.Join(conds, " and "))

	return sql, vals
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

func checkErr(err error) {
	if err != nil {
		panic(err.Error()) // TODO: proper error handling
	}
}
