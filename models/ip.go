package model

type IP struct {
	ID int `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	ServiceId int `json:"service,omitempty"`
	Allowed string `json:"allowed,omitempty"`
	Status string `json:"status,omitempty"`
}

type IP2 struct {
	ID int `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Service Service2 `json:"service,omitempty"`
	Allowed string `json:"allowed,omitempty"`
	Status string `json:"status,omitempty"`
}
