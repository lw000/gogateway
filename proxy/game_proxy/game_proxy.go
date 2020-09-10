package game_proxy

import "log"

type GameProxy struct {
}

func Create() *GameProxy {
	return &GameProxy{}
}

func (gp *GameProxy) Start() error {
	log.Println("GameProxy Start")
	return nil
}

func (gp *GameProxy) Stop() {
	log.Println("GameProxy Stop")
}
