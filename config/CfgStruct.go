package config

type CfgStruct struct {
	Backend struct {
		Host   string `json:"host"`
		Path   string `json:"path"`
		Scheme string `json:"scheme"`
	} `json:"backend"`
	Debug   int64 `json:"debug"`
	Servers struct {
		Server []struct {
			Blacklist   []string `json:"blacklist"`
			Listen      int64    `json:"listen"`
			Servername  []string `json:"servername"`
			Ssl         string   `json:"ssl"`
			SslCertfile string   `json:"ssl_certfile"`
			SslKeyfile  string   `json:"ssl_keyfile"`
			Whitelist   []string `json:"whitelist"`
		} `json:"server"`
	} `json:"servers"`
}
