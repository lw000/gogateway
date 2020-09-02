package backend

import (
	"github.com/golang/protobuf/proto"
	"github.com/lw000/gocommon/network/ws/packet"
	"github.com/lw000/gocommon/utils"
	log "github.com/sirupsen/logrus"
	"gogateway/backend/ws"
	"gogateway/constant"
	"gogateway/global"
	"gogateway/protos/serve"
	"gogateway/users"
	"sync"
)

type Server struct {
	wsc *client.WsClient
}

const (
	MaxMessageSize    = 1024
	MessageBufferSize = 1024
)

var (
	_serveInstance     *Server
	_serveInstanceOnce sync.Once
)

func New() *Server {
	cfg := client.Config{MaxMessageSize: MaxMessageSize, MessageBufferSize: MessageBufferSize}
	serve := &Server{
		wsc: client.New(cfg),
	}
	serve.init()
	return serve
}

func Instance() *Server {
	_serveInstanceOnce.Do(func() {
		_serveInstance = New()
	})
	return _serveInstance
}

func (s *Server) init() *Server {
	s.wsc.HandleMessageBinary(s.onMessageBinaryHandler)
	s.wsc.HandleError(s.onErrorHandler)
	return s
}

func (s *Server) Start() error {
	err := s.wsc.Open(global.ProjectConfig.BackendConf.Scheme, global.ProjectConfig.BackendConf.Host, global.ProjectConfig.BackendConf.Path)
	if err != nil {
		log.Error(err)
		return err
	}

	go s.wsc.Run()

	err = s.registerService()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *Server) Stop() {
	if s == nil {
		return
	}
	s.wsc.Close()
}

func (s *Server) WriteProtoMessage(mid, sid uint16, clientId uint32, pb proto.Message) error {
	return s.wsc.WriteProtoMessage(mid, sid, clientId, pb)
}

func (s *Server) WriteBinaryMessage(mid, sid uint16, clientId uint32, data []byte) error {
	return s.wsc.WriteBinaryMessage(mid, sid, clientId, data)
}

// 注册服务到路由服务器
func (s *Server) registerService() error {
	req := Tserve.ReqRegService{
		ServerId: constant.GatewayServerId,
		SvrType:  constant.GatewayServerType,
	}
	err := s.wsc.WriteProtoMessage(constant.MdmGatewayService, constant.SubGatewayServiceRegister, tyutils.HashCode(tyutils.UUID()), &req)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (s *Server) onMessageBinaryHandler(msg []byte) error {
	pk, err := typacket.NewPacketWithData(msg)
	if err != nil {
		log.Error("接收到非法协议")
		return err
	}

	switch pk.Mid() {
	case constant.MdmGatewayService:
		switch pk.Sid() {
		case constant.SubGatewayServiceRegister:
			var ack Tserve.AckRegService
			if err = proto.Unmarshal(pk.Data(), &ack); err != nil {
				log.Error(err)
				return err
			}

			if ack.Result != 0 {
				log.Error(ack.Errmsg)
				return err
			}
			log.Infof("服务注册成功 %+v", ack)
		}
	default:
		ack := typacket.NewPacket(pk.Mid(), pk.Sid(), pk.ClientId())
		if err = ack.Encode(pk.Data()); err != nil {
			log.Error(err)
			return err
		}
		err = users.Instance().WriteMessage(pk.ClientId(), ack.Data())
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func (s *Server) onErrorHandler(w *client.WsClient, err error) {
}
