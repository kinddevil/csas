package api

import (
	"webserver"
)

// Initialize our routes
var Routes = webserver.Routes{

	webserver.Route{
		"GetSchool",          // Name
		"GET",                // HTTP method
		"/school/{schoolId}", // Route pattern
		GetSchool,
	},

	webserver.Route{
		"AddSchool", // Name
		"POST",      // HTTP method
		"/school",   // Route pattern
		AddSchool,
	},

	webserver.Route{
		"DelSchool",          // Name
		"DELETE",             // HTTP method
		"/school/{schoolId}", // Route pattern
		DelSchool,
	},

	webserver.Route{
		"EditSchool", // Name
		"PUT",        // HTTP method
		"/school",    // Route pattern
		EditSchool,
	},

	webserver.Route{
		"ListSchools", // Name
		"GET",         // HTTP method
		"/school",     // Route pattern
		ListSchools,
	},
}
