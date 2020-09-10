package master_proxy

import "log"

type MasterProxy struct {
}

func Create() *MasterProxy {
	return &MasterProxy{}
}

func (mp *MasterProxy) Start() error {
	log.Println("MasterProxy Start")
	return nil
}

func (mp *MasterProxy) Stop() {
	log.Println("MasterProxy Stop")
}
