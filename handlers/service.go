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

func CreateService(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Add Service ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)
	data, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewReader(data))
	log.Printf("%s\n", data)

	var service model.Service
	var resp *model.Service2
	err := decodeJSONBody(w, r, &service)
	if err == nil {
		resp, err = CreateService1(r, &service)
	}
	httpSendResponse(w, 0, resp, err)
}

func CreateService1(r *http.Request, service *model.Service) (*model.Service2, error) {
	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err := fmt.Errorf("domain %s %d unknown", domainName, domainId)
		return nil, err
	}

	return db.InsertService(domainId, service)
}

// ReadServices is an httpHandler for route GET /services
func ReadServices(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Get All Services ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp []*model.Service2

	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		resp = db.SelectServices(domainId)
	}
	httpSendResponse(w, 0, resp, err)
}

// ReadService is an httpHandler for route GET /services/{id}
func ReadService(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Get Service ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp *model.Service2

	serviceName, serviceId := reqNameOrId(r)
	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		resp = db.SelectService(domainId, serviceName, serviceId)
	}
	httpSendResponse(w, 0, resp, err)
}

// UpdateService is an httpHandler for route PUT /services
func UpdateService(w http.ResponseWriter, r *http.Request) {
	// Logging
	log.Printf("============== Update Services ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)
	data, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewReader(data))
	log.Printf("%s\n", data)

	var err error
	var resp *model.Service2

	serviceName, serviceId := reqNameOrId(r)
	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		// Decode the request body
		var service model.Service
		err = decodeJSONBody(w, r, &service)
		if err == nil {
			resp = db.UpdateService(domainId, serviceName, serviceId, &service)
		}
	}

	httpSendResponse(w, 0, resp, err)
}

// DeleteService is an httpHandler for route DELETE /service
func DeleteService(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Delete Services ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var err error
	var resp *model.Service2

	serviceName, serviceId := reqNameOrId(r)
	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		resp := db.DeleteService(domainId, serviceName, serviceId)
		if resp == nil {
			err = fmt.Errorf("unknown Service")
		}
		// resp is the service obejct for the deleted service
	}
	httpSendResponse(w, 0, resp, err)
}
