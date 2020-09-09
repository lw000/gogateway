package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/lw000/gocommon/network/ws/hub"
	"github.com/lw000/gocommon/network/ws/packet"
	"github.com/lw000/gocommon/utils"
	log "github.com/sirupsen/logrus"
	"gogateway/constants"
	"gogateway/protos/msg"
	"gogateway/test/client"
	"gogateway/test/config"
	"math/rand"
	"time"
)

var (
	cfg *config.JsonConfig
)

type EventHandler struct {
	tyhub.Handler
}

func (ev *EventHandler) Receiver(conn *websocket.Conn, pk *typacket.Packet) {
	var ack Tmsg.AckTestMessage
	if err := proto.Unmarshal(pk.Data(), &ack); err != nil {
		log.Println(err)
		return
	}
	log.Printf("接收：uid:%d, message:%s", ack.GetCode(), ack.GetMsg())
}

func TestMessage(c *client.WsClient, uid uint32) {
	t := time.NewTicker(time.Millisecond * time.Duration(cfg.Millisecond))
	defer t.Stop()
	for {
		select {
		case <-t.C:
			{
				req := Tmsg.ReqTestMessage{Uid: uid, Msg: "test message 1"}
				err := c.WriteProtoMessage(constants.MdmClient, constants.SubClientMsg, &req, &EventHandler{})
				if err != nil {
					log.Println(err)
					return
				}
			}

			{
				req := Tmsg.ReqTestMessage{Uid: uid, Msg: "test message 2"}
				err := c.WriteProtoMessage(constants.MdmClient, constants.SubClientMsg, &req, &EventHandler{})
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

func main() {
	config.ConfigLocalFilesystemLogger("log", "client", time.Hour*24*365, time.Hour*24)

	var err error
	cfg, err = config.LoadJsonConfig("conf.json")
	if err != nil {
		log.Panic(err)
	}

	log.Printf("%+v", cfg)

	rand.Seed(time.Now().Unix())

	for i := 1; i <= cfg.Count; i++ {
		cli := client.New()
		cli.HandleError(func(w *client.WsClient, e error) {
			log.Println("disconnected")
		})

		err = cli.Open(cfg.WsConf.Scheme, cfg.WsConf.Host, cfg.WsConf.Path)
		if err != nil {
			log.Println(err)
			continue
		}

		uid := tyutils.HashCode(tyutils.UUID())

		log.Printf("connected. [uid=%d]", uid)

		go cli.Run()

		if cfg.Send {
			go TestMessage(cli, uint32(uid))
		}

		time.Sleep(time.Microsecond * time.Duration(100))
	}
	select {}
}
