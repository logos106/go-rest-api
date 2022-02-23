package handler

import (
	"fmt"
	"net/http"

	db "github.com/saroopmathur/rest-api/db"
	m "github.com/saroopmathur/rest-api/models"
)

// "GetPolicies", "GET", "/policies",
// Get all policies
// Called by either Admin from web app or User from client app
func GetPolicies(w http.ResponseWriter, r *http.Request) {
	var err error
	var policy *m.Policy
	u := reqUser(r)
	if u.Domain.ID == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain Unknown")
	} else {
		switch u.Role {
		case db.ROLE_USER:
			policy, err = db.GetUserPolicy(u.Domain.ID, u.Name, u.ID)
		case db.ROLE_ADMIN, db.ROLE_POWERADMIN:
			domainName, domainId := reqDomain(r)
			policy, err = db.GetAllPolicies(domainId)
			fmt.Printf("GetAllPolicies: domain [%s %d] %v\n", domainName, domainId, policy)
		case db.ROLE_SERVICE:
			//policy, err = db.GetUserPolicy(u.Domain.ID, u.Name, u.ID)
			//fmt.Printf("GetPolicies: domain [%s %d] %v\n", domainName, domainId, policies)
		}
	}
	httpSendResponse(w, 0, policy, err)
}

// "GetPolicy", "GET", "/policies/{id}",
// Get policies for the specified username or userId (If caller is an Admin)
// Must be Admin to call this API
func GetPolicy(w http.ResponseWriter, r *http.Request) {
	var err error
	var policy *m.Policy

	u := reqUser(r)

	switch u.Role {
	case db.ROLE_ADMIN, db.ROLE_POWERADMIN:
		userName, userId := reqNameOrId(r)
		domainName, domainId := reqDomain(r)
		if domainId == 0 {
			// Unknown Domain
			err = fmt.Errorf("domain Unknown")
		} else {
			policy, err = db.GetUserPolicy(domainId, userName, userId)
			fmt.Printf("GetUserPolicy: User %s %d domain %s %v\n", userName, userId, domainName, policy)
		}
	case db.ROLE_USER:
	case db.ROLE_SERVICE:
	}
	httpSendResponse(w, 0, policy, err)
}
