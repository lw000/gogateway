package users

import (
	"github.com/olahol/melody"
	"gogateway/errors"
	"sync"
)

type ClientSessions struct {
	sessions sync.Map
}

var (
	clientSessions     *ClientSessions
	clientSessionsOnce sync.Once
)

func New() *ClientSessions {
	return &ClientSessions{}
}

func Instance() *ClientSessions {
	clientSessionsOnce.Do(func() {
		clientSessions = New()
	})
	return clientSessions
}

func (cs *ClientSessions) Add(sessionId uint32, s *melody.Session) {
	cs.sessions.Store(sessionId, s)
}

func (cs *ClientSessions) Remove(sessionId uint32) {
	cs.sessions.Delete(sessionId)
}

func (cs *ClientSessions) WriteMessage(clientId uint32, data []byte) error {
	value, exists := cs.sessions.Load(clientId)
	if !exists {
		return errors.New(0, "客户端不存在")
	}

	session, ok := value.(*melody.Session)
	if ok {
		if !session.IsClosed() {
			return session.WriteBinary(data)
		}
	}
	return errors.New(0, "未知错误")
}
