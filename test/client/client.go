package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/lw000/gocommon/network/ws/hub"
	"github.com/lw000/gocommon/network/ws/packet"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type envelope struct {
	mt  int
	msg []byte
}

// Config configuration struct.
type Config struct {
	MaxMessageSize    int64 // Maximum size in bytes of a message.
	MessageBufferSize int   // The max amount of messages that can be in a sessions buffer before it starts dropping them.
}

var (
	DefaultDialer = &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
	}

	TlsDialer = &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
	}

	DefaultConfig = &Config{
		MaxMessageSize:    1024,
		MessageBufferSize: 1024,
	}
)

type WsClient struct {
	sync.RWMutex
	config         *Config
	open           bool
	conn           *websocket.Conn
	hub            *tyhub.Hub
	onMessageError func(w *WsClient, e error)
}

func New() *WsClient {
	return &WsClient{
		config: DefaultConfig,
		hub:    tyhub.New(),
	}
}

func NewWithConfig(config *Config) *WsClient {
	if config == nil {
		config = DefaultConfig
	}
	return &WsClient{
		config: config,
		hub:    tyhub.New(),
	}
}

func (w *WsClient) Open(scheme string, host, path string) error {
	u := url.URL{Scheme: scheme, Host: host, Path: path}

	log.Info("connecting to ", u.String())

	var (
		err  error
		resp *http.Response
	)

	if scheme == "wss" {
		w.conn, resp, err = TlsDialer.Dial(u.String(), nil)
	} else if scheme == "ws" {
		w.conn, resp, err = DefaultDialer.Dial(u.String(), nil)
	} else {
		return errors.New(fmt.Sprintf("未知Scheme:%s", scheme))
	}

	if err != nil {
		log.Error(err)
		return err
	}

	if resp != nil {
		// log.Info(fmt.Sprintf("%+v", resp))
	}

	w.conn.SetPingHandler(func(appData string) error {
		return w.pong()
	})

	w.conn.SetPongHandler(func(appData string) error {
		return w.ping()
	})

	w.open = true

	return nil
}

func (w *WsClient) Closed() bool {
	w.Lock()
	defer w.Unlock()
	return !w.open
}

func (w *WsClient) HandleError(fn func(w *WsClient, e error)) {
	w.onMessageError = fn
}

func (w *WsClient) write(msg *envelope) error {
	if w.Closed() {
		return errors.New("ws is closed")
	}

	w.Lock()
	defer w.Unlock()

	err := w.conn.WriteMessage(msg.mt, msg.msg)
	if err != nil {
		w.onMessageError(w, err)
		return err
	}
	return nil
}

func (w *WsClient) WriteTextMessage(mid, sid uint16, data []byte, handler tyhub.Handler) error {
	if w.Closed() {
		return errors.New("ws is closed")
	}

	pk := typacket.NewPacket(mid, sid, 0)
	if len(data) > 0 {
		err := pk.Encode(data)
		if err != nil {
			return errors.New("构建数据包错误")
		}
	}

	w.hub.RegisterHandler(mid, sid, handler)

	return w.write(&envelope{mt: websocket.BinaryMessage, msg: pk.Data()})
}

func (w *WsClient) WriteProtoMessage(mid, sid uint16, pb proto.Message, handler tyhub.Handler) error {
	if w.Closed() {
		return errors.New("ws is closed")
	}

	data, err := proto.Marshal(pb)
	if err != nil {
		return err
	}

	pk := typacket.NewPacket(mid, sid, 0)
	err = pk.Encode(data)
	if err != nil {
		return errors.New("构建数据包错误")
	}

	w.hub.RegisterHandler(mid, sid, handler)

	return w.write(&envelope{mt: websocket.BinaryMessage, msg: pk.Data()})
}

func (w *WsClient) WriteBinaryMessage(mid, sid uint16, data []byte, handler tyhub.Handler) error {
	if w.Closed() {
		return errors.New("ws is closed")
	}

	pk := typacket.NewPacket(mid, sid, 0)

	if len(data) > 0 {
		err := pk.Encode(data)
		if err != nil {
			return errors.New("构建数据包错误")
		}
	}

	w.hub.RegisterHandler(mid, sid, handler)

	return w.write(&envelope{mt: websocket.BinaryMessage, msg: pk.Data()})
}

func (w *WsClient) ping() error {
	if w.Closed() {
		return errors.New("ws is closed")
	}
	// log.Info("PingMessage")
	return w.write(&envelope{mt: websocket.PingMessage, msg: []byte{}})
}

func (w *WsClient) pong() error {
	if w.Closed() {
		return errors.New("ws is closed")
	}
	// log.Info("PongMessage")
	return w.write(&envelope{mt: websocket.PongMessage, msg: []byte{}})
}

func (w *WsClient) Run() {
	defer func() {
		log.Error("ws ReadMessage exit")
	}()

	if w.Closed() {
		log.Error(errors.New("ws is closed"))
		return
	}

	w.conn.SetReadLimit(w.config.MaxMessageSize)

	for {
		mt, msg, err := w.conn.ReadMessage()
		if err != nil {
			w.open = false
			log.Error(err)
			w.onMessageError(w, err)
			return
		}

		if mt != websocket.BinaryMessage {
			return
		}

		if err = w.hub.DispatchMessage(w.conn, msg); err != nil {
			log.Error(err)
			return
		}
	}
}

func (w *WsClient) Close() {
	if w == nil {
		return
	}

	if w.Closed() {
		return
	}
	w.Lock()
	w.open = false
	_ = w.conn.Close()
	w.Unlock()
}
