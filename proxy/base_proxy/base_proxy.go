package base_proxy

type BaseProxy interface {
	Start() error
	Stop()
}
