package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/lw000/gocommon/network/ws/packet"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	TlsDialer = &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		TLSClientConfig:  &tls.Config{InsecureSkipVerify: true},
	}

	DefaultDialer = &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
	}
)

var (
	DefaultConfig = Config{
		MaxMessageSize:    1024,
		MessageBufferSize: 1024,
	}
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

type WsClient struct {
	sync.RWMutex
	scheme, host, path string
	config             Config
	open               bool
	conn               *websocket.Conn
	done               chan struct{}
	output             chan *envelope
	onMessage          func(msg []byte) error
	onMessageBinary    func(msg []byte) error
	onMessageError     func(w *WsClient, err error)
}

func New(config Config) *WsClient {
	c := &WsClient{
		config: config,
		done:   make(chan struct{}),
	}
	c.init()
	return c
}

func (w *WsClient) init() {
	w.output = make(chan *envelope, w.config.MessageBufferSize)
}

func (w *WsClient) Open(scheme string, host, path string) error {
	w.scheme = scheme
	w.host = host
	w.path = path

	if err := w.reconnect(); err != nil {
		return err
	}

	if w.open {
		go w.writePump()
		go w.checking()
	}

	return nil
}

func (w *WsClient) Closed() bool {
	w.Lock()
	defer w.Unlock()
	return !w.open
}

func (w *WsClient) HandleMessage(fn func(msg []byte) error) {
	w.onMessage = fn
}

func (w *WsClient) HandleMessageBinary(fn func(msg []byte) error) {
	w.onMessageBinary = fn
}

func (w *WsClient) HandleError(fn func(w *WsClient, e error)) {
	w.onMessageError = fn
}

func (w *WsClient) reconnect() error {
	u := url.URL{Scheme: w.scheme, Host: w.host, Path: w.path}

	log.Info("connecting to ", u.String())

	var (
		err  error
		resp *http.Response
	)

	if w.scheme == "wss" {
		w.conn, resp, err = TlsDialer.Dial(u.String(), nil)
	} else if w.scheme == "ws" {
		w.conn, resp, err = DefaultDialer.Dial(u.String(), nil)
	} else {
		return errors.New(fmt.Sprintf("scheme:%s error", w.scheme))
	}
	if err != nil {
		log.Error(err)
		return err
	}

	if resp != nil {
		// log.Info("%+v", resp)
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

func (w *WsClient) checking() {
	defer func() {
		if x := recover(); x != nil {
			log.Error(x)
		}
		log.Error("ws checking exit")
	}()

	ticker := time.NewTicker(time.Second * time.Duration(5))
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if !w.open {
				err := w.Open(w.scheme, w.host, w.path)
				if err != nil {
					log.Error(err)
				}
			}
		case <-w.done:
			return
		}
	}
}

func (w *WsClient) writePump() {
	defer func() {
		if x := recover(); x != nil {
			log.Error(x)
		}
		log.Error("ws writePump exit")
	}()

loop:
	for {
		select {
		case msg, ok := <-w.output:
			if !ok {
				break loop
			}

			err := w.conn.WriteMessage(msg.mt, msg.msg)
			if err != nil {
				w.onMessageError(w, err)
				break loop
			}

			if msg.mt == websocket.CloseMessage {
				break loop
			}
		case <-w.done:
			break loop
		}
	}
}

func (w *WsClient) write(msg *envelope) error {
	if w.Closed() {
		return errors.New("ws is closed")
	}

	w.Lock()
	defer w.Unlock()
	err := w.conn.WriteMessage(msg.mt, msg.msg)
	if err != nil {
		w.open = false
		w.done <- struct{}{}
		return err
	}
	return nil
}

func (w *WsClient) WriteTextMessage(data []byte) error {
	if w.Closed() {
		return errors.New("ws is closed")
	}

	select {
	case w.output <- &envelope{mt: websocket.TextMessage, msg: data}:
	default:
		w.onMessageError(w, errors.New("ws message buffer is full"))
	}
	// log.Info("in: %d", len(w.output))
	return nil
	return w.write(&envelope{mt: websocket.BinaryMessage, msg: data})
}

func (w *WsClient) WriteProtoMessage(mid, sid uint16, clientId uint32, pb proto.Message) error {
	if w.Closed() {
		return errors.New("ws is closed")
	}

	pk := typacket.NewPacket(mid, sid, clientId)
	data, err := proto.Marshal(pb)
	if err != nil {
		return err
	}

	err = pk.Encode(data)
	if err != nil {
		return errors.New("构建数据包错误")
	}

	select {
	case w.output <- &envelope{mt: websocket.BinaryMessage, msg: pk.Data()}:
	default:
		w.onMessageError(w, errors.New("ws message buffer is full"))
	}
	// log.Info("in: %d", len(w.output))
	return nil
	return w.write(&envelope{mt: websocket.BinaryMessage, msg: pk.Data()})
}

func (w *WsClient) WriteBinaryMessage(mid, sid uint16, clientId uint32, data []byte) (err error) {
	if w.Closed() {
		return errors.New("ws is closed")
	}

	pk := typacket.NewPacket(mid, sid, clientId)
	err = pk.Encode(data)
	if err != nil {
		return errors.New("构建数据包错误")
	}

	select {
	case w.output <- &envelope{mt: websocket.BinaryMessage, msg: pk.Data()}:
	default:
		err = errors.New("ws message buffer is full")
		w.onMessageError(w, err)
	}
	// log.Info("in: %d", len(w.output))
	return
	return w.write(&envelope{mt: websocket.BinaryMessage, msg: pk.Data()})
}

func (w *WsClient) ping() error {
	if w.Closed() {
		return errors.New("ws is closed")
	}

	select {
	case w.output <- &envelope{mt: websocket.PingMessage, msg: []byte{}}:
	default:
		w.onMessageError(w, errors.New("ws message buffer is full"))
	}
	return nil
	return w.write(&envelope{websocket.PingMessage, []byte{}})
}

func (w *WsClient) pong() error {
	if w.Closed() {
		return errors.New("ws is closed")
	}

	select {
	case w.output <- &envelope{mt: websocket.PongMessage, msg: []byte{}}:
	default:
		w.onMessageError(w, errors.New("ws message buffer is full"))
	}
	return nil
	return w.write(&envelope{websocket.PongMessage, []byte{}})
}

func (w *WsClient) Run() {
	if w.Closed() {
		log.Error("ws is closed")
		return
	}

	defer func() {
		if x := recover(); x != nil {
			log.Error(x)
		}
	}()

	w.conn.SetReadLimit(w.config.MaxMessageSize)

	for {
		mt, msg, err := w.conn.ReadMessage()
		if err != nil {
			w.open = false
			log.Error(err)
			w.onMessageError(w, err)
			return
		}

		if mt == websocket.TextMessage {
			err = w.onMessage(msg)
			if err != nil {
				log.Error(err)
				return
			}
		}

		if mt == websocket.BinaryMessage {
			err = w.onMessageBinary(msg)
			if err != nil {
				log.Error(err)
				return
			}
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
	if err := w.conn.Close(); err != nil {
		log.Error(err)
	}
	close(w.output)
	close(w.done)
	w.Unlock()
}
