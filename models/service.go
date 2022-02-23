package model

type Service struct {
	Name      string `json:"name,omitempty"`
	Password  string `json:"password,omitempty"`
	WGKey     string `json:"wg_key,omitempty"`
	PublicIP  string `json:"public_ip,omitempty"`
	VirtualIP string `json:"virtual_ip,omitempty"`
	LocalIP   string `json:"local_ip,omitempty"`
}

type Service2 struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Domain    Domain `json:"-"`
	Password  string `json:"-"`
	Icon      int    `json:"-"`
	Status    string `json:"-"`
	WGKey     string `json:"wg_key,omitempty"`
	PublicIP  string `json:"public_ip,omitempty"`
	VirtualIP string `json:"virtual_ip,omitempty"`
	LocalIP   string `json:"local_ip,omitempty"`
	SessionID string `json:"-"`
}
