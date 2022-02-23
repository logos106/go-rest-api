package model

type PolicyApp struct {
	Name      string       `json:"name,omitempty"`
	AllowedIPs string       `json:"allowed_ips,omitempty"`
	IsUserPolicy bool      `json:"is_user_policy,omitempty"`
}

type ServiceNode struct {
	Name      string       `json:"name,omitempty"`
	WGKey     string       `json:"wg_key,omitempty"`
	VirtualIP string       `json:"virtual_ip,omitempty"`
	PublicIP  string       `json:"public_ip,omitempty"`
	LocalIP   string       `json:"local_ip,omitempty"`
	Apps      []*PolicyApp       `json:"apps,omitempty"`
}

type Policy struct {
	ServiceNodes   map[string]*ServiceNode  `json:"services,omitempty"`
}
