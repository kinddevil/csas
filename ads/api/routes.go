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
		"GetAdById",  // Name
		"GET",        // HTTP method
		"/ad/{adId}", // Route pattern
		GetAdById,
	},

	webserver.Route{
		"GetAds", // Name
		"GET",    // HTTP method
		"/ads",   // Route pattern
		GetAds,
	},

	webserver.Route{
		"InsertAd", // Name
		"POST",     // HTTP method
		"/ad",      // Route pattern
		InsertAd,
	},

	webserver.Route{
		"UpdateAd", // Name
		"PUT",      // HTTP method
		"/ad",      // Route pattern
		UpdateAd,
	},

	webserver.Route{
		"AdsUpload",         // Name
		"POST",              // HTTP method
		"/ad/{adId}/upload", // Route pattern
		UploadAd,
	},
}
