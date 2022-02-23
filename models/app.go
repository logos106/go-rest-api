package model

// Application data as sent by UI
type AppReq struct {
	Name string `json:"name,omitempty"`
	ServiceName string `json:"service_name,omitempty"`
	ServiceId int `json:"service_id,omitempty"`
	AllowedIPs string `json:"allowed_ips,omitempty"`
}

type App struct {
	ID int `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	ServiceName string `json:"service_name,omitempty"`
	ServiceId int `json:"service_id,omitempty"`
	AllowedIPs string `json:"allowed_ips,omitempty"`
}
