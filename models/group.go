package model

type Group struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Group2 struct {
	ID     int    `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Domain Domain `json:"-"`
	Count  int    `json:"count,omitempty"`
	Role   string `json:"-"`
	Status string `json:"-"`
}

type GroupAccess struct {
	ID     int    `json:"id,omitempty"`
	Group  int    `json:"group,omitempty"`
	App    int    `json:"allowed,omitempty"`
	Status string `json:"status,omitempty"`
}

type GroupAccess2 struct {
	ID    int    `json:"id,omitempty"`
	Group string `json:"group,omitempty"`
	Apps  []App  `json:"apps,omitempty"`
}
