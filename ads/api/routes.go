package api

import (
	"webserver"
)

// Initialize our routes
var Routes = webserver.Routes{

	webserver.Route{
		"GetAccount", // Name
		"GET",        // HTTP method
		"/accounts/{accountId}", // Route pattern
		GetAccount,
	},

	webserver.Route{
		"GetAd",       // Name
		"GET",         // HTTP method
		"/ads/{adId}", // Route pattern
		GetAd,
	},

	webserver.Route{
		"AdsUpload",   // Name
		"POST",        // HTTP method
		"/ads/upload", // Route pattern
		UploadAds,
	},
}
