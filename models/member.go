package model

type Member struct {
	ID  int `json:"id,omitempty"`
	Group int `json:"group,omitempty"`
	User int `json:"user,omitempty"`
}
