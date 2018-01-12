package api

import (
	"baseinfo"
	"cas-dicts/model"
	"cas-dicts/service"
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
	// "strings"
	// "html/template"
	// "net"
)

var MysqlClient service.IDictsClient

type Sizer interface {
	Size() int64
}

func GetDict(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method)
	var dictId, _ = strconv.Atoi(mux.Vars(r)["dictId"])
	log.Println(dictId)

	data := MysqlClient.GetDictById(int64(dictId))
	ret, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(ret)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ret))

	return
}

func GetDicts(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	log.Println("queries...", queries)

	dtype := queries.Get("type")
	page, _ := strconv.Atoi(queries.Get("page"))
	items, _ := strconv.Atoi(queries.Get("items"))

	log.Println(page, items, "page and items...")
	// if page is 0, then return all
	data := MysqlClient.GetAllDicts(page, items, dtype)

	ret, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(ret)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ret))
	return
}

func AddDict(w http.ResponseWriter, r *http.Request) {
	username := baseinfo.GetUsernameFromHeader(r)
	if username == "" {
		log.Println("user name is empty to add dict")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dict := new(model.Dict)
	log.Println(dict)
	log.Println("insert dict body...", r.Body)
	err := json.NewDecoder(r.Body).Decode(dict)
	log.Println(err)
	log.Println("json dict decoded...", dict)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("dict...", dict)

	ret, succ := MysqlClient.InsertDict(username, dict.Key, dict.Desc, dict.Type)
	w.Header().Set("Content-Type", "application/json")

	if !succ {
		w.WriteHeader(503)
		w.Write([]byte(strconv.Itoa(-1)))
		return
	}

	adId, _ := ret.LastInsertId()
	log.Println("insert dict ret...", ret, succ)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatInt(adId, 10)))
	return
}

func EditDict(w http.ResponseWriter, r *http.Request) {
	username := baseinfo.GetUsernameFromHeader(r)
	if username == "" {
		log.Println("user name is empty to add dict")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dict := new(model.Dict)
	log.Println(dict)
	log.Println("edit dict body...", r.Body)
	err := json.NewDecoder(r.Body).Decode(dict)
	log.Println(err)
	log.Println("json dict decoded...", dict)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("dict id", dict.Id)
	if dict.Id == 0 {
		http.Error(w, "dict id is null", http.StatusBadRequest)
		log.Println("dictionary id is null for update...")
		return
	}

	ret, succ := MysqlClient.UpdateDict(dict.Id, username, dict.Key, dict.Desc, dict.Type)
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

func DelDict(w http.ResponseWriter, r *http.Request) {
	var dictId, _ = strconv.Atoi(mux.Vars(r)["dictId"])
	log.Println(dictId)

	// ret, succ := MysqlClient.DelDictByIdReal(int64(dictId))
	ret, succ := MysqlClient.DelDictById(int64(dictId))
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

func DelDicts(w http.ResponseWriter, r *http.Request) {
	data := make(map[string][]int64)
	log.Println("ids...", data)
	log.Println("delete dicts body...", r.Body)
	err := json.NewDecoder(r.Body).Decode(&data)
	log.Println(err)
	log.Println("json dict ids decoded...", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if ids, ok := data["ids"]; ok {
		ret, succ := MysqlClient.DelDicts(ids)
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
