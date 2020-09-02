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

func (c *ClientSessions) Add(sessionId uint32, s *melody.Session) {
	c.sessions.Store(sessionId, s)
}

func (c *ClientSessions) Remove(sessionId uint32) {
	c.sessions.Delete(sessionId)
}

func (c *ClientSessions) WriteMessage(clientId uint32, data []byte) error {
	v, exists := c.sessions.Load(clientId)
	if !exists {
		return errors.New(0, "客户端不存在")
	}
	session, ok := v.(*melody.Session)
	if !ok {
		return errors.New(0, "内部错误")
	}

	if !session.IsClosed() {
		return session.WriteBinary(data)
	}

	return nil
}
