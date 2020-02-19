package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/lw000/gocommon/network/ws/packet"
	"github.com/olahol/melody"
	log "github.com/sirupsen/logrus"
	"gogateway/constant"
	"gogateway/wsserve/config"
	"gogateway/wsserve/control"
	"net/http"
	"time"
)

var (
	m    *melody.Melody
	addr = flag.String("addr", ":8830", "http service address")
)

//
// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  10240,
// 	WriteBufferSize: 10240,
// 	CheckOrigin:     func(r *http.Request) bool { return true },
// } // use default options
//
// func handleConnect(conn *websocket.Conn) {
// 	defer func() {
// 		log.Info("%s exit", conn.RemoteAddr().String())
// 		er := conn.Close()
// 		if er != nil {
// 			log.Error("%v", er)
// 		}
// 	}()
//
// 	for {
// 		mt, message, er := conn.ReadMessage()
// 		if er != nil {
// 			log.Error(er)
// 			break
// 		}
// 		switch mt {
// 		case websocket.PingMessage:
// 			er = conn.WriteMessage(websocket.PongMessage, []byte{})
// 			if er != nil {
// 				break
// 			}
// 		case websocket.PongMessage:
// 			er = conn.WriteMessage(websocket.PingMessage, []byte{})
// 			if er != nil {
// 				break
// 			}
// 		case websocket.BinaryMessage:
// 			er = control.GetHub().DispatchMessage(conn, message)
// 			if er != nil {
// 				break
// 			}
// 		case websocket.TextMessage:
// 			er = conn.WriteMessage(websocket.TextMessage, message)
// 			if er != nil {
// 				break
// 			}
// 		default:
// 			er = conn.WriteMessage(mt, message)
// 			if er != nil {
// 				break
// 			}
// 		}
// 	}
// }
//
// func wsHandler(w http.ResponseWriter, r *http.Request) {
// 	conn, er := upgrader.Upgrade(w, r, w.Header())
// 	if er != nil {
// 		log.Error("upgrade: %s", er.Error())
// 		return
// 	}
// 	go handleConnect(conn)
// }

func wsHandler(w http.ResponseWriter, r *http.Request) {
	er := m.HandleRequest(w, r)
	if er != nil {
		er = m.CloseWithMsg([]byte(er.Error()))
	}
}

func main() {
	flag.Parse()

	config.ConfigLocalFilesystemLogger("log", "wsserve", time.Hour*24*365, time.Hour*24)

	engine := gin.Default()

	m = melody.New()
	m.Config = &melody.Config{
		WriteWait:         10 * time.Second,
		PongWait:          30 * time.Second,
		PingPeriod:        (30 * time.Second * 9) / 10,
		MaxMessageSize:    1024,
		MessageBufferSize: 1024,
	}

	engine.GET("/ws", func(c *gin.Context) {
		wsHandler(c.Writer, c.Request)
	})

	m.HandleMessageBinary(func(s *melody.Session, msg []byte) {
		pk, er := typacket.NewPacketWithData(msg)
		if er != nil {
			return
		}

		mid := pk.Mid()
		sid := pk.Sid()
		switch mid {
		case constant.MDM_GATEWAY_SERVICE:
			switch sid {
			case constant.SUB_GATEWAY_SERVICE_REGISTER:
				control.OnRegisterService(s, pk)
			default:

			}
		case constant.MDM_CLIENT:
			switch sid {
			case constant.SUB_CLIENT_MSG:
				control.OnMessage(s, pk)
			default:

			}
		default:

		}
	})

	m.HandleConnect(func(s *melody.Session) {
		log.WithField("RemoteAddr", s.Request.RemoteAddr).Info("客户端·连接")
	})

	m.HandleDisconnect(func(s *melody.Session) {
		log.WithField("RemoteAddr", s.Request.RemoteAddr).Info("客户端·断开")
	})

	m.HandleError(func(s *melody.Session, e error) {
		log.WithField("RemoteAddr", s.Request.RemoteAddr).Info(e.Error())
	})
	log.Error(engine.Run(*addr))

	// http.HandleFunc("/", wsHandler)
	// log.Info(*addr)
	// log.Error(http.ListenAndServe(*addr, nil))
}
