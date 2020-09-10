package backend

import (
	"github.com/golang/protobuf/proto"
	"github.com/lw000/gocommon/network/ws/packet"
	"github.com/lw000/gocommon/utils"
	log "github.com/sirupsen/logrus"
	"gogateway/backend/ws"
	"gogateway/constants"
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

func (serve *Server) init() *Server {
	serve.wsc.HandleMessageBinary(serve.onMessageBinaryHandler)
	serve.wsc.HandleError(serve.onErrorHandler)
	return serve
}

func (serve *Server) Start() error {
	err := serve.wsc.Open(global.ProjectConfig.BackendConf.Scheme,
		global.ProjectConfig.BackendConf.Host,
		global.ProjectConfig.BackendConf.Path)
	if err != nil {
		log.Error(err)
		return err
	}

	go serve.wsc.Run()

	err = serve.registerService()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (serve *Server) Stop() {
	if serve == nil {
		return
	}
	serve.wsc.Close()
}

func (serve *Server) WriteProtoMessage(mid, sid uint16, clientId uint32, pb proto.Message) error {
	return serve.wsc.WriteProtoMessage(mid, sid, clientId, pb)
}

func (serve *Server) WriteBinaryMessage(mid, sid uint16, clientId uint32, data []byte) error {
	return serve.wsc.WriteBinaryMessage(mid, sid, clientId, data)
}

// 注册服务到中心服务器
func (serve *Server) registerService() error {
	req := Tserve.ReqRegService{
		ServerId: constants.GatewayServerId,
		SvrType:  constants.GatewayServerType,
	}
	return serve.WriteProtoMessage(constants.MdmGatewayService, constants.SubGatewayServiceRegister, tyutils.HashCode(tyutils.UUID()), &req)
}

func (serve *Server) onMessageBinaryHandler(msg []byte) error {
	pk, err := typacket.NewPacketWithData(msg)
	if err != nil {
		log.Error("非法协议")
		return err
	}

	switch pk.Mid() {
	case constants.MdmGatewayService:
		switch pk.Sid() {
		case constants.SubGatewayServiceRegister:
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
		err = users.Instance().SendClientMessage(pk.ClientId(), ack.Data())
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func (serve *Server) onErrorHandler(w *client.WsClient, err error) {
}
