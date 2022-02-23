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

func ValidateAdminReq(admin *model.Admin) error {
	if IsValidName(admin.Name) {
		return nil
	}
	return fmt.Errorf("invalid Name")
}

// CreateAdmin is an httpHandler for route POST /admins
func CreateAdmin(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Add Admin ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)
	data, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewReader(data))
	log.Printf("%s\n", data)

	var admin model.Admin
	var resp *model.Admin2
	err := decodeJSONBody(w, r, &admin)
	if err == nil {
		err = ValidateAdminReq(&admin)
		if err == nil {
			resp, err = CreateAdmin1(r, &admin)
		}
	}
	if resp != nil {
		fmt.Printf("CreateAdmin: %s %v\n", admin.Name, *resp)
	} else {
		fmt.Printf("CreateAdmin: %s NULL resp\n", admin.Name)
	}
	httpSendResponse(w, 0, resp, err)
}

func CreateAdmin1(r *http.Request, admin *model.Admin) (*model.Admin2, error) {
	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err := fmt.Errorf("domain %s %d unknown", domainName, domainId)
		return nil, err
	}

	return db.InsertAdmin(domainId, admin.Name, admin.Password)
}

// ReadAdmins is an httpHandler for route GET /admins
func ReadAdmins(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Read All Admins ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	domainName, domainId := reqDomain(r)
	var resp []*model.Admin2
	var err error
	if domainId == 0 {
		// Unknown Domain
		err = fmt.Errorf("domain %s %d unknown", domainName, domainId)
	} else {
		resp = db.SelectAdmins(domainId)
	}
	httpSendResponse(w, 0, resp, err)
}

// ReadAdmin is an httpHandler for route GET /admins/{id}
func ReadAdmin(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Get an Admin ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	resp, code, err := ReadAdmin1(r)
	if resp != nil {
		fmt.Printf("ReadAdmin: %v\n", *resp)
	} else {
		fmt.Printf("ReadAdmin: NULL resp %v\n", err)
	}
	httpSendResponse(w, code, resp, err)
}

func ReadAdmin1(r *http.Request) (*model.Admin2, int, error) {
	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err := fmt.Errorf("domain %s %d unknown", domainName, domainId)
		return nil, http.StatusNotFound, err
	}
	adminName, adminId := reqNameOrId(r)
	fmt.Printf("ReadAdmin: %s %d\n", adminName, adminId)
	admin := db.SelectAdmin(domainId, adminName, adminId)
	if admin == nil {
		fmt.Printf("ReadAdmin: %s %d SelectAdmin return NULL\n", adminName, adminId)
		return nil, http.StatusNotFound, nil
	}
	fmt.Printf("ReadAdmin: %s %d SelectAdmin return %v\n", adminName, adminId, *admin)
	return admin, http.StatusOK, nil
}

// CreateAdmin is an httpHandler for route PUT /admins
func UpdateAdmin(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Update Admin ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)
	data, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewReader(data))
	log.Printf("%s\n", data)

	var admin model.Admin
	var resp *model.Admin2
	err := decodeJSONBody(w, r, &admin)
	if err == nil {
		err = ValidateAdminReq(&admin)
		if err == nil {
			resp, err = UpdateAdmin1(r, &admin)
		}
	}
	httpSendResponse(w, 0, resp, err)
}

func UpdateAdmin1(r *http.Request, admin *model.Admin) (*model.Admin2, error) {
	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err := fmt.Errorf("domain %s %d unknown", domainName, domainId)
		return nil, err
	}
	adminName, adminId := reqNameOrId(r)
	resp := db.UpdateAdmin(domainId, adminName, adminId, admin)
	return resp, nil
}

// DeleteAdmin is an httpHandler for route DELETE /admin
func DeleteAdmin(w http.ResponseWriter, r *http.Request) {
	err := DeleteAdmin1(r)
	httpSendResponse(w, 0, nil, err)
}

func DeleteAdmin1(r *http.Request) error {
	log.Printf("============== Delete Admin ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	domainName, domainId := reqDomain(r)
	if domainId == 0 {
		// Unknown Domain
		err := fmt.Errorf("domain %s %d unknown", domainName, domainId)
		return err
	}
	adminName, adminId := reqNameOrId(r)
	return db.DeleteAdmin(domainId, adminName, adminId)
}
