package proxy

type InputParams struct {
	Type     string `json:"type"`
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}
