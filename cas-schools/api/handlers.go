package api

import (
	"cas-schools/model"
	"cas-schools/service"
	// // "baseinfo"
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

var MysqlClient service.ISchoolsClient

type Sizer interface {
	Size() int64
}

func GetSchool(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method)
	var schoolId, _ = strconv.Atoi(mux.Vars(r)["schoolId"])
	log.Println(schoolId)

	data := MysqlClient.GetSchoolById(int64(schoolId))
	ret, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(ret)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ret))

	return
}

func ListSchools(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	log.Println("queries...", queries)

	page, _ := strconv.Atoi(queries.Get("page"))
	items, _ := strconv.Atoi(queries.Get("items"))

	log.Println(page, items, "page and items...")
	// if page is 0, then return all
	data := MysqlClient.GetAllSchools(page, items)

	ret, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(ret)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ret))
	return
}

func AddSchool(w http.ResponseWriter, r *http.Request) {
	school := new(model.School)
	log.Println(school)
	log.Println("insert school body...", r.Body)
	err := json.NewDecoder(r.Body).Decode(school)
	log.Println(err)
	log.Println("json school decoded...", school)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = school.CheckSchool(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("school...", school)

	// name, contact, phone, province, city, county, provinceCode, cityCode, countyCode, addr, fax, email, web, post, from, to, contractId, contract string, isPayment, isLock bool, teacherNo, studentNo int
	ret, succ := MysqlClient.InsertSchool(school.Name, school.Contact, school.Phone, school.Province, school.City, school.County,
		school.ProvinceCode, school.CityCode, school.CountyCode, school.Address, school.Fax, school.Email, school.Web,
		school.Post, school.StartTime, school.ExpireTime, school.ContractId, school.Contract, school.IsPayment,
		school.IsLock, school.TeacherNo, school.StudentNo)

	adId, _ := ret.LastInsertId()
	log.Println("insert school ret...", ret, succ)
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

func EditSchool(w http.ResponseWriter, r *http.Request) {
	school := new(model.School)
	log.Println(school)
	log.Println("insert school body...", r.Body)
	err := json.NewDecoder(r.Body).Decode(school)
	log.Println(err)
	log.Println("json school decoded...", school)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = school.CheckSchool(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("school...", school)

	if school.Id == 0 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("school id is null for update...")
		return
	}

	ret, succ := MysqlClient.UpdateSchool(school.Id, school.Name, school.Contact, school.Phone, school.Province, school.City, school.County,
		school.ProvinceCode, school.CityCode, school.CountyCode, school.Address, school.Fax, school.Email, school.Web,
		school.Post, school.StartTime, school.ExpireTime, school.ContractId, school.Contract, school.IsPayment,
		school.IsLock, school.TeacherNo, school.StudentNo)
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

func DelSchool(w http.ResponseWriter, r *http.Request) {
	var schoolId, _ = strconv.Atoi(mux.Vars(r)["schoolId"])
	log.Println(schoolId)

	ret, succ := MysqlClient.DelSchoolById(int64(schoolId))
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
