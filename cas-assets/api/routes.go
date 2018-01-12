package api

import (
	"webserver"
)

// Initialize our routes
var Routes = webserver.Routes{

	webserver.Route{
		"UploadAssets", // Name
		"POST",         // HTTP method
		"/upload",      // Route pattern
		UploadAssetsToOss,
	},
}
