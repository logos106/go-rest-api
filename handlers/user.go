package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	db "github.com/saroopmathur/rest-api/db"
	model "github.com/saroopmathur/rest-api/models"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Add User ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)
	data, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewReader(data))
	log.Printf("%s\n", data)

	var user model.User
	var resp *model.User2
	err := decodeJSONBody(w, r, &user)
	if err == nil {
		resp, err = CreateUser1(r, &user)
	}
	httpSendResponse(w, 0, resp, err)
}

func CreateUser1(r *http.Request, user *model.User) (*model.User2, error) {
	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err := fmt.Errorf("domain %s %d unknown", domainName, domainId)
		return nil, err
	}

	return db.InsertUser(domainId, user)
}

// ReadUsers is an httpHandler for route GET /users
func ReadUsers(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Get All Users ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp []*model.User2

	offset, limit := reqPageInfo(r)

	search := reqSearchString(r)

	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		resp = db.SelectUsers(domainId, offset, limit, search)
	}

	httpSendResponse(w, 0, resp, err)
}

// ReadUser is an httpHandler for route GET /users/{id}
func ReadUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Get User By Id ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp *model.User2

	userName, userId := reqNameOrId(r)
	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		resp = db.SelectUser(domainId, userName, userId)
		fmt.Printf("ReadUser %s %d Domain %s %v\n", userName, userId, domainName, resp)
	}
	httpSendResponse(w, 0, resp, err)
}

// ReadUserGroups is an httpHandler for route GET /users/groups/{id}
func ReadUserGroups(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Get User Groups ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp []*model.Group2

	domainName, domainId := reqDomain(r)
	userName, userId := reqNameOrId(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		resp = db.GetUserGroups(domainId, userName, userId)
		fmt.Printf("ReadUserGroups User %s %d Domain %s %v\n", userName, userId, domainName, resp)
	}

	httpSendResponse(w, 0, resp, err)
}

// UpdateUser is an httpHandler for route PUT /users
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Update User ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)
	data, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewReader(data))
	log.Printf("%s\n", data)

	var err error
	var resp *model.User2

	userName, userId := reqNameOrId(r)
	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		// Decode the request body
		var user model.User
		err = decodeJSONBody(w, r, &user)
		if err == nil {
			resp = db.UpdateUser(domainId, userName, userId, &user)
		}
		fmt.Printf("Update User %s %d Domain %s %v\n", userName, userId, domainName, resp)
	}

	httpSendResponse(w, 0, resp, err)
}

// DeleteUser is an httpHandler for route DELETE /user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Delete User By Id ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp *model.User2

	userName, userId := reqNameOrId(r)
	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		resp := db.DeleteUser(domainId, userName, userId)
		if resp == nil {
			err = fmt.Errorf("unknown User")
		}
		// resp is the user obejct for the deleted user
		fmt.Printf("Delete User %s %d Domain %s %v\n", userName, userId, domainName, resp)
	}
	httpSendResponse(w, 0, resp, err)
}
