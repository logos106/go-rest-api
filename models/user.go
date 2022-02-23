package model

// User data as sent by UI
type User struct {
	Name      string `json:"name,omitempty"`
	Password  string `json:"password,omitempty"`
	WGKey     string `json:"wg_key,omitempty"`
	PublicIP  string `json:"public_ip,omitempty"`
	VirtualIP string `json:"virtual_ip,omitempty"`
	LocalIP   string `json:"local_ip,omitempty"`
}

type User2 struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Domain    Domain `json:"-"`
	Group     Group  `json:"-"`
	Password  string `json:"-"`
	WGKey     string `json:"wg_key,omitempty"`
	PublicIP  string `json:"public_ip,omitempty"`
	VirtualIP string `json:"virtual_ip,omitempty"`
	LocalIP   string `json:"local_ip,omitempty"`
	Role      string `json:"-"`
	Status    string `json:"-"`
	SessionID string `json:"-"`
}

type UserAccess struct {
	ID     int    `json:"id,omitempty"`
	User   int    `json:"user,omitempty"`
	App    int    `json:"allowed,omitempty"`
	Status string `json:"status,omitempty"`
}

type UserAccess2 struct {
	ID   int    `json:"id,omitempty"`
	User string `json:"user,omitempty"`
	Apps []App  `json:"apps,omitempty"`
}
