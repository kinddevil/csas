package api

import (
	"webserver"
)

const (
	Prefix string = "/advertising"
)

// Initialize our routes
var Routes = webserver.Routes{

	webserver.Route{
		"GetAccount", // Name
		"GET",        // HTTP method
		Prefix + "/accounts/{accountId}", // Route pattern
		GetAccount,
	},

	webserver.Route{
		"GetAdById", // Name
		"GET",       // HTTP method
		Prefix + "/ad/{adId}", // Route pattern
		GetAdById,
	},

	webserver.Route{
		"GetAds",        // Name
		"GET",           // HTTP method
		Prefix + "/ads", // Route pattern
		GetAds,
	},

	webserver.Route{
		"InsertAd",     // Name
		"POST",         // HTTP method
		Prefix + "/ad", // Route pattern
		InsertAd,
	},

	webserver.Route{
		"InsertAd",     // Name
		"PUT",          // HTTP method
		Prefix + "/ad", // Route pattern
		UpdateAd,
	},

	webserver.Route{
		"AdsUpload", // Name
		"POST",      // HTTP method
		Prefix + "/ad/{adId}/upload", // Route pattern
		UploadAd,
	},
}
