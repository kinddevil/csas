package webserver

import (
	"log"
	"net/http"
)

func StartWebServer(port string, routes *Routes) {
	r := NewRouter(routes)
	http.Handle("/", r)
	log.Println("Starting HTTP service at " + port)
	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Println(err.Error())
	}
}
