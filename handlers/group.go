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

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Add Group ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)
	data, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewReader(data))
	log.Printf("%s\n", data)

	var group model.Group
	var resp *model.Group2
	err := decodeJSONBody(w, r, &group)
	if err == nil {
		resp, err = CreateGroup1(r, &group)
	}
	httpSendResponse(w, 0, resp, err)
}

func CreateGroup1(r *http.Request, group *model.Group) (*model.Group2, error) {
	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err := fmt.Errorf("domain %s %d unknown", domainName, domainId)
		return nil, err
	}

	return db.InsertGroup(domainId, group)
}

// ReadGroups is an httpHandler for route GET /groups
func ReadGroups(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Get All Groups ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp []*model.Group2

	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		resp = db.SelectGroups(domainId)
	}

	httpSendResponse(w, 0, resp, err)
}

// ReadGroup is an httpHandler for route GET /groups/{id}
func ReadGroup(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Get a Group By Id ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp *model.Group2

	groupName, groupId := reqNameOrId(r)
	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		resp = db.SelectGroup(domainId, groupName, groupId)
	}
	httpSendResponse(w, 0, resp, err)
}

// ReadGroupUsers is an httpHandler for route GET /groups/users/{id}
func ReadGroupUsers(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Read Group Users ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp []*model.User2

	groupName, groupId := reqNameOrId(r)
	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		resp = db.GetGroupUsers(domainId, groupName, groupId)
	}
	httpSendResponse(w, 0, resp, err)
}

// UpdateGroup is an httpHandler for route PUT /groups
func UpdateGroup(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Update Group ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)
	data, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewReader(data))
	log.Printf("%s\n", data)

	var err error
	var resp *model.Group2

	groupName, groupId := reqNameOrId(r)
	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		// Decode the request body
		var group model.Group
		err = decodeJSONBody(w, r, &group)
		if err == nil {
			resp = db.UpdateGroup(domainId, groupName, groupId, &group)
		}
	}

	httpSendResponse(w, 0, resp, err)
}

// DeleteGroup is an httpHandler for route DELETE /group
func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Delete Group ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp *model.Group2

	groupName, groupId := reqNameOrId(r)
	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		resp, rowsAffected := db.DeleteGroup(domainId, groupName, groupId)
		if resp == nil {
			err = fmt.Errorf("unknown Group")
		} else {
			fmt.Printf("Delete Group %s %d domain %s, deleted %d members\n",
				groupName, groupId, domainName, rowsAffected)
		}
		// resp is the group obejct for the deleted group
	}
	httpSendResponse(w, 0, resp, err)
}
