package frontend

import (
	"github.com/lw000/gocommon/network/ws/packet"
	"github.com/lw000/gocommon/utils"
	"github.com/olahol/melody"
	log "github.com/sirupsen/logrus"
	"gogateway/agent"
	"gogateway/backend"
	"gogateway/users"
	"net/http"
	"sync"
	"time"
)

type MsgHooksFunc = func(pk *typacket.Packet) bool

type Server struct {
	m             *melody.Melody
	messagesHooks []MsgHooksFunc
}

var (
	serve     *Server
	serveOnce sync.Once
)

func New() *Server {
	serve := &Server{
		m: melody.New(),
	}
	return serve.init()
}

func Instance() *Server {
	serveOnce.Do(func() {
		serve = New()
	})
	return serve
}

func (serve *Server) init() *Server {
	serve.m.Config = &melody.Config{
		WriteWait:         10 * time.Second,
		PongWait:          30 * time.Second,
		PingPeriod:        (30 * time.Second * 9) / 10,
		MaxMessageSize:    1024 * 32,
		MessageBufferSize: 1024 * 32,
	}

	serve.m.HandleConnect(serve.onConnectHandler)
	serve.m.HandleMessageBinary(serve.onBinaryMessageHandler)
	serve.m.HandleDisconnect(serve.onDisconnectHandler)
	serve.m.HandleError(serve.onErrorHandler)
	return serve
}

func (serve *Server) Start() error {

	return nil
}

func (serve *Server) Stop() {
	if serve == nil {
		return
	}
	if err := serve.m.Close(); err != nil {
		log.Error(err)
	}
}

func (serve *Server) AddMessageHook(fun ...MsgHooksFunc) {
	serve.messagesHooks = append(serve.messagesHooks, fun...)
}

func (serve *Server) HandleRequest(w http.ResponseWriter, r *http.Request) error {
	return serve.m.HandleRequestWithKeys(w, r, nil)
}

func (serve *Server) onConnectHandler(s *melody.Session) {
	clientId := tyutils.HashCode(tyutils.UUID())
	s.Set("clientId", clientId)
	users.Instance().Add(clientId, agent.New(clientId, s))
	log.Infof("客户端连接, clientId:%d", clientId)
}

func (serve *Server) onBinaryMessageHandler(s *melody.Session, msg []byte) {
	clientId := serve.getClientId(s)
	if clientId <= 0 {
		_ = s.CloseWithMsg([]byte("error"))
		return
	}

	var (
		err error
		pk  *typacket.Packet
	)
	pk, err = typacket.NewPacketWithData(msg)
	if err != nil {
		_ = s.CloseWithMsg([]byte("core error"))
		return
	}

	var allowForward = true
	for _, hook := range serve.messagesHooks {
		allowForward = hook(pk)
		if !allowForward {
			break
		}
	}

	if !allowForward {
		_ = s.CloseWithMsg([]byte("core error"))
		return
	}

	// if pk.CheckCode() != 123456 {
	// 	err = s.CloseWithMsg([]byte("Illegal Packet"))
	// 	if err != nil {
	// 		log.Error(err)
	// 	}
	// 	return
	// }
	// log.Info("clientId:%d, pk:%+v", clientId, pk)

	switch true {
	case true:
		err = backend.Instance().WriteBinaryMessage(pk.Mid(), pk.Sid(), clientId, pk.Data())
		if err != nil {
			log.Error(err)
			err = s.CloseWithMsg([]byte("error"))
			if err != nil {
				log.Error(err)
			}
		}
	default:
		err = backend.Instance().WriteBinaryMessage(pk.Mid(), pk.Sid(), clientId, pk.Data())
		if err != nil {
			log.Error(err)
			err = s.CloseWithMsg([]byte("error"))
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func (serve *Server) onDisconnectHandler(s *melody.Session) {
	clientId := serve.getClientId(s)
	users.Instance().Remove(clientId)
	log.Infof("客户端断开, clientId:%d", clientId)
}

func (serve *Server) onErrorHandler(s *melody.Session, err error) {
	sessionId := serve.getClientId(s)
	log.Infof("客户端错误, clientId: %d, err:%serve", sessionId, err.Error())
}

func (serve *Server) getClientId(s *melody.Session) uint32 {
	v, exists := s.Get("clientId")
	if exists {
		return v.(uint32)
	}
	return 0
}
