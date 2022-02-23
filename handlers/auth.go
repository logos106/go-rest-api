package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/saroopmathur/rest-api/db"
	m "github.com/saroopmathur/rest-api/models"
)

type LoginResp struct {
	Token    string `json:"Token"`
	Role     string `json:"role,omitempty"`
	Domain   string `json:"domain,omitempty"`
	LogoFile string `json:"logo_file,omitempty"`
	IconFile string `json:"icon_file,omitempty"`
}

func GetLogoFile(domainName string, domainId int) string {
	return fmt.Sprintf("images/logos/%d.png", domainId)
}

func GetIconFile(domainName string, domainId int, iconId string) string {
	return fmt.Sprintf("images/icons/%s.png", iconId)
}

func Login(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Xpress-SessionId")
	role := r.Header.Get("Xpress-Role")
	iconId := r.Header.Get("Xpress-IconId")
	domainName, domainId := reqDomain(r)

	// Send the token as reponse
	resp := &LoginResp{}
	resp.Token = token
	resp.Role = role
	resp.Domain = domainName
	resp.LogoFile = GetLogoFile(domainName, domainId)
	resp.IconFile = GetIconFile(domainName, domainId, iconId)
	httpSendResponse(w, 0, resp, nil)
}

func UserCheckPassword(u *m.User2, pass string) bool {
	fmt.Printf("Stored=%s, Input=%s\n", u.Password, pass)
	return strings.Compare(u.Password, pass) == 0
}

func ServiceCheckPassword(s *m.Service2, pass string) bool {
	return strings.Compare(s.Password, pass) == 0
}

func AdminCheckPassword(a *m.Admin2, pass string) bool {
	return strings.Compare(a.Password, pass) == 0
}

func Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Xpress-SessionId")
	if token != "" {
		db.TokenInvalidate(token)
	}
	httpSendResponse(w, 0, nil, nil)
}
