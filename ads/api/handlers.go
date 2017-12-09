package api

import (
	"ads/model"
	"ads/service"
	"dbclient"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	// "html/template"
	// "net"
)

var DBClient dbclient.IBoltClient
var MysqlClient service.IAdsClient

type Sizer interface {
	Size() int64
}

func GetAccount(w http.ResponseWriter, r *http.Request) {

	log.Println(r.Method)
	// Read the 'accountId' path parameter from the mux map
	var accountId = mux.Vars(r)["accountId"]

	// Read the account struct BoltDB
	account, err := DBClient.QueryAccount(accountId)

	// If err, return a 404
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// If found, marshal into JSON, write headers and content
	data, _ := json.Marshal(account)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func InsertAd(w http.ResponseWriter, r *http.Request) {
	ad := new(model.Advertising)
	log.Println(ad)
	log.Println(r.Body)
	err := json.NewDecoder(r.Body).Decode(&ad)
	log.Println(err)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(ad)
	ret := MysqlClient.InsertAd(ad.ImageName, ad.ImageLink, ad.SchoolIds, ad.Province, ad.City, ad.Title, ad.DisplayPages)
	fmt.Println("insert ret...", ret)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatBool(ret)))
}

func UpdateAd(w http.ResponseWriter, r *http.Request) {
	ret := MysqlClient.UpdateAd(1, "test", "test", "1,2,3", "bj", "bj", "cahnge test ad", 1)
	fmt.Println("update ret...", ret)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatBool(ret)))
}

func GetAdById(w http.ResponseWriter, r *http.Request) {
	var adId, _ = strconv.Atoi(mux.Vars(r)["adId"])
	data := MysqlClient.GetAdById(int64(adId))

	ret, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Content-Length", strconv.Itoa(12))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ret))
}

func GetAds(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	log.Println("queries...", queries)
	page, _ := strconv.Atoi(queries.Get("page"))
	items, _ := strconv.Atoi(queries.Get("items"))
	// if page is 0, then return all
	data := MysqlClient.GetAllAds(page, items)

	ret, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Content-Length", strconv.Itoa(12))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ret))
}

func UploadAd(w http.ResponseWriter, r *http.Request) {
	var adId, _ = strconv.Atoi(mux.Vars(r)["adId"])
	log.Println(adId)

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		log.Println("parse body error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := r.MultipartForm
	files := m.File["data"]

	log.Println("multipar form...", m.Value)
	log.Println(r.Body)
	for i := 0; i < len(files); i++ {
		log.Println("upload file...", files[i])
		file := files[i]
		fileName := strings.ToLower(file.Filename)
		log.Println(fileName)

		if match, _ := regexp.MatchString("\\.(png|jpeg|jpg)$", fileName); !match {
			// if !(strings.HasSuffix(fileName, ".png") || strings.HasSuffix(fileName, ".jpeg") || strings.HasSuffix(fileName, ".jpg")) {
			ferr := "file format error, not png, jpeg or jpg"
			log.Println(ferr)
			http.Error(w, ferr, http.StatusInternalServerError)
			return
		}

		log.Println(file.Header)
		log.Println(reflect.TypeOf(file))

		infile, _ := file.Open()
		log.Println(infile.(Sizer).Size())
		defer infile.Close()

		fmt.Fprintf(w, "%v", file.Header)
		path := "./assets/adimg/"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Println("path does not exist, create path", path)
			os.MkdirAll(path, os.ModePerm)
		}
		f, err := os.OpenFile(path+file.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println("open file error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		io.Copy(f, infile)

		//TODO: support multiple upload with checks
		break
	}

	w.Header().Set("Content-Type", "application/json")
	return
}

func UploadAds(w http.ResponseWriter, r *http.Request) {

	log.Println(r.Method)
	// http://www.giantflyingsaucer.com/blog/?p=5635

	// login := r.FormValue("login")
	// dec := json.NewDecoder(req.Body)

	// body, err := ioutil.ReadAll(req.Body)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(string(body))
	// var t test_struct
	// err = json.Unmarshal(body, &t)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(t.Test)

	// err := json.NewDecoder(r.Body).Decode(&u)
	// if err != nil {
	// 	http.Error(w, err.Error(), 400)
	// 	return
	// }
	// fmt.Println(u.Id)

	//router.HandleFunc("/movie/{imdbKey}", handleMovie).Methods("GET", "DELETE")
	// switch req.Method {
	// case "GET":
	// 	outgoingJSON, error := json.Marshal(movie)
	// 	if error != nil {
	// 		log.Println(error.Error())
	// 		http.Error(res, error.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	fmt.Fprint(res, string(outgoingJSON))
	// case "DELETE":
	// 	delete(movies, imdbKey)
	// 	res.WriteHeader(http.StatusNoContent)
	// }

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		log.Println("parse body error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := r.MultipartForm
	files := m.File["data"]

	fmt.Println(files)
	fmt.Println(m.Value)
	for i := 0; i < len(files); i++ {
		fmt.Println(files[i])
		file := files[i]
		fmt.Println(file.Filename)
		fmt.Println(file.Header)
		fmt.Println(reflect.TypeOf(file))
	}

	file := files[0]
	fmt.Println("File...")
	fmt.Println(file)
	infile, _ := file.Open()
	fmt.Println(infile.(Sizer).Size())
	defer infile.Close()

	// file, handler, err := r.FormFile("uploadfile")
	// if err != nil {
	// 	log.Println("get file error", err)
	// 	return
	// }
	// defer file.Close()

	// fmt.Fprintf(w, "%v", handler.Header)
	// f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)

	fmt.Fprintf(w, "%v", file.Header)
	f, err := os.OpenFile("./test/"+file.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("open file error", err)
		return
	}
	defer f.Close()
	io.Copy(f, infile)

	// ret := make(map[string]string)
	// ret["status"] = "ok"
	// data, _ := json.Marshal(ret)
	// log.Println(data)
	// log.Println(r.Method)
	// log.Println("uploads...")
	w.Header().Set("Content-Type", "application/json")
	// panic("pnc")
	// w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	// w.WriteHeader(http.StatusOK)
	// log.Println(http.StatusBadGateway)
	// w.Write(data)
	return
}
