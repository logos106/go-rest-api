package handler

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/saroopmathur/rest-api/db"
	model "github.com/saroopmathur/rest-api/models"
)

type AccessResp struct {
	RowsAffected int `json:"rows_affected,omitempty"`
}

// Access Control

// "UserAccess", "GET", "/users/access"
// List allowed applications for all users
func UserAccessAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== User Get Access for All ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp *[]model.UserAccess2
	domainName, domainId := reqDomain(r)
	if domainId == 0 { // Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		log.Printf("UserAccessAll: Domain:[%s %d]\n", domainName, domainId)
		resp = db.SelectUserAccessAll()
	}

	httpSendResponse(w, 0, resp, err)
}

// "UserAccessAll", "GET", "/users/access/{id}"
// List allowed applications for user {id}
func UserAccess(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== User Get Access ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp *[]model.App
	domainName, domainId := reqDomain(r)
	userName, userId := reqNameOrId(r)
	if domainId == 0 { // Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		log.Printf("UserAccess: Domain:[%s %d]\n", domainName, domainId)
		resp = db.SelectUserAccess(domainId, userName, userId)
	}

	httpSendResponse(w, 0, resp, err)
}

// "UserAddAccess", "POST", "/users/access/{id}/{id2}"
// Allow application {id2} to be accessible by user {id}
func UserAddAccess(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== User Add Access ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp *model.UserAccess
	domainName, domainId := reqDomain(r)
	userName, userId := reqNameOrId(r)
	appName, appId := reqNameOrId2(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		log.Printf("UserAddAccess: Domain:[%s %d] User:[%s %d] App:[%s %d]\n",
			domainName, domainId, userName, userId, appName, appId)
		resp, err = db.InsertUac(domainId, userName, userId, appName, appId)
	}

	httpSendResponse(w, 0, resp, err)
}

// "UserDelAccess", "DELETE", "/users/access/{id}/{id2}"
// Remove access to application {id2} for user {id}
func UserDelAccess(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== User Del Access ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var count int
	var resp *AccessResp

	domainName, domainId := reqDomain(r)
	userName, userId := reqNameOrId(r)
	appName, appId := reqNameOrId2(r)
	if domainId == 0 { // Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		log.Printf("UserDelAccess: Domain:[%s %d] User:[%s %d] App:[%s %d]\n",
			domainName, domainId, userName, userId, appName, appId)

		count = db.DeleteUac(domainId, userName, userId, appName, appId)
	}

	if err == nil {
		resp = &AccessResp{}
		resp.RowsAffected = count
	}

	httpSendResponse(w, 0, resp, err)
}

///////////////////////////////////////////////////////////

// "GroupAccessAll", "GET", "/groups/access/{id}"
// List allowed applications for all groups
func GroupAccessAll(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Group Access List All ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp *[]model.GroupAccess2
	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		log.Printf("SelectUserAccessAll: Domain:[%s %d]\n", domainName, domainId)
		resp = db.SelectGroupAccessAll()
	}

	httpSendResponse(w, 0, resp, err)
}

// "GroupAccess", "GET", "/groups/access/{id}"
// List allowed applications for group {id}
func GroupAccess(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== User Access List ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp *[]model.App
	domainName, domainId := reqDomain(r)
	groupName, grouprId := reqNameOrId(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		log.Printf("SelectGroupAccess: Domain:[%s %d]\n", domainName, domainId)
		resp = db.SelectGroupAccess(domainId, groupName, grouprId)
	}

	httpSendResponse(w, 0, resp, err)
}

// "GroupAddAccess", "POST", "/groups/access/{id}/{id2}"
// Allow application {id2} to be accessible by group {id}
func GroupAddAccess(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Group Add Access ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp *model.GroupAccess
	domainName, domainId := reqDomain(r)
	groupName, groupId := reqNameOrId(r)
	appName, appId := reqNameOrId2(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		log.Printf("GroupAddAccess: Domain:[%s %d] User:[%s %d] App:[%s %d]\n",
			domainName, domainId, groupName, groupId, appName, appId)
		resp, err = db.InsertGac(domainId, groupName, groupId, appName, appId)
	}

	httpSendResponse(w, 0, resp, err)
}

// "GroupDelAccess", "DELETE", "/groups/access/{id}/{id2}"
// Remove access to application {id2} for group {id}
func GroupDelAccess(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Group Del Access ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var count int
	var resp *AccessResp

	domainName, domainId := reqDomain(r)
	groupName, groupId := reqNameOrId(r)
	appName, appId := reqNameOrId2(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		log.Printf("GroupDelAccess: Domain:[%s %d] User:[%s %d] App:[%s %d]\n",
			domainName, domainId, groupName, groupId, appName, appId)
		count = db.DeleteGac(domainId, groupName, groupId, appName, appId)
	}

	if err == nil {
		resp = &AccessResp{}
		resp.RowsAffected = count
	}

	httpSendResponse(w, 0, resp, err)
}
