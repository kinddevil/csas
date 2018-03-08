package api

import (
	"baseinfo"
	"cas-courses/model"
	"cas-courses/service"
	// "dbclient"
	"encoding/json"
	// "fmt"
	"github.com/gorilla/mux"
	// "hash/fnv"
	// fnv.New32a() h.Sum32() https://play.golang.org/p/_J2YysdEqE
	// "io"
	"log"
	"net/http"
	// "time"
	"errors"
	// "os"
	// "reflect"
	// "regexp"
	"strconv"
	// "strings"
	// "html/template"
	// "net"
)

var MysqlClient service.ICourseClient

type Sizer interface {
	Size() int64
}

type BaseInfo struct {
	SchoolName string
	SchoolId   int64
	UserType   string
}

var (
	SCHOOL_ADMINS = map[string]int{
		"school":       1,
		"school_admin": 1,
	}
)

func GetCourse(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method)
	var courseId, _ = strconv.Atoi(mux.Vars(r)["cid"])
	log.Println(courseId)

	username := baseinfo.GetUsernameFromHeader(r)
	sid, _, utype := MysqlClient.GetBaseInfo(username)
	if utype == "" {
		msg := "there is no user or invalid user with no type"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	data := MysqlClient.GetCourseById(int64(courseId), sid)
	ret, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(ret)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ret))

	return
}

func GetCourses(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	log.Println("queries...", queries)

	username := baseinfo.GetUsernameFromHeader(r)
	sid, _, utype := MysqlClient.GetBaseInfo(username)
	if utype == "" {
		msg := "there is no user or invalid user with no type"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	ctype := queries.Get("type")
	visible := queries.Get("is_visible")
	page, _ := strconv.Atoi(queries.Get("page"))
	items, _ := strconv.Atoi(queries.Get("items"))

	log.Println(page, items, ctype, visible, "page and items...")
	// if page is 0, then return all
	data := MysqlClient.GetAllCourses(page, items, ctype, sid, visible)

	ret, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(ret)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ret))
	return
}

func AddCourse(w http.ResponseWriter, r *http.Request) {
	username := baseinfo.GetUsernameFromHeader(r)
	if username == "" {
		log.Println("user name is empty to add course")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sid, sname, utype := MysqlClient.GetBaseInfo(username)
	if utype == "" {
		msg := "there is no user or invalid user with no type"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	course := new(model.Course)
	log.Println(course)
	log.Println("insert course body...", r.Body)
	err := json.NewDecoder(r.Body).Decode(course)
	log.Println(err)
	log.Println("json course decoded...", course)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("course...", course)

	if perr := checkPermission(&BaseInfo{sname, sid, utype}, &course.SchoolId); perr != nil {
		http.Error(w, perr.Error(), http.StatusBadRequest)
		return
	}

	ret, succ := MysqlClient.InsertCourse(course.Name, course.SchoolId,
		course.Desc, course.Type, course.IsActive, course.Start, course.End, course.Events, course.IsVisible)
	w.Header().Set("Content-Type", "application/json")

	if !succ {
		w.WriteHeader(503)
		w.Write([]byte(strconv.Itoa(-1)))
		return
	}

	adId, _ := ret.LastInsertId()
	log.Println("insert course ret...", ret, succ)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatInt(adId, 10)))
	return
}

func EditCourse(w http.ResponseWriter, r *http.Request) {
	username := baseinfo.GetUsernameFromHeader(r)
	if username == "" {
		log.Println("user name is empty to add course")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	course := new(model.Course)
	log.Println(course)
	log.Println("edit course body...", r.Body)
	err := json.NewDecoder(r.Body).Decode(course)
	log.Println(err)
	log.Println("json course decoded...", course)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("course id", course.Id)
	if course.Id == 0 {
		http.Error(w, "course id is null", http.StatusBadRequest)
		log.Println("course id is null for update...")
		return
	}

	sid, sname, utype := MysqlClient.GetBaseInfo(username)
	if utype == "" {
		msg := "there is no user or invalid user with no type"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if perr := checkPermission(&BaseInfo{sname, sid, utype}, &course.SchoolId); perr != nil {
		http.Error(w, perr.Error(), http.StatusBadRequest)
		return
	}

	ret, succ := MysqlClient.UpdateCourse(course.Id, sid, course.Name, course.SchoolId,
		course.Desc, course.Type, course.IsActive, course.Start, course.End, course.Events, course.IsVisible)
	w.Header().Set("Content-Type", "application/json")

	if !succ {
		w.WriteHeader(503)
		w.Write([]byte("-1"))
		return
	}

	affected, _ := ret.RowsAffected()
	log.Println("affected...", affected)
	log.Println("update ret...", ret, succ)

	if affected >= 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("1"))
	} else {
		w.WriteHeader(503)
		w.Write([]byte("-1"))
	}
	return
}

func DelCourse(w http.ResponseWriter, r *http.Request) {
	var courseId, _ = strconv.Atoi(mux.Vars(r)["cid"])
	log.Println(courseId)

	username := baseinfo.GetUsernameFromHeader(r)
	if username == "" {
		log.Println("user name is empty to add course")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sid, _, utype := MysqlClient.GetBaseInfo(username)
	if utype == "" {
		msg := "there is no user or invalid user with no type"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// ret, succ := MysqlClient.DelCourseByIdReal(int64(CourseId))
	ret, succ := MysqlClient.DelCourseById(int64(courseId), sid)
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

func DelCourses(w http.ResponseWriter, r *http.Request) {
	data := make(map[string][]int64)
	log.Println("ids...", data)
	log.Println("delete course body...", r.Body)
	err := json.NewDecoder(r.Body).Decode(&data)
	log.Println(err)
	log.Println("json course ids decoded...", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := baseinfo.GetUsernameFromHeader(r)
	if username == "" {
		log.Println("user name is empty to add course")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sid, _, utype := MysqlClient.GetBaseInfo(username)
	if utype == "" {
		msg := "there is no user or invalid user with no type"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if ids, ok := data["ids"]; ok {
		ret, succ := MysqlClient.DelCourses(ids, sid)
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

func checkPermission(user *BaseInfo, schoolId *int64) error {
	log.Println(user, schoolId, *schoolId, "checkper")
	if _, ok := SCHOOL_ADMINS[user.UserType]; user == nil || !ok {
		return errors.New("User does not have permission!")
	}

	if schoolId != nil && *schoolId != user.SchoolId {
		return errors.New("User does not have permission for other school user!")
	}
	return nil
}
