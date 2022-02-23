package model

type Admin struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

type Admin2 struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Domain    Domain `json:"-"`
	Password  string `json:"-"`
	Role      string `json:"-"`
	Status    string `json:"-"`
	SessionID string `json:"-"`
}
