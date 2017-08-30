package service

import (
	"dbclient"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	// "html/template"
	// "net"
)

var DBClient dbclient.IBoltClient
var MysqlClient dbclient.IMysqlClient

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

func UploadAds(w http.ResponseWriter, r *http.Request) {

	log.Println(r.Method)

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
