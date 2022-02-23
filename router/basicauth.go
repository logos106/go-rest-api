package router

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/saroopmathur/rest-api/db"
	h "github.com/saroopmathur/rest-api/handlers"
	m "github.com/saroopmathur/rest-api/models"
)

func setReqHeaders(r *http.Request, role string, name string, domainName string, id int, domainId int, sessionId string) {
	r.Header.Add("Xpress-User", name)
	r.Header.Add("Xpress-UserId", fmt.Sprintf("%d", id))
	r.Header.Add("Xpress-Domain", domainName)
	r.Header.Add("Xpress-DomainId", fmt.Sprintf("%d", domainId))
	r.Header.Add("Xpress-SessionId", sessionId)
	r.Header.Add("Xpress-Role", role)
	r.Header.Add("Xpress-IconId", "23") // TODO - get actual value from database
}

func adminLoginMiddleware(r *http.Request) error {
	log.Printf("============== Login as Admin ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var a *m.Admin2
	var err error

	name, pass, ok := r.BasicAuth()
	if ok {
		a = db.GetAdminByName(name)
		if a == nil || !h.AdminCheckPassword(a, pass) {
			err = fmt.Errorf("invalid admin username or password")
		}
	} else {
		err = fmt.Errorf("please enter your username and password")
	}

	if err != nil {
		fmt.Printf("Admin Basic Authentication Failed: %s %s - %v\n", name, pass, err)
		return err
	}

	if a.Domain.ID == db.POWERDOMAIN {
		a.Role = db.ROLE_POWERADMIN
		db.SetAnyDomain(&a.Domain)
	}

	db.GenerateAndSaveAdminToken(a)
	log.Printf("Admin Login Successful: Role=%s %s@%s [%d %d] %s\n", a.Role, a.Name, a.Domain.Name, a.ID, a.Domain.ID, a.SessionID)
	setReqHeaders(r, a.Role, a.Name, a.Domain.Name, a.ID, a.Domain.ID, a.SessionID)
	return nil
}

func serviceLoginMiddleware(r *http.Request) error {
	log.Printf("============== Login as Service ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	// Check if this is login request
	var s *m.Service2
	var err error
	// Login
	name, pass, ok := r.BasicAuth()
	if ok {
		log.Printf("Got Service Login Request: %s %s\n", name, pass)
		s = db.GetServiceByName(name)
		if s == nil || !h.ServiceCheckPassword(s, pass) {
			err = fmt.Errorf("invalid name or password")
		}
	} else {
		err = fmt.Errorf("please enter your username and password")
	}
	if err != nil {
		return err
	}
	db.GenerateAndSaveServiceToken(s)
	log.Printf("Service Login Successful: %s@%s [%d %d] %s\n", s.Name, s.Domain.Name, s.ID, s.Domain.ID, s.SessionID)
	setReqHeaders(r, db.ROLE_SERVICE, s.Name, s.Domain.Name, s.ID, s.Domain.ID, s.SessionID)
	return nil
}

func userLoginMiddleware(r *http.Request) error {
	log.Printf("============== Login as User ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	// Check if this is login request
	var u *m.User2
	var err error

	user, pass, ok := r.BasicAuth()
	if ok {
		log.Printf("Got User Login Request: %s %s\n", user, pass)
		u = db.GetUserByName(user)
		if u == nil || !h.UserCheckPassword(u, pass) {
			err = fmt.Errorf("invalid username or password")
		}
	} else {
		err = fmt.Errorf("please enter your username and password")
	}
	if err != nil {
		return err
	}
	db.GenerateAndSaveUserToken(u)
	log.Printf("User Login Successful: %s@%s [%d %d] %s\n", u.Name, u.Domain.Name, u.ID, u.Domain.ID, u.SessionID)
	setReqHeaders(r, db.ROLE_USER, u.Name, u.Domain.Name, u.ID, u.Domain.ID, u.SessionID)
	return nil
}

func tokenLoginMiddleware(r *http.Request) error {
	// Check Token
	token := GetToken(r.Header.Get("Authorization"))
	if token == "" {
		err := fmt.Errorf("Unauthorized")
		return err
	}
	//log.Printf("Login Authorization Token: %s\n", token)

	// First check if this is a user token
	u := db.GetUserByToken(token)
	if u != nil {
		// User logged in successfully
		log.Printf("Token User Login Successful: %s@%s [%d %d] %s\n", u.Name, u.Domain.Name, u.ID, u.Domain.ID, u.SessionID)
		setReqHeaders(r, db.ROLE_USER, u.Name, u.Domain.Name, u.ID, u.Domain.ID, token)
		return nil
	}

	// Next check if this is a service token
	s := db.GetServiceByToken(token)
	if s != nil {
		// Service logged in successfully
		log.Printf("Token Service Login Successful: %s@%s [%d %d] %s\n", s.Name, s.Domain.Name, s.ID, s.Domain.ID, s.SessionID)
		setReqHeaders(r, db.ROLE_SERVICE, s.Name, s.Domain.Name, s.ID, s.Domain.ID, token)
		return nil
	}

	// Next check if this is an admin token
	a := db.GetAdminByToken(token)
	if a != nil {
		// Admin logged in successfully
		//log.Printf("Token Admin Login Successful: Role=%s %s@%s [%d %d] %s\n",
		//	a.Role, a.Name, a.Domain.Name, a.ID, a.Domain.ID, a.SessionID)
		setReqHeaders(r, a.Role, a.Name, a.Domain.Name, a.ID, a.Domain.ID, token)
		return nil
	}
	err := fmt.Errorf("session Expired. Login Again")
	log.Printf("Token Login Failed: %s %s\n", token, err.Error())
	return err
}

var APIBase string = "/api/v1"

func BasicAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		url := r.URL.String()
		if url == APIBase+"/login" {
			err = userLoginMiddleware(r)
		} else if url == APIBase+"/servicelogin" {
			err = serviceLoginMiddleware(r)
		} else if url == APIBase+"/adminlogin" {
			err = adminLoginMiddleware(r)
		} else {
			err = tokenLoginMiddleware(r)
			if err == nil {
				role := r.Header.Get("Xpress-Role")
				if strings.HasPrefix(url, APIBase+"/userapi/") {
					if role != db.ROLE_USER {
						// client API only for User role
						err = fmt.Errorf("unauthorized for this API - Bad Role '%s' Must be a User", role)
					}
				} else if strings.HasPrefix(url, APIBase+"/serviceapi/") {
					if role != db.ROLE_SERVICE {
						// Service API only for Service role
						err = fmt.Errorf("unauthorized for this API - Bad Role '%s' Must be a Service", role)
					}
				} else if role != db.ROLE_ADMIN && role != db.ROLE_POWERADMIN {
					// Must be Admin for all other APIs
					err = fmt.Errorf("unauthorized for this API - Bad Role '%s' Must be an Admin", role)
				}
			}
		}

		if err != nil {
			w.Header().Set("WWW-Authenticate", fmt.Sprintf("Basic realm=\"%s\"", err.Error()))
			w.WriteHeader(401)
			w.Write([]byte("Unauthorized.\n"))
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}

func GetToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, "Bearer")
	if len(parts) != 2 {
		return ""
	}

	token := strings.TrimSpace(parts[1])
	if len(token) < 1 {
		return ""
	}

	return token
}
