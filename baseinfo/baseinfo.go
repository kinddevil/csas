package baseinfo

import (
	"database/sql"
	"dbclient"
	"log"
	"net/http"
	// "reflect"
	"strconv"
	"strings"
)

const (
	userTable = "user"
)

type UserInfo struct {
	Username   string
	SchoolId   int64
	SchoolName string
	Type       string
}

func GetUsernameFromHeader(r *http.Request) string {
	username := r.Header.Get("X-Requested-Un")
	log.Println("user name is", username)
	// for k, v := range r.Header {
	// 	log.Println(k, v)
	// }
	return strings.Trim(username, " ")
}

func GetSchoolInfoFromUser(db *sql.DB, username string) (int64, string, string) {
	users := GetSchoolsFromUser(db, username)
	if len(users) > 0 {
		user := users[0].(*UserInfo)
		return user.SchoolId, user.SchoolName, user.Type
	}
	return 0, "", ""
}

func GetSchoolsFromUser(db *sql.DB, username string) []interface{} {
	sqlStr := "select username, school_id, school_name, type from " + userTable + " where username = ?"
	ret := dbclient.Query(db, sqlStr, func(m map[string]string) interface{} {
		user := &UserInfo{}
		user.Username = m["username"]
		user.SchoolName = m["school_name"]
		user.Type = m["type"]
		schoolId, _ := strconv.Atoi(m["school_id"])
		user.SchoolId = int64(schoolId)
		return user
	}, username)
	// return map[int64]string{
	// 	0: "",
	// }
	return ret
}
