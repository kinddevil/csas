package dbclient

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"sync"
	"time"
)

var (
	mysqlClient *MysqlClient
	monce       sync.Once
)

type IMysqlClient interface {
	Open()
	Close()
	Ping()
	Seed()
}

type MysqlClient struct {
	db *sql.DB
}

func InitMysql() *MysqlClient {
	monce.Do(func() {
		mysqlClient = &MysqlClient{}
	})
	return mysqlClient
}

func (client *MysqlClient) Open() {
	var err error
	client.db, err = sql.Open("mysql", "user:pass@tcp(localhost:3306)/db?charset=utf8&parseTime=true")
	// client.db, err = sql.Open("mysql", "casuser:Cassuser365@tcp(db4free.net:3306)/db?charset=utf8&parseTime=true")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	client.db.SetMaxOpenConns(5)
	client.db.SetMaxIdleConns(5)
}

func (client *MysqlClient) Close() {
	if client.db != nil {
		client.db.Close()
	}
}

func (client *MysqlClient) Ping() {
	err := client.db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}

func (client *MysqlClient) Seed() {
	// RowsAffected changes, err := ret.RowsAffected()
	stmtOut0, _ := client.db.Prepare("CREATE TABLE IF NOT EXISTS users_test (id bigint primary key, name text)")
	stmtOut1, _ := client.db.Prepare("truncate table users_test")
	stmtOut0.Exec()
	stmtOut1.Exec()
	tx, _ := client.db.Begin()
	for i := 0; i < 100; i++ {
		sql := "insert into users_test value(" + strconv.Itoa(i) + ", 'user00" + strconv.Itoa(i) + "')"
		fmt.Println(sql)
		stmtOut2, _ := tx.Prepare(sql)
		stmtOut2.Exec()
	}
	tx.Commit()
}

func test() {
	db, err := sql.Open("mysql", "user:pass@tcp(localhost:3306)/db?charset=utf8&parseTime=true")

	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	tx, err := db.Begin()

	stmtOut, err := tx.Prepare("SELECT id FROM user WHERE name = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	var id int                              // we "scan" the result in here// Query the square-number of 13
	err = stmtOut.QueryRow("bar").Scan(&id) // WHERE number = 13
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Printf("The id 1 is: %d \n", id)

	// Execute the query
	rows, err := tx.Query("SELECT * FROM user")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
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

	time.Sleep(10000 * time.Millisecond)
}
