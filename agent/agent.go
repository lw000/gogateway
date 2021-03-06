package agent

import (
	"github.com/olahol/melody"
	"time"
)

type Agent struct {
	userId    uint32
	clientId  uint32
	session   *melody.Session
	connTime  time.Time
	loginTime time.Time
	leaveTime time.Time
}

func New(clientId uint32, session *melody.Session) *Agent {
	return &Agent{
		clientId: clientId,
		session:  session,
		connTime: time.Now(),
	}
}

func (a *Agent) ClientIdId() uint32 {
	return a.clientId
}

func (a *Agent) SendClientMessage(data []byte) error {
	if !a.session.IsClosed() {
		return a.session.WriteBinary(data)
	}
	return nil
}

func (a *Agent) Login() error {
	a.loginTime = time.Now()
	return nil
}

func (a *Agent) Leave() error {
	a.leaveTime = time.Now()
	return nil
}
