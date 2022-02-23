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

type MemberResp struct {
	RowsAffected int `json:"rows_affected,omitempty"`
}

func AddGroupMembers(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Add Users to Group ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)
	data, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewReader(data))
	log.Printf("%s\n", data)

	var users []string
	var err error
	var addCount int
	var resp *MemberResp

	domainName, domainId := reqDomain(r)
	groupName, groupId := reqNameOrId(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		err = decodeJSONBody(w, r, &users)
		// Remove duplicates
		users1 := RemoveDuplicateValues(users)

		fmt.Printf("AddGroupMembers{%s %d %s %d] %v\n", groupName, groupId, domainName, domainId, users1)
		if err == nil {
			addCount = db.AddGroupMembers(domainId, groupName, groupId, users1)
		}
	}

	if err == nil {
		resp = &MemberResp{}
		resp.RowsAffected = addCount
	}

	httpSendResponse(w, 0, resp, err)
}

func RemoveDuplicateValues(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}

	return list
}

// ReadGroupMembers is an httpHandler for route GET /groupmembers/{id}
func ReadGroupMembers(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Get Group Members ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp []*model.User2

	domainName, domainId := reqDomain(r)
	groupName, groupId := reqNameOrId(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		fmt.Printf("ReadGroupMembers{%s %d %d]\n", groupName, groupId, domainId)
		resp = db.SelectGroupMembers(domainId, groupName, groupId)
	}
	httpSendResponse(w, 0, resp, err)
}

// RemoveGroupMember is an httpHandler for route POST /groupmembers/remove/{id}
func RemoveGroupMembers(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Remove Users from Group ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var count int
	var resp *MemberResp

	domainName, domainId := reqDomain(r)
	groupName, groupId := reqNameOrId(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		var users []string
		err = decodeJSONBody(w, r, &users)
		fmt.Printf("DeleteGroupMembers{%s %d %d] %v\n", groupName, groupId, domainId, users)
		if err == nil {
			count = db.RemoveGroupMembers(domainId, groupName, groupId, users)
		}
	}

	if err == nil {
		resp = &MemberResp{}
		resp.RowsAffected = count
	}

	httpSendResponse(w, 0, resp, err)
}
