package webserver

import (
	"log"
	"net/http"
)

func StartWebServer(port string, routes *Routes) {
	s := http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/")))
	r := NewRouter(routes)
	r.PathPrefix("/assets/").Handler(s)
	// r.PathPrefix("/assets").Handler(http.FileServer(http.Dir("./assets/")))
	http.Handle("/", r)
	log.Println("Starting HTTP service at " + port)
	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Println(err.Error())
	}
}
