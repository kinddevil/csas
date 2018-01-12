package api

import (
	// "cas-assets/model"
	"cas-assets/service"
	// // "baseinfo"
	// "dbclient"
	"cas-assets/oss"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"path/filepath"
	// "hash/fnv"
	// fnv.New32a() h.Sum32() https://play.golang.org/p/_J2YysdEqE
	"github.com/satori/go.uuid"
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

var MysqlClient service.IAssetsClient

const (
	file_size_limit = 512
)

type Sizer interface {
	Size() int64
}

func UploadAssetsToOss(w http.ResponseWriter, r *http.Request) {
	// username := baseinfo.GetUsernameFromHeader(r)
	// if username == "" {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	// sid, _, _ := MysqlClient.GetBaseInfo(username)
	// sidS := strconv.FormatInt(sid, 10)

	// Start upload files...

	ret := make([]map[string]string, 0, 9)

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		log.Println("parse body error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := r.MultipartForm
	files := m.File["file_upload_name"]
	var keys []string
	if val, ok := m.Value["file_upload_update"]; ok {
		keys = val
	}

	log.Println("multipar form...", m.Value, reflect.TypeOf(m.Value))
	log.Println("files and keys ...", files, keys)
	// log.Println(r.Body)

	for i := 0; i < len(files); i++ {
		// log.Println("upload file...", files[i])
		file := files[i]
		fileName := strings.ToLower(file.Filename)
		log.Println("filename:", fileName, "header:", file.Header)
		log.Println("type:", reflect.TypeOf(file))

		if match, _ := regexp.MatchString("\\.(png|jpeg|jpg|gif)$", fileName); !match {
			// if !(strings.HasSuffix(fileName, ".png") || strings.HasSuffix(fileName, ".jpeg") || strings.HasSuffix(fileName, ".jpg")) {
			ferr := "file format error, not png, jpeg or jpg"
			log.Println(ferr)
			http.Error(w, ferr, http.StatusInternalServerError)
			return
		}

		id := ""
		if len(keys) > i {
			id = keys[i]
		} else {
			id = uuid.NewV4().String()
		}
		log.Println("id...", id)

		infile, _ := file.Open()
		filesize := infile.(Sizer).Size() / 1000
		log.Println("size...", infile.(Sizer).Size(), filesize) // byte
		defer infile.Close()

		if filesize > file_size_limit {
			serr := fmt.Sprintf("file size exceed %dkb", file_size_limit)
			log.Println(serr)
			http.Error(w, serr, http.StatusInternalServerError)
			return
		}

		// fmt.Fprintf(w, "%v", file.Header)
		// err = bucket.PutObjectFromFile("assets/"+id, [path], options...)
		path := oss.ImgPrefix + id + filepath.Ext(fileName)
		if err := oss.Bucket.PutObject(oss.ImgPrefix+id+filepath.Ext(fileName), infile, oss.ImgOptions...); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// db operation
		_, succ := MysqlClient.InsertOrUpdateAssets(id, fileName, oss.Region, path, int(filesize))
		if !succ {
			dberr := "save error, please retry"
			log.Println(dberr)
			http.Error(w, dberr, http.StatusInternalServerError)
			return
		}

		res := make(map[string]string)
		res["filename"] = fileName
		res["id"] = id
		res["path"] = path
		ret = append(ret, res)

		//TODO: support multiple upload with checks
		break
	}
	// End upload files...

	rdata, _ := json.Marshal(ret)
	w.Header().Set("Content-Type", "application/json")
	w.Write(rdata)
	return
}

func UploadAssets(w http.ResponseWriter, r *http.Request) {
	// username := baseinfo.GetUsernameFromHeader(r)
	// if username == "" {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	// sid, _, _ := MysqlClient.GetBaseInfo(username)
	// sidS := strconv.FormatInt(sid, 10)
	fname := ""

	var adId, _ = strconv.Atoi(mux.Vars(r)["adId"])
	log.Println(adId)

	// Start upload files...
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		log.Println("parse body error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := r.MultipartForm
	files := m.File["file_upload_name"]
	// file_upload_name file_upload_update

	log.Println("multipar form...", m.Value)
	log.Println(r.Body)
	for i := 0; i < len(files); i++ {
		log.Println("upload file...", files[i])
		file := files[i]
		fileName := strings.ToLower(file.Filename)
		fname = fileName
		log.Println(fileName)

		if match, _ := regexp.MatchString("\\.(png|jpeg|jpg|gif)$", fileName); !match {
			// if !(strings.HasSuffix(fileName, ".png") || strings.HasSuffix(fileName, ".jpeg") || strings.HasSuffix(fileName, ".jpg")) {
			ferr := "file format error, not png, jpeg or jpg"
			log.Println(ferr)
			http.Error(w, ferr, http.StatusInternalServerError)
			return
		}

		log.Println(file.Header)
		log.Println(reflect.TypeOf(file))

		infile, _ := file.Open()
		log.Println("size...", infile.(Sizer).Size()) // byte
		defer infile.Close()

		fmt.Fprintf(w, "%v", file.Header)
		path := "./assets/adimg/"
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Println("path does not exist, create path", path)
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				log.Println("create dir error", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
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
	// End upload files...

	log.Println("filename...", fname)

	if adId <= 0 {
		log.Println("No ads for files, create new ad item")
		// ret, _ := MysqlClient.InsertEmptyAd()
		// aid, _ := ret.LastInsertId()
		// adId = int(aid)
	}

	// MysqlClient.SaveUploadFiles(adId, fname)

	// fileRet, _ := MysqlClient.SaveUploadFiles(adId, fname)
	// log.Println(fileRet.LastInsertId())

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(strconv.Itoa(adId)))
	return
}
