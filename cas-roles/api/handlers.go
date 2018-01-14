package api

import (
	"baseinfo"
	"cas-roles/model"
	"cas-roles/service"
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
	"strconv"
	"strings"
	// "html/template"
	// "net"
)

var MysqlClient service.IRolesClient

type Sizer interface {
	Size() int64
}

func GetRole(w http.ResponseWriter, r *http.Request) {
	username := baseinfo.GetUsernameFromHeader(r)
	if username == "" {
		log.Println("user name is empty to add role")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	log.Println(r.Method)
	var id, _ = strconv.Atoi(mux.Vars(r)["roleId"])
	log.Println(id)

	data := MysqlClient.GetRoleById(int64(id), username)
	ret, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(ret)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ret))

	return
}

func GetRoles(w http.ResponseWriter, r *http.Request) {
	username := baseinfo.GetUsernameFromHeader(r)
	if username == "" {
		log.Println("user name is empty to add role")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	queries := r.URL.Query()
	log.Println("queries...", queries)

	page, _ := strconv.Atoi(queries.Get("page"))
	items, _ := strconv.Atoi(queries.Get("items"))

	log.Println(page, items, "page and items...")
	// if page is 0, then return all
	data := MysqlClient.GetAllRoles(page, items, username)

	ret, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(ret)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ret))
	return
}

func AddRole(w http.ResponseWriter, r *http.Request) {
	username := baseinfo.GetUsernameFromHeader(r)
	if username == "" {
		log.Println("user name is empty to add role")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	role := new(model.Role)
	log.Println(role)
	log.Println("insert role body...", r.Body)
	err := json.NewDecoder(r.Body).Decode(role)
	log.Println(err)
	log.Println("json role decoded...", role)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("role...", role)

	userids := []string{}
	log.Println(role.UserIds, "userids...")
	if role.UserIds != "" {
		userids = strings.Split(role.UserIds, ",")
		for i, val := range userids {
			userids[i] = strings.Trim(val, " ")
		}
	}
	log.Println(userids, "userids...")

	ret, succ := MysqlClient.InsertRole(role.Name, username, role.Permissions, userids)
	w.Header().Set("Content-Type", "application/json")

	if !succ {
		w.WriteHeader(503)
		w.Write([]byte(strconv.Itoa(-1)))
		return
	}

	adId, _ := ret.LastInsertId()
	log.Println("insert role ret...", ret, succ)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatInt(adId, 10)))
	return
}

func EditRole(w http.ResponseWriter, r *http.Request) {
	username := baseinfo.GetUsernameFromHeader(r)
	if username == "" {
		log.Println("user name is empty to add role")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	role := new(model.Role)
	log.Println(role)
	log.Println("edit role body...", r.Body)
	err := json.NewDecoder(r.Body).Decode(role)
	log.Println(err)
	log.Println("json role decoded...", role)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("role id", role.Id)
	if role.Id == 0 {
		http.Error(w, "role id is null", http.StatusBadRequest)
		log.Println("role id is null for update...")
		return
	}

	userids := []string{}
	if role.UserIds != "" {
		userids = strings.Split(role.UserIds, ",")
		for i, val := range userids {
			userids[i] = strings.Trim(val, " ")
		}
	}
	log.Println("username...", username)
	ret, succ := MysqlClient.UpdateRole(role.Id, role.Name, username, role.Permissions, userids)
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

func DelRoles(w http.ResponseWriter, r *http.Request) {
	data := make(map[string][]int64)
	log.Println("ids...", data)
	log.Println("delete roles body...", r.Body)
	err := json.NewDecoder(r.Body).Decode(&data)
	log.Println(err)
	log.Println("json role ids decoded...", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ids, ok := data["ids"]; ok {
		ret, succ := MysqlClient.DelRoles(ids)
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
