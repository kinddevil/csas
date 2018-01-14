package api

import (
	"webserver"
)

// Initialize our routes
var Routes = webserver.Routes{

	webserver.Route{
		"GetRole",        // Name
		"GET",            // HTTP method
		"/role/{roleId}", // Route pattern
		GetRole,
	},

	webserver.Route{
		"GetRoles", // Name
		"GET",      // HTTP method
		"/role",    // Route pattern
		GetRoles,
	},

	webserver.Route{
		"AddRole", // Name
		"POST",    // HTTP method
		"/role",   // Route pattern
		AddRole,
	},

	webserver.Route{
		"EditRole", // Name
		"PUT",      // HTTP method
		"/role",    // Route pattern
		EditRole,
	},

	webserver.Route{
		"DelRoles", // Name
		"DELETE",   // HTTP method
		"/role",    // Route pattern
		DelRoles,
	},
}
