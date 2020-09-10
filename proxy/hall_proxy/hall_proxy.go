package hall_proxy

import "log"

type HallProxy struct {
}

func Create() *HallProxy {
	return &HallProxy{}
}

func (hp *HallProxy) Start() error {
	log.Println("HallProxy Start")
	return nil
}

func (hp *HallProxy) Stop() {
	log.Println("HallProxy Stop")
}
