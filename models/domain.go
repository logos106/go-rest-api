package model

type DomainReq struct {
	Name string `json:"name,omitempty"`
}

type Domain struct {
	ID  int `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Status string `json:"-"`
}
