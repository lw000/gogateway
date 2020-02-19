package config

import (
	"encoding/json"
	"io/ioutil"
)

type Server struct {
	Listen      int64
	Ssl         string
	SslCertfile string
	SslKeyfile  string
	Servername  []string
	Blacklist   []string
	Whitelist   []string
}

type Servers struct {
	Server []Server
}

type BackendConfig struct {
	Host   string `json:"host"`
	Scheme string `json:"scheme"`
	Path   string `json:"path"`
}

type JsonConfig struct {
	Debug       int64
	Servers     Servers
	BackendConf BackendConfig
}

func New() *JsonConfig {
	return &JsonConfig{}
}

func LoadJsonConfig(file string) (*JsonConfig, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var ccf CfgStruct
	if err = json.Unmarshal(data, &ccf); err != nil {
		return nil, err
	}

	cfg := New()
	cfg.Debug = ccf.Debug

	for _, v := range ccf.Servers.Server {
		cfg.Servers.Server = append(cfg.Servers.Server, Server{
			v.Listen,
			v.Ssl,
			v.SslCertfile,
			v.SslKeyfile,
			v.Servername,
			v.Blacklist,
			v.Whitelist,
		})
	}

	// 后端服务配置
	cfg.BackendConf.Host = ccf.Backend.Host
	cfg.BackendConf.Scheme = ccf.Backend.Scheme
	cfg.BackendConf.Path = ccf.Backend.Path

	return cfg, nil
}
