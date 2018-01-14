package main

import (
	// "encoding/json"
	"fmt"
	"reflect"
	"strings"
	// "github.com/gorilla/mux"
	// "io"
	// "log"
	// "net/http"
	// "os"
	// "strconv"
	// "html/template"
	// "net"
)

type UserInfo struct {
	Ids   int
	Names string
	pri   string
}

func GetFieldMap(obj interface{}) (ret map[string]string) {
	val := reflect.ValueOf(obj).Elem()
	ret = make(map[string]string)
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		fmt.Println(typeField.Type)
		key := strings.ToLower(typeField.Name)
		if typeField.PkgPath != "" {
			// Private method
			continue
		} else {
			ret[key] = typeField.Name
		}
	}
	return
}

func main() {
	user := UserInfo{1, "name", "other"}
	ret := GetFieldMap(&user)
	fmt.Println(ret, len(ret))
}
