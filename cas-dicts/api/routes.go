package api

import (
	"webserver"
)

// Initialize our routes
var Routes = webserver.Routes{

	webserver.Route{
		"GetDict", // Name
		"GET",     // HTTP method
		"/dictionary/{dictId}", // Route pattern
		GetDict,
	},

	webserver.Route{
		"GetDicts",    // Name
		"GET",         // HTTP method
		"/dictionary", // Route pattern
		GetDicts,
	},

	webserver.Route{
		"AddDict",     // Name
		"POST",        // HTTP method
		"/dictionary", // Route pattern
		AddDict,
	},

	webserver.Route{
		"EditDict",    // Name
		"PUT",         // HTTP method
		"/dictionary", // Route pattern
		EditDict,
	},

	webserver.Route{
		"DelDict",              // Name
		"DELETE",               // HTTP method
		"/dictionary/{dictId}", // Route pattern
		DelDict,
	},
}
