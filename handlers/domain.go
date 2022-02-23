package handler

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/saroopmathur/rest-api/db"
	model "github.com/saroopmathur/rest-api/models"
)

func CreateDomain(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Add Domain ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var domain model.DomainReq
	var resp *model.Domain
	var err error
	var code int
	if !reqIsSuperuser(r) {
		code = http.StatusUnauthorized
		err = fmt.Errorf("Unauthorized")
	} else {
		err = decodeJSONBody(w, r, &domain)
		if err == nil {
			resp, err = db.InsertDomain(&domain)
		}
	}
	httpSendResponse(w, code, resp, err)
}

// ReadDomains is an httpHandler for route GET /domains
func ReadDomains(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Get All Domains ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp []*model.Domain
	var code int

	if !reqIsSuperuser(r) {
		code = http.StatusUnauthorized
		err = fmt.Errorf("Unauthorized")
	} else {
		resp = db.SelectDomains()
	}
	httpSendResponse(w, code, resp, err)
}

// ReadDomain is an httpHandler for route GET /domains/{id}
func ReadDomain(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Get Domain By Id ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp *model.Domain
	var domainId int
	var domainName string
	var code int

	domainName, domainId = reqNameOrId(r)
	if !reqIsSuperuser(r) {
		// Can only read its own domain
		myDomainName, myDomainId := reqDomain(r)
		if myDomainId == domainId {
			domainName = myDomainName
		} else if myDomainName == domainName {
			domainId = myDomainId
		} else {
			// Not Authorized
			code = http.StatusUnauthorized
			err = fmt.Errorf("Unauthorized")
		}
	}
	if err == nil {
		if domainId == 0 && domainName == "" {
			// Unknown Domain
			err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
		} else {
			resp = db.SelectDomain(domainId, domainName)
			if resp != nil {
				fmt.Printf("SelectDomain %s %d returned %v\n", domainName, domainId, *resp)
			} else {
				fmt.Printf("SelectDomain %s %d returned nil\n", domainName, domainId)
			}
		}
	}
	if err != nil {
		fmt.Printf("SelectDomain %s %d failed: %v\n", domainName, domainId, err)
	}
	httpSendResponse(w, code, resp, err)
}

// UpdateDomain is an httpHandler for route PUT /domains
func UpdateDomain(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Update Domain By Id ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp *model.Domain
	var domainId int
	var domainName string
	var code int

	domainName, domainId = reqNameOrId(r)
	currDomainName, currDomainId := reqDomain(r)

	// Can only read its own domain
	if currDomainId == 0 && currDomainName == "" {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", currDomainName, currDomainId)
	} else if reqIsSuperuser(r) {
		// OK to update
	} else if currDomainId == domainId {
		domainName = currDomainName
	} else if currDomainName == domainName {
		domainId = currDomainId
	} else {
		// Not Authorized
		fmt.Printf("UpdateDomain: Cant update [%s %d]. Admin is of domain [%s %d]\n",
			domainName, domainId, currDomainName, currDomainId)
		code = http.StatusUnauthorized
		err = fmt.Errorf("Unauthorized")
	}
	if err == nil {
		// Decode the request body
		var domain model.DomainReq
		err = decodeJSONBody(w, r, &domain)
		if err == nil {
			resp = db.UpdateDomain(domainId, domainName, &domain)
		}
	}

	httpSendResponse(w, code, resp, err)
}

// ChangeDomain is an httpHandler for route GET /changedomain/{id}
func ChangeDomain(w http.ResponseWriter, r *http.Request) {
	var err error
	var newDomainId int
	var newDomainName string
	var code int

	if !reqIsSuperuser(r) {
		code = http.StatusUnauthorized
		err = fmt.Errorf("Unauthorized")
	} else {
		currDomainName, currDomainId := reqDomain(r)
		newDomainName, newDomainId = reqNameOrId(r)
		if newDomainId == 0 && newDomainName == "" {
			// Unknown Domain - Bad Request
			err = fmt.Errorf("neither Domain name nor ID specified")
		} else if currDomainId == 0 || currDomainName == "" {
			// Unknown Domain
			err = fmt.Errorf("domain %s %d unknown", currDomainName, currDomainId)
		} else if (newDomainName == currDomainName) || (newDomainId == currDomainId) {
			// Nothing to do
			fmt.Printf("ChangeDomain from [%s %d] to [%s %d] - No change needed\n",
				currDomainName, currDomainId,
				newDomainName, newDomainId)
		} else {
			// Change to this domain
			sessionId := r.Header.Get("Xpress-SessionId")
			err = db.ChangeDomain(sessionId, newDomainId, newDomainName)
			fmt.Printf("ChangeDomain from [%s %d] to [%s %d] %s err=%v\n",
				currDomainName, currDomainId,
				newDomainName, newDomainId, sessionId, err)
			if err == nil {
				//
				// Send response same as Login response
				//
				resp := db.SelectDomain(newDomainId, newDomainName)
				r.Header.Set("Xpress-IconId", "23") // TODO - get actual value from database
				r.Header.Set("Xpress-Domain", resp.Name)
				r.Header.Set("Xpress-DomainId", fmt.Sprintf("%d", resp.ID))
				Login(w, r)
				return
			}
		}
	}

	httpSendResponse(w, code, nil, err)
}

// DeleteDomain is an httpHandler for route DELETE /domain
func DeleteDomain(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Delete Domain ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp *model.Domain
	var code int

	if !reqIsSuperuser(r) {
		code = http.StatusUnauthorized
		err = fmt.Errorf("Unauthorized")
	} else {
		domainName, domainId := reqNameOrId(r)
		resp = db.DeleteDomain(domainId, domainName)
	}
	httpSendResponse(w, code, resp, err)
}
