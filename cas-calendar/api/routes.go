package api

import (
	"webserver"
)

// Initialize our routes
var Routes = webserver.Routes{

	webserver.Route{
		"GetCalendar",     // Name
		"GET",             // HTTP method
		"/calendar/{cid}", // Route pattern
		GetCalendar,
	},

	webserver.Route{
		"GetCalendars", // Name
		"GET",          // HTTP method
		"/calendar",    // Route pattern
		GetCalendars,
	},

	webserver.Route{
		"AddCalendar", // Name
		"POST",        // HTTP method
		"/calendar",   // Route pattern
		AddCalendar,
	},

	webserver.Route{
		"EditCalendar", // Name
		"PUT",          // HTTP method
		"/calendar",    // Route pattern
		EditCalendar,
	},

	webserver.Route{
		"DelCalendar",     // Name
		"DELETE",          // HTTP method
		"/calendar/{cid}", // Route pattern
		DelCalendar,
	},

	webserver.Route{
		"DelCalendars", // Name
		"DELETE",       // HTTP method
		"/calendar",    // Route pattern
		DelCalendars,
	},

	webserver.Route{
		"Test",  // Name
		"GET",   // HTTP method
		"/test", // Route pattern
		Test,
	},
}
