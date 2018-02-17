package api

import (
	"webserver"
)

// Initialize our routes
var Routes = webserver.Routes{

	webserver.Route{
		"GetUser",          // Name
		"GET",              // HTTP method
		"/user/{username}", // Route pattern
		GetUser,
	},

	webserver.Route{
		"GetUsers", // Name
		"GET",      // HTTP method
		"/user",    // Route pattern
		GetUsers,
	},

	webserver.Route{
		"AddUser", // Name
		"POST",    // HTTP method
		"/user",   // Route pattern
		AddUser,
	},

	webserver.Route{
		"ResetPwd",  // Name
		"POST",      // HTTP method
		"/resetpwd", // Route pattern
		ResetPwd,
	},

	webserver.Route{
		"UpdatePwd",  // Name
		"POST",       // HTTP method
		"/updatepwd", // Route pattern
		UpdatePwd,
	},

	webserver.Route{
		"EditUser", // Name
		"PUT",      // HTTP method
		"/user",    // Route pattern
		EditUser,
	},

	webserver.Route{
		"DelUsers", // Name
		"DELETE",   // HTTP method
		"/user",    // Route pattern
		DelUsers,
	},
}
