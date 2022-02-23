package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/golang/gddo/httputil/header"
	"github.com/gorilla/mux"
	"github.com/saroopmathur/rest-api/db"
	model "github.com/saroopmathur/rest-api/models"
)

const (
	PAGE_OFFSET = 0
	PAGE_LIMIT  = 1000000
)

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}

func reqUser(r *http.Request) *model.User2 {
	u := model.User2{}
	u.Name = r.Header.Get("Xpress-User")
	u.ID, _ = strconv.Atoi(r.Header.Get("Xpress-UserId"))
	u.Domain.Name = r.Header.Get("Xpress-Domain")
	u.Domain.ID, _ = strconv.Atoi(r.Header.Get("Xpress-DomainId"))
	u.SessionID = r.Header.Get("Xpress-SessionId")
	u.Role = r.Header.Get("Xpress-Role")
	return &u
}

var NameRegexp *regexp.Regexp

const NAME_REGEX = "^[A-Za-z]([A-Za-z0-9@/:%+$#.|;<>?& _-]*[A-Za-z0-9])?$"

func init() {
	var err error
	NameRegexp, err = regexp.Compile(NAME_REGEX)
	if err != nil {
		log.Fatalf("Invlaid REGEXP %v\n", err)
	}
}

func IsValidName(s string) bool {
	return NameRegexp.Match([]byte(s))
}

// func reqRole(r *http.Request) string {
// 	role := r.Header.Get("Xpress-Role")
// 	return role
// }

func reqIsSuperuser(r *http.Request) bool {
	role := r.Header.Get("Xpress-Role")
	return role == db.ROLE_POWERADMIN
}

func reqDomain(r *http.Request) (string, int) {
	var domainId int
	var err error

	domainName := r.Header.Get("Xpress-Domain")
	idStr := r.Header.Get("Xpress-DomainId")
	if idStr != "" {
		domainId, err = strconv.Atoi(idStr)
		if err != nil {
			domainId = 0
		}
	}
	//fmt.Printf("REQ: Domain:%s DomainId:%d\n", domainName, domainId)
	return domainName, domainId
}

func reqNameOrId(r *http.Request) (string, int) {
	var name string
	params := mux.Vars(r)
	str := params["id"]
	id, err := strconv.Atoi(str)
	if err != nil {
		id = 0
		name = str
		if !IsValidName(name) {
			name = "INVALID"
		}
	} else if id == 0 {
		name = "INVALID"
	}
	//fmt.Printf("REQ: Name:%s Id:%d\n", name, id)
	return name, id
}

func reqNameOrId2(r *http.Request) (string, int) {
	var name string
	params := mux.Vars(r)
	str := params["id2"]
	id, err := strconv.Atoi(str)
	if err != nil {
		id = 0
		name = str
		if !IsValidName(name) {
			name = "INVALID"
		}
	} else if id == 0 {
		name = "INVALID"
	}
	//fmt.Printf("REQ: Name:%s Id:%d\n", name, id)
	return name, id
}

func reqPageInfo(r *http.Request) (int, int) {
	offset := PAGE_OFFSET
	keys, ok := r.URL.Query()["offset"]
	if ok {
		offset, _ = strconv.Atoi(keys[0])
	}

	limit := PAGE_LIMIT
	keys, ok = r.URL.Query()["limit"]
	if ok {
		limit, _ = strconv.Atoi(keys[0])
	}

	return offset, limit
}

func reqSearchString(r *http.Request) string {
	search := ""

	keys, ok := r.URL.Query()["search"]
	if ok {
		search = keys[0]
	}

	return search
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func httpSendResponse(w http.ResponseWriter, code int, resp interface{}, err error) {
	if code == 0 {
		// code not specified, determine based on error
		if err == nil {
			code = http.StatusOK
		} else {
			if strings.Contains(err.Error(), "unique constraint") {
				code = http.StatusConflict
			} else if strings.Contains(err.Error(), "Unauthorized") {
				code = http.StatusUnauthorized
			} else {
				code = http.StatusBadRequest
			}
		}
	}

	if resp == nil && err != nil {
		// build resp based on error message
		out := &Response{}
		out.Data.Status = code
		out.Code = ""
		out.Message = err.Error()

		resp = out
	}

	if isNil(resp) {
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(code)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}
