package router

import (
	"net/http"

	handler "github.com/saroopmathur/rest-api/handlers"
)

// Route type description
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes contains all routes
type Routes []Route

// For auth and others
var routes0 = Routes{
	Route{
		"Login",
		"GET",
		"/login",
		handler.Login,
	},
	Route{
		"Login",
		"GET",
		"/servicelogin",
		handler.Login,
	},
	Route{
		"Login",
		"GET",
		"/adminlogin",
		handler.Login,
	},
	Route{
		"Login",
		"GET",
		"/logout",
		handler.Logout,
	},
}

// For domain
var routes1 = Routes{
	Route{
		"CreateDomain",
		"POST",
		"/domains",
		handler.CreateDomain,
	},
	Route{
		"ReadDomains",
		"GET",
		"/domains",
		handler.ReadDomains,
	},
	Route{
		"ReadDomain",
		"GET",
		"/domains/{id}",
		handler.ReadDomain,
	},
	Route{
		"UpdateDomain",
		"PUT",
		"/domains/{id}",
		handler.UpdateDomain,
	},
	Route{
		"DeleteDomain",
		"DELETE",
		"/domains/{id}",
		handler.DeleteDomain,
	},
	Route{
		"ChangeDomain",
		"POST",
		"/changedomain/{id}",
		handler.ChangeDomain,
	},
}

// For admin
var routes2 = Routes{
	Route{
		"CreateAdmin",
		"POST",
		"/admins",
		handler.CreateAdmin,
	},
	Route{
		"ReadAdmins",
		"GET",
		"/admins",
		handler.ReadAdmins,
	},
	Route{
		"ReadAdmin",
		"GET",
		"/admins/{id}",
		handler.ReadAdmin,
	},
	Route{
		"UpdateAdmin",
		"PUT",
		"/admins/{id}",
		handler.UpdateAdmin,
	},
	Route{
		"DeleteAdmin",
		"DELETE",
		"/admins/{id}",
		handler.DeleteAdmin,
	},
}

// For user
var routes3 = Routes{
	Route{
		"CreateUser",
		"POST",
		"/users",
		handler.CreateUser,
	},
	Route{
		"ReadUsers",
		"GET",
		"/users",
		handler.ReadUsers,
	},
	Route{
		"UserAccessAll",
		"GET",
		"/users/access",
		handler.UserAccessAll,
	},
	Route{
		"ReadUser",
		"GET",
		"/users/{id}",
		handler.ReadUser,
	},
	Route{
		"UpdateUser",
		"PUT",
		"/users/{id}",
		handler.UpdateUser,
	},
	Route{
		"DeleteUser",
		"DELETE",
		"/users/{id}",
		handler.DeleteUser,
	},
	Route{
		"ReadUserGroups",
		"GET",
		"/users/groups/{id}",
		handler.ReadUserGroups,
	},
}

// For service
var routes4 = Routes{
	Route{
		"CreateService",
		"POST",
		"/services",
		handler.CreateService,
	},
	Route{
		"ReadServices",
		"GET",
		"/services",
		handler.ReadServices,
	},
	Route{
		"ReadService",
		"GET",
		"/services/{id}",
		handler.ReadService,
	},
	Route{
		"UpdateService",
		"PUT",
		"/services/{id}",
		handler.UpdateService,
	},
	Route{
		"DeleteService",
		"DELETE",
		"/services/{id}",
		handler.DeleteService,
	},
}

// For allowed app
var routes5 = Routes{
	Route{
		"CreateApp",
		"POST",
		"/apps",
		handler.CreateApp,
	},
	Route{
		"ReadApp",
		"GET",
		"/apps/{id}",
		handler.ReadApp,
	},
	Route{
		"ReadApp",
		"GET",
		"/apps2/{id}",
		handler.ReadApp2,
	},
	Route{
		"ReadApps",
		"GET",
		"/apps",
		handler.ReadApps,
	},
	Route{
		"UpdateApp",
		"PUT",
		"/apps/{id}",
		handler.UpdateApp,
	},
	Route{
		"DeleteApps",
		"DELETE",
		"/apps/{id}",
		handler.DeleteApp,
	},
	// Policies
	Route{
		"GetPolicy",
		"GET",
		"/policies/{id}",
		handler.GetPolicy,
	},
	Route{
		"GetPolicies",
		"GET",
		"/policies",
		handler.GetPolicies,
	},

	// Access Control
	Route{
		"UserAccess",
		"GET",
		"/users/access/{id}",
		handler.UserAccess,
	},
	Route{
		"UserAddAccess",
		"POST",
		"/users/access/{id}/{id2}",
		handler.UserAddAccess,
	},
	Route{
		"UserDelAccess",
		"DELETE",
		"/users/access/{id}/{id2}",
		handler.UserDelAccess,
	},
	Route{
		"GroupAccess",
		"GET",
		"/groups/access/{id}",
		handler.GroupAccess,
	},
	Route{
		"GroupAddAccess",
		"POST",
		"/groups/access/{id}/{id2}",
		handler.GroupAddAccess,
	},
	Route{
		"GroupDelAccess",
		"DELETE",
		"/groups/access/{id}/{id2}",
		handler.GroupDelAccess,
	},
}

// For user group
var routes6 = Routes{
	Route{
		"CreateGroup",
		"POST",
		"/groups",
		handler.CreateGroup,
	},
	Route{
		"ReadGroups",
		"GET",
		"/groups",
		handler.ReadGroups,
	},
	Route{
		"GroupAccessAll",
		"GET",
		"/groups/access",
		handler.GroupAccessAll,
	},
	Route{
		"ReadGroup",
		"GET",
		"/groups/{id}",
		handler.ReadGroup,
	},
	Route{
		"UpdateGroup",
		"PUT",
		"/groups/{id}",
		handler.UpdateGroup,
	},
	Route{
		"DeleteGroup",
		"DELETE",
		"/groups/{id}",
		handler.DeleteGroup,
	},
	Route{
		"ReadGroupUsers",
		"GET",
		"/groups/users/{id}",
		handler.ReadGroupUsers,
	},
}

// For group memeber
var routes7 = Routes{
	Route{
		"AddGroupMembers",
		"POST",
		"/groupmembers/add/{id}",
		handler.AddGroupMembers,
	},
	Route{
		"RemoveGroupMember",
		"POST",
		"/groupmembers/remove/{id}",
		handler.RemoveGroupMembers,
	},
	Route{
		"ReadGroupMembers",
		"GET",
		"/groupmembers/{id}",
		handler.ReadGroupMembers,
	},
}
