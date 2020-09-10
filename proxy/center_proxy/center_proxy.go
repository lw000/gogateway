package center_proxy

import "log"

type CenterProxy struct {
}

func Create() *CenterProxy {
	return &CenterProxy{}
}

func (gp *CenterProxy) Start() error {
	log.Println("CenterProxy Start")
	return nil
}

func (gp *CenterProxy) Stop() {
	log.Println("CenterProxy Stop")
}
