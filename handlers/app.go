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

func CreateApp(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Add App ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)
	data, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewReader(data))
	log.Printf("%s\n", data)

	var app model.AppReq
	var resp *model.App
	var err error

	domainName, domainId := reqDomain(r)
	if domainId == 0 { // Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		err = decodeJSONBody(w, r, &app)
		if err == nil {
			resp, err = db.InsertApp(domainId, &app)
		}
	}

	httpSendResponse(w, 0, resp, err)
}

// ReadApps is an httpHandler for route GET /apps
func ReadApps(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Get All Apps ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp []*model.App

	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		resp = db.SelectApps(domainId)
	}
	httpSendResponse(w, 0, resp, err)
}

// ReadApp is an httpHandler for route GET /apps/{id}
func ReadApp(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Get App By Id ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp *model.App

	domainName, domainId := reqDomain(r)
	appName, appId := reqNameOrId(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		resp = db.SelectApp(domainId, appName, appId)
		fmt.Printf("ReadApp %s %d Domain %s %v\n", appName, appId, domainName, resp)
	}

	httpSendResponse(w, 0, resp, err)
}

// ReadApp is an httpHandler for route GET /apps/{id}
func ReadApp2(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Get App By Service Id ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp *model.App

	domainName, domainId := reqDomain(r)
	serviceName, serviceId := reqNameOrId(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		resp = db.SelectApp2(domainId, serviceName, serviceId)
		fmt.Printf("ReadApp %s %d Domain %s %v\n", serviceName, serviceId, domainName, resp)
	}

	httpSendResponse(w, 0, resp, err)
}

// UpdateApp is an httpHandler for route PUT /apps
func UpdateApp(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Update App ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)
	data, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewReader(data))
	log.Printf("%s\n", data)

	var err error
	var resp *model.App

	appName, appId := reqNameOrId(r)
	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		// Decode the request body
		var app model.AppReq
		err = decodeJSONBody(w, r, &app)
		if err == nil {
			resp = db.UpdateApp(domainId, appName, appId, &app)
		}
		fmt.Printf("Update App %s %d Domain %s %v\n", appName, appId, domainName, resp)
	}

	httpSendResponse(w, 0, resp, err)
}

// DeleteApp is an httpHandler for route DELETE /app
func DeleteApp(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Delete App By Id ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp *model.App

	appName, appId := reqNameOrId(r)
	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		resp := db.DeleteApp(domainId, appName, appId)
		if resp == nil {
			err = fmt.Errorf("unknown App")
		}
		// resp is the app obejct for the deleted app
		fmt.Printf("Delete App %s %d Domain %s %v\n", appName, appId, domainName, resp)
	}

	httpSendResponse(w, 0, resp, err)
}
