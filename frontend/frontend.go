package frontend

import (
	"github.com/lw000/gocommon/network/ws/packet"
	"github.com/lw000/gocommon/utils"
	"github.com/olahol/melody"
	log "github.com/sirupsen/logrus"
	"gogateway/backend"
	"gogateway/users"
	"net/http"
	"sync"
	"time"
)

type MsgHooks = func(pk *typacket.Packet) bool

type Server struct {
	m             *melody.Melody
	messagesHooks []MsgHooks
}

var (
	serveInstance     *Server
	serveInstanceOnce sync.Once
)

func New() *Server {
	svr := &Server{
		m: melody.New(),
	}
	return svr.init()
}

func Instance() *Server {
	serveInstanceOnce.Do(func() {
		serveInstance = New()
	})
	return serveInstance
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

func (serve *Server) AddMessageHook(fun MsgHooks) {
	serve.messagesHooks = append(serve.messagesHooks, fun)
}

func (serve *Server) HandleRequest(w http.ResponseWriter, r *http.Request) error {
	return serve.m.HandleRequestWithKeys(w, r, nil)
}

func (serve *Server) onConnectHandler(s *melody.Session) {
	sessionId := tyutils.HashCode(tyutils.UUID())
	s.Set("sessionId", sessionId)
	users.Instance().Add(sessionId, s)
	log.Infof("客户端连接, sessionId:%d", sessionId)
}

func (serve *Server) onErrorHandler(s *melody.Session, e error) {
	sessionId := serve.getSessionId(s)
	log.Infof("客户端错误, sessionId: %d, err:%s", sessionId, e.Error())
}

func (serve *Server) onBinaryMessageHandler(s *melody.Session, msg []byte) {
	var err error
	value, exists := s.Get("sessionId")
	if !exists {
		_ = s.CloseWithMsg([]byte("error"))
		return
	}
	sessionId := value.(uint32)

	var pk *typacket.Packet
	pk, err = typacket.NewPacketWithData(msg)
	if err != nil {
		_ = s.CloseWithMsg([]byte("error"))
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
		_ = s.CloseWithMsg([]byte("error"))
		return
	}

	// if pk.CheckCode() != 123456 {
	// 	err = s.CloseWithMsg([]byte("Illegal Packet"))
	// 	if err != nil {
	// 		log.Error(err)
	// 	}
	// 	return
	// }

	// log.Info("sessionId:%d, pk:%+v", sessionId, pk)

	switch true {
	case true:
		err = backend.Instance().WriteBinaryMessage(pk.Mid(), pk.Sid(), sessionId, pk.Data())
		if err != nil {
			log.Error(err)
			err = s.CloseWithMsg([]byte("error"))
			if err != nil {
				log.Error(err)
			}
		}
	default:
		err = backend.Instance().WriteBinaryMessage(pk.Mid(), pk.Sid(), sessionId, pk.Data())
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
	sessionId := serve.getSessionId(s)
	users.Instance().Remove(sessionId)
	log.Infof("客户端断开, sessionId:%d", sessionId)
}

func (serve *Server) getSessionId(s *melody.Session) uint32 {
	v, exists := s.Get("sessionId")
	if exists {
		return v.(uint32)
	}
	return 0
}
