package api

import (
	"webserver"
)

// Initialize our routes
var Routes = webserver.Routes{

	webserver.Route{
		"GetCourse",     // Name
		"GET",             // HTTP method
		"/course/{cid}", // Route pattern
		GetCourse,
	},

	webserver.Route{
		"GetCourses", // Name
		"GET",          // HTTP method
		"/course",    // Route pattern
		GetCourses,
	},

	webserver.Route{
		"AddCourse", // Name
		"POST",        // HTTP method
		"/course",   // Route pattern
		AddCourse,
	},

	webserver.Route{
		"EditCourse", // Name
		"PUT",          // HTTP method
		"/course",    // Route pattern
		EditCourse,
	},

	webserver.Route{
		"DelCourse",     // Name
		"DELETE",          // HTTP method
		"/course/{cid}", // Route pattern
		DelCourse,
	},

	webserver.Route{
		"DelCourses", // Name
		"DELETE",       // HTTP method
		"/course",    // Route pattern
		DelCourses,
	},
}
