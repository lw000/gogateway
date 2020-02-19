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
	c *client.WsClient
}

var (
	serveInstance     *Server
	serveInstanceOnce sync.Once
)

func New() *Server {
	s := &Server{
		c: client.New(client.Config{MaxMessageSize: 1024, MessageBufferSize: 1024}),
	}
	s.init()
	return s
}

func Instance() *Server {
	serveInstanceOnce.Do(func() {
		serveInstance = New()
	})
	return serveInstance
}

func (serve *Server) init() *Server {
	serve.c.HandleMessageBinary(serve.onMessageBinaryHandler)
	serve.c.HandleError(serve.onErrorHandler)
	return serve
}

func (serve *Server) Start() error {
	err := serve.c.Open(global.ProjectConfig.BackendConf.Scheme, global.ProjectConfig.BackendConf.Host, global.ProjectConfig.BackendConf.Path)
	if err != nil {
		log.Error(err)
		return err
	}

	go serve.c.Run()

	err = serve.registerService()
	if err != nil {
		log.Info(err)
		return err
	}

	return nil
}

func (serve *Server) Stop() {
	if serve == nil {
		return
	}
	serve.c.Close()
}

func (serve *Server) WriteProtoMessage(mid, sid uint16, clientId uint32, pb proto.Message) error {
	return serve.c.WriteProtoMessage(mid, sid, clientId, pb)
}

func (serve *Server) WriteBinaryMessage(mid, sid uint16, clientId uint32, data []byte) error {
	return serve.c.WriteBinaryMessage(mid, sid, clientId, data)
}

// 注册服务到路由服务器
func (serve *Server) registerService() error {
	req := Tserve.ReqRegService{ServerId: constant.GATEWAY_SERVER_ID, SvrType: constant.GATEWAY_SERVER_TYPE}
	err := serve.c.WriteProtoMessage(constant.MDM_GATEWAY_SERVICE, constant.SUB_GATEWAY_SERVICE_REGISTER, tyutils.HashCode(tyutils.UUID()), &req)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (serve *Server) onMessageBinaryHandler(msg []byte) error {
	pk, err := typacket.NewPacketWithData(msg)
	if err != nil {
		log.Error("接收到非法协议")
		return err
	}

	switch pk.Mid() {
	case constant.MDM_GATEWAY_SERVICE:
		switch pk.Sid() {
		case constant.SUB_GATEWAY_SERVICE_REGISTER:
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
		}
	}
	return nil
}

func (serve *Server) onErrorHandler(w *client.WsClient, err error) {
}
