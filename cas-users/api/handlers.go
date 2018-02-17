package api

import (
	"baseinfo"
	"cas-users/model"
	"cas-users/service"
	// "dbclient"
	"encoding/json"
	// "fmt"
	"github.com/gorilla/mux"
	// "hash/fnv"
	// fnv.New32a() h.Sum32() https://play.golang.org/p/_J2YysdEqE
	// "io"
	"log"
	"net/http"
	// "os"
	// "reflect"
	// "regexp"
	"errors"
	"strconv"
	// "strings"
	// "html/template"
	// "net"
)

type Operation int

const (
	OP_ADD Operation = iota
	OP_EDIT
	OP_DEL
	OP_VIEW
	OP_LIST
)

var MysqlClient service.IUsersClient

type Sizer interface {
	Size() int64
}

type BaseInfo struct {
	SchoolId int64
	UserName string
	UserType string
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	log.Println(r.Method)
	uname := mux.Vars(r)["username"]
	log.Println(uname)

	data := MysqlClient.GetUserByUsername(uname)
	ret, _ := json.Marshal(data)

	// check permission
	username := baseinfo.GetUsernameFromHeader(r)
	sid, sname, stype := MysqlClient.GetBaseInfo(username)
	if stype == "" {
		msg := "there is no user or invalid user with no type"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	var target *BaseInfo
	user := data.(map[string]interface{})
	if len(user) > 0 {
		school_id, _ := strconv.ParseInt(user["school_id"].(string), 10, 64)
		target = &BaseInfo{school_id, user["school_name"].(string), user["type"].(string)}
	}
	err := validOperation(OP_VIEW, &BaseInfo{sid, sname, stype}, target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(ret)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ret))

	return
}

//Parameter: type
//Super admin can only see school-admin
//School admin can only see users in same school
func GetUsers(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	log.Println("queries...", queries)

	username := baseinfo.GetUsernameFromHeader(r)
	sid, sname, stype := MysqlClient.GetBaseInfo(username)
	if stype == "" {
		msg := "there is no user or invalid user with no type"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	err := validOperation(OP_LIST, &BaseInfo{sid, sname, stype}, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tarType := queries.Get("type")
	page, _ := strconv.Atoi(queries.Get("page"))
	items, _ := strconv.Atoi(queries.Get("items"))

	log.Println(page, items, "page and items...")
	// if page is 0, then return all
	data := MysqlClient.GetAllUsers(page, items, sid, stype, tarType)

	ret, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(ret)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ret))
	return
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	username := baseinfo.GetUsernameFromHeader(r)
	if username == "" {
		log.Println("user name is empty from auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sid, sname, stype := MysqlClient.GetBaseInfo(username)
	if stype == "" {
		msg := "there is no user or invalid user with no type"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	user := new(model.User)
	log.Println(user)
	log.Println("insert user body...", r.Body)
	err := json.NewDecoder(r.Body).Decode(user)
	log.Println(err)
	log.Println("json user decoded...", user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("user...", user)

	err = validOperation(OP_ADD, &BaseInfo{sid, sname, stype}, &BaseInfo{user.SchoolId, user.Username, user.Type})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ret, succ := MysqlClient.InsertUser(user.Username, user.Name, user.Email, user.RoleIds,
		user.Password, user.SchoolName, user.Type, user.Phone, user.SchoolId, user.Activated, user.IsLock)

	adId, _ := ret.LastInsertId()
	log.Println("insert user ret...", ret, succ)
	w.Header().Set("Content-Type", "application/json")
	if succ {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.FormatInt(adId, 10)))
	} else {
		w.WriteHeader(503)
		w.Write([]byte(strconv.Itoa(-1)))
	}
	return
}

func EditUser(w http.ResponseWriter, r *http.Request) {
	user := new(model.User)
	log.Println(user)
	log.Println("edit user body...", r.Body)
	err := json.NewDecoder(r.Body).Decode(user)
	log.Println(err)
	log.Println("json user decoded...", user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("username...", user.Username)
	if user.Username == "" {
		http.Error(w, "user id is null", http.StatusBadRequest)
		log.Println("user id is null for update...")
		return
	}

	username := baseinfo.GetUsernameFromHeader(r)
	sid, sname, stype := MysqlClient.GetBaseInfo(username)
	if stype == "" {
		msg := "there is no user or invalid user with no type"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	err = validOperation(OP_EDIT, &BaseInfo{sid, sname, stype}, &BaseInfo{user.SchoolId, user.Username, user.Type})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ret, succ := MysqlClient.UpdateUser(user.Username, user.Name, user.Email, user.RoleIds,
		user.SchoolName, user.Type, user.Phone, user.SchoolId, user.Activated, user.IsLock)
	affected, _ := ret.RowsAffected()
	log.Println("affected...", affected)
	log.Println("update ret...", ret, succ)
	w.Header().Set("Content-Type", "application/json")
	if succ && affected >= 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("1"))
	} else {
		w.WriteHeader(503)
		w.Write([]byte("-1"))
	}
	return
}

func ResetPwd(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&data)
	log.Println(err)
	log.Println("json user pwd decoded...")

	// Validate permission
	username := baseinfo.GetUsernameFromHeader(r)
	log.Println("user name...", username)
	sid, sname, utype := MysqlClient.GetBaseInfo(username)
	if utype == "" {
		msg := "there is no user or invalid user with no type"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	tsid, tsname, tutype := MysqlClient.GetBaseInfo(data["username"])

	err = validOperation(OP_EDIT, &BaseInfo{sid, sname, utype}, &BaseInfo{tsid, tsname, tutype})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ret, succ := MysqlClient.UpdatePwd(data["username"], data["password"])
	affected, _ := ret.RowsAffected()
	log.Println("affected...", affected)
	log.Println("update ret...", ret, succ)
	w.Header().Set("Content-Type", "application/json")
	if succ && affected >= 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("1"))
	} else {
		w.WriteHeader(503)
		w.Write([]byte("-1"))
	}
	return
}

func UpdatePwd(w http.ResponseWriter, r *http.Request) {

	username := baseinfo.GetUsernameFromHeader(r)

	data := make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&data)
	log.Println(err)
	log.Println("json user pwd decoded...")

	if err = MysqlClient.CheckPwd(username, data["opassword"]); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ret, succ := MysqlClient.UpdatePwd(username, data["password"])
	affected, _ := ret.RowsAffected()
	log.Println("affected...", affected)
	log.Println("update ret...", ret, succ)
	w.Header().Set("Content-Type", "application/json")
	if succ && affected >= 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("1"))
	} else {
		w.WriteHeader(503)
		w.Write([]byte("-1"))
	}
	return
}

func DelUsers(w http.ResponseWriter, r *http.Request) {
	username := baseinfo.GetUsernameFromHeader(r)
	sid, _, stype := MysqlClient.GetBaseInfo(username)
	if stype == "" {
		msg := "there is no user or invalid user with no type"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	data := make(map[string][]string)
	log.Println("ids...", data)
	log.Println("delete users body...", r.Body)
	err := json.NewDecoder(r.Body).Decode(&data)
	log.Println(err)
	log.Println("json user ids decoded...", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ids, ok := data["ids"]; ok {
		ret, succ := MysqlClient.DelUsers(ids, sid, stype)
		affected, _ := ret.RowsAffected()
		if succ && affected >= 0 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("1"))
		} else {
			w.WriteHeader(503)
			w.Write([]byte("-1"))
		}
		return
	} else {
		http.Error(w, "no ids in body to delete...", http.StatusBadRequest)
		log.Println("no ids in body to delete......")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("1"))
	return
}

func validOperation(op Operation, origin, target *BaseInfo) error {

	isStaff := func(user *BaseInfo) bool {
		return isTeacher(user.UserType) || isStudent(user.UserType)
	}
	inSameSchool := func(user1, user2 *BaseInfo) bool {
		return user1 != nil && user2 != nil && user1.SchoolId == user2.SchoolId
	}
	isSamePerson := func(user1, user2 *BaseInfo) bool {
		return user1 != nil && user2 != nil && user1.UserName == user2.UserName
	}
	isSuperAdmin := func(user *BaseInfo) bool {
		return user != nil && user.UserType == "admin"
	}
	isSchoolAdmin := func(user *BaseInfo) int {
		if user == nil {
			return 0
		}
		if user.UserType == "school_admin" {
			return 1
		} else if user.UserType == "school" {
			return 2
		} else {
			return 0
		}
	}

	isSuperAdminOp := func(user1, user2 *BaseInfo) bool {
		return isSuperAdmin(user1) && (isSchoolAdmin(user2) == 1 || isSuperAdmin(user2))
	}
	isSchoolAdminOp := func(user1, user2 *BaseInfo) bool {
		return isSchoolAdmin(user1) > 0 && !isSuperAdmin(user2) &&
			inSameSchool(user1, user2)
	}

	isAdminOpt := func(user1, user2 *BaseInfo) bool {
		return isSuperAdminOp(user1, user2) || isSchoolAdminOp(user1, user2)
	}

	if op == OP_ADD { // admin only
		if !isAdminOpt(origin, target) {
			log.Println("user is not admin to add...")
			return errors.New("user does not have operation permission")
		}
	} else if op == OP_EDIT {
		if isStaff(origin) && !isSamePerson(origin, target) {
			log.Println("user is not admin that cannot edit others information...")
			return errors.New("user does not have operation permission")
		} else if !isAdminOpt(origin, target) {
			log.Println("user is not admin to edit...")
			return errors.New("user does not have operation permission")
		}
	} else if op == OP_DEL { // The same as add
		if !isSuperAdmin(origin) || isSchoolAdmin(origin) == 0 {
			log.Println("user is not admin to delete or not in the same domain...")
			return errors.New("user does not have operation permission")
		}
	} else if op == OP_VIEW { // The same as edit
		if isStaff(origin) && !isSamePerson(origin, target) {
			log.Println("user is not admin that cannot del others information...")
			return errors.New("user does not have operation permission")
		} else if !isAdminOpt(origin, target) {
			log.Println("user is not admin to view or not in the same domain...")
			return errors.New("user does not have operation permission")
		}
	} else if op == OP_LIST { // The same as add
		if !isSuperAdmin(origin) && isSchoolAdmin(origin) == 0 {
			log.Println("user is not admin to list...")
			return errors.New("user does not have operation permission")
		}
	} else {
		return errors.New("unknown operation")
	}
	return nil
}

func isTeacher(uType string) bool {
	return uType == "teacher"
}

func isStudent(uType string) bool {
	return uType == "student"
}
