package proxy

type InputParams struct {
	Type     string   `json:"type"`
	Ip       string   `json:"ip"`
	Port     string   `json:"port"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	ReqUrls  []string `json:"req_urls"`
}

type RequestInfo struct {
	Proxy struct {
		Type     string `json:"type"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"proxy"`
	Request struct {
		Method string              `json:"method"`
		Urls   []string            `json:"urls"`
		Header map[string][]string `json:"header"`
		Body   string              `json:"body"`
	} `json:"request"`
	Timeout    int  `json:"timeout"`
	HiddenBody bool `json:"hidden_body"`
}
