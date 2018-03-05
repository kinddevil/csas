package api

import (
	"baseinfo"
	"cas-calendar/model"
	"cas-calendar/service"
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

var MysqlClient service.ICalendarClient

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

func GetCalendar(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method)
	var calendarId, _ = strconv.Atoi(mux.Vars(r)["cid"])
	log.Println(calendarId)

	username := baseinfo.GetUsernameFromHeader(r)
	sid, _, utype := MysqlClient.GetBaseInfo(username)
	if utype == "" {
		msg := "there is no user or invalid user with no type"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	data := MysqlClient.GetCalendarById(int64(calendarId), sid)
	ret, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(ret)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ret))

	return
}

func GetCalendars(w http.ResponseWriter, r *http.Request) {
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
	data := MysqlClient.GetAllCalendars(page, items, ctype, sid, visible)

	ret, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(ret)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ret))
	return
}

func AddCalendar(w http.ResponseWriter, r *http.Request) {
	username := baseinfo.GetUsernameFromHeader(r)
	if username == "" {
		log.Println("user name is empty to add Calendar")
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

	calendar := new(model.Calendar)
	log.Println(calendar)
	log.Println("insert Calendar body...", r.Body)
	err := json.NewDecoder(r.Body).Decode(calendar)
	log.Println(err)
	log.Println("json Calendar decoded...", calendar)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Calendar...", calendar)

	if perr := checkPermission(&BaseInfo{sname, sid, utype}, &calendar.SchoolId); perr != nil {
		http.Error(w, perr.Error(), http.StatusBadRequest)
		return
	}

	ret, succ := MysqlClient.InsertCalendar(calendar.Name, calendar.SchoolId,
		calendar.Desc, calendar.Type, calendar.IsActive, calendar.Start, calendar.End, calendar.Events, calendar.IsVisible)
	w.Header().Set("Content-Type", "application/json")

	if !succ {
		w.WriteHeader(503)
		w.Write([]byte(strconv.Itoa(-1)))
		return
	}

	adId, _ := ret.LastInsertId()
	log.Println("insert Calendar ret...", ret, succ)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatInt(adId, 10)))
	return
}

func EditCalendar(w http.ResponseWriter, r *http.Request) {
	username := baseinfo.GetUsernameFromHeader(r)
	if username == "" {
		log.Println("user name is empty to add Calendar")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	calendar := new(model.Calendar)
	log.Println(calendar)
	log.Println("edit Calendar body...", r.Body)
	err := json.NewDecoder(r.Body).Decode(calendar)
	log.Println(err)
	log.Println("json Calendar decoded...", calendar)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Calendar id", calendar.Id)
	if calendar.Id == 0 {
		http.Error(w, "calendar id is null", http.StatusBadRequest)
		log.Println("calendar id is null for update...")
		return
	}

	sid, sname, utype := MysqlClient.GetBaseInfo(username)
	if utype == "" {
		msg := "there is no user or invalid user with no type"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if perr := checkPermission(&BaseInfo{sname, sid, utype}, &calendar.SchoolId); perr != nil {
		http.Error(w, perr.Error(), http.StatusBadRequest)
		return
	}

	ret, succ := MysqlClient.UpdateCalendar(calendar.Id, sid, calendar.Name, calendar.SchoolId,
		calendar.Desc, calendar.Type, calendar.IsActive, calendar.Start, calendar.End, calendar.Events, calendar.IsVisible)
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

func DelCalendar(w http.ResponseWriter, r *http.Request) {
	var CalendarId, _ = strconv.Atoi(mux.Vars(r)["cid"])
	log.Println(CalendarId)

	username := baseinfo.GetUsernameFromHeader(r)
	if username == "" {
		log.Println("user name is empty to add Calendar")
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

	// ret, succ := MysqlClient.DelCalendarByIdReal(int64(CalendarId))
	ret, succ := MysqlClient.DelCalendarById(int64(CalendarId), sid)
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

func DelCalendars(w http.ResponseWriter, r *http.Request) {
	data := make(map[string][]int64)
	log.Println("ids...", data)
	log.Println("delete Calendars body...", r.Body)
	err := json.NewDecoder(r.Body).Decode(&data)
	log.Println(err)
	log.Println("json Calendar ids decoded...", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := baseinfo.GetUsernameFromHeader(r)
	if username == "" {
		log.Println("user name is empty to add Calendar")
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
		ret, succ := MysqlClient.DelCalendars(ids, sid)
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

func Test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("this is a test for calendar"))
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
