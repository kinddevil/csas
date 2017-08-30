package service

import (
	"webserver"
)

// Initialize our routes
var Routes = webserver.Routes{

	webserver.Route{
		"GetAccount", // Name
		"GET",        // HTTP method
		"/accounts/{accountId}", // Route pattern
		// func(w http.ResponseWriter, r *http.Request) {
		//    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		//    w.Write([]byte("{\"result\":\"OK\"}"))
		//},
		GetAccount,
	},

	webserver.Route{
		"AdsUpload",   // Name
		"POST",        // HTTP method
		"/ads/upload", // Route pattern
		// func(w http.ResponseWriter, r *http.Request) {
		//    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		//    w.Write([]byte("{\"result\":\"OK\"}"))
		//},
		UploadAds,
	},
}
