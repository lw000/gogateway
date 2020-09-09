package users

import (
	"gogateway/agent"
	"gogateway/errors"
	"sync"
)

type ClientAgents struct {
	agents sync.Map
}

var (
	clientAgents     *ClientAgents
	clientAgentsOnce sync.Once
)

func New() *ClientAgents {
	return &ClientAgents{}
}

func Instance() *ClientAgents {
	clientAgentsOnce.Do(func() {
		clientAgents = New()
	})
	return clientAgents
}

func (c *ClientAgents) Add(clientId uint32, agent *agent.Agent) {
	c.agents.Store(clientId, agent)
}

func (c *ClientAgents) Remove(clientId uint32) {
	c.agents.Delete(clientId)
}

func (c *ClientAgents) WriteMessage(clientId uint32, data []byte) error {
	v, exists := c.agents.Load(clientId)
	if !exists {
		return errors.New(0, "客户端不存在")
	}
	clientAgent, ok := v.(*agent.Agent)
	if !ok {
		return errors.New(0, "内部错误")
	}
	return clientAgent.WriteMessage(data)
}
