package control

import (
	"github.com/golang/protobuf/proto"
	"github.com/lw000/gocommon/network/ws/packet"
	"github.com/olahol/melody"
	log "github.com/sirupsen/logrus"
	"gogateway/protos/msg"
	"gogateway/protos/serve"
)

// var (
// 	hub *tyhub.Hub
// )

func init() {
	// registerHub()
}

func AckMessage(s *melody.Session, data []byte) {
	if err := s.WriteBinary(data); err != nil {
		log.Error(err)
	}
}

func OnRegisterService(s *melody.Session, pk *typacket.Packet) {
	var req Tserve.ReqRegService
	if err := proto.Unmarshal(pk.Data(), &req); err != nil {
		log.Error(err)
		return
	}
	log.Infof("[%s] %+v req:%+v", s.Request.RemoteAddr, pk, req)

	data, err := proto.Marshal(&Tserve.AckRegService{Result: 0, Errmsg: ""})
	if err != nil {
		log.Error(err)
		return
	}
	ack := typacket.NewPacket(pk.Mid(), pk.Sid(), pk.ClientId())
	if err = ack.Encode(data); err == nil {
		AckMessage(s, ack.Data())
	}
}

func OnMessage(s *melody.Session, pk *typacket.Packet) {
	var req Tmsg.ReqTestMessage
	if err := proto.Unmarshal(pk.Data(), &req); err != nil {
		log.Error(err)
		return
	}
	log.Infof("[%s] %+v req:%+v", s.Request.RemoteAddr, pk, req)

	data, err := proto.Marshal(&Tmsg.AckTestMessage{Code: req.GetUid(), Msg: req.GetMsg()})
	if err != nil {
		log.Error(err)
		return
	}
	ack := typacket.NewPacket(pk.Mid(), pk.Sid(), pk.ClientId())
	if err = ack.Encode(data); err == nil {
		AckMessage(s, ack.Data())
	}
}

// func registerHub() {
// 	hub := tyhub.NewHub()
// 	hub.AddHandler(constant.MDM_GATEWAY_SERVICE, constant.SUB_GATEWAY_SERVICE_REGISTER, &tyhub.Handler{Fn:onRegisgerService})
// 	hub.AddHandler(constant.MDM_CLIENT, constant.SUB_CLIENT_MSG, &tyhub.Handler{Fn:onMessage})
// }
//
//
// func GetHub() *tyhub.Hub {
// 	return hub
// }
