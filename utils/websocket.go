package utils

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ontology-oracle/log"
)

type WebSocketOptions struct {
	HeartbeatInterval time.Duration
	HeartbeatPkg      []byte
}

type WebSocketClient struct {
	host              string
	opts              *WebSocketOptions
	conn              *websocket.Conn
	recvCh            chan []byte
	existCh           chan interface{}
	lastHeartbeatTime time.Time
	lock              sync.RWMutex
	status            bool
}

func NewWebSocketClient(host string, opts ...*WebSocketOptions) *WebSocketClient {
	var options *WebSocketOptions
	if len(opts) == 0 {
		options = &WebSocketOptions{
			HeartbeatInterval: 30 * time.Second,
			HeartbeatPkg:      []byte(`{"Action":"heartbeat", "SubscribeEvent":true}`),
		}
	} else {
		options = opts[0]
	}
	return &WebSocketClient{
		host:              host,
		opts:              options,
		recvCh:            make(chan []byte, 1),
		existCh:           make(chan interface{}, 0),
		lastHeartbeatTime: time.Now(),
	}
}

func (ws *WebSocketClient) Connet() (recvCh chan []byte, existCh chan interface{}, err error) {
	ws.conn, _, err = websocket.DefaultDialer.Dial(ws.host, nil)
	if err != nil {
		return nil, nil, err
	}
	ws.status = true
	go ws.doRecv()
	go ws.heartbeat()
	return ws.recvCh, ws.existCh, nil
}

func (ws *WebSocketClient) updateHearbeatTime() {
	ws.lock.Lock()
	defer ws.lock.Unlock()
	ws.lastHeartbeatTime = time.Now()
}

func (ws *WebSocketClient) getHeartbeatTime() time.Time {
	ws.lock.RLock()
	defer ws.lock.RUnlock()
	return ws.lastHeartbeatTime
}

func (ws *WebSocketClient) Send(data []byte) error {
	ws.lock.RLock()
	defer ws.lock.RUnlock()
	if !ws.status {
		return fmt.Errorf("WebSocket connect has already closed")
	}
	return ws.conn.WriteMessage(websocket.TextMessage, data)
}

func (ws *WebSocketClient) doRecv() {
	defer close(ws.recvCh)
	for {
		_, data, err := ws.conn.ReadMessage()
		if err != nil {
			if ws.Status() {
				log.Errorf("WebSocketClient host:%v ReadMessage error:%v", ws.host, err.Error())
				ws.Close()
			}
			return
		}
		ws.updateHearbeatTime()
		ws.recvCh <- data
	}
}

func (ws *WebSocketClient) heartbeat() {
	err := ws.Send(ws.opts.HeartbeatPkg)
	if err != nil {
		log.Errorf("WebSocketClient send heartbeat error:%v", err.Error())
	}
	timer := time.NewTimer(ws.opts.HeartbeatInterval)
	defer timer.Stop()
	for {
		select {
		case <-ws.existCh:
			return
		case <-timer.C:
			err := ws.Send(ws.opts.HeartbeatPkg)
			if err != nil {
				log.Errorf("WebSocketClient send heartbeat error:%v", err.Error())
			}
			timer.Reset(ws.opts.HeartbeatInterval)
		}
	}
}

func (ws *WebSocketClient) Status() bool {
	ws.lock.RLock()
	defer ws.lock.RUnlock()
	return ws.status
}

func (ws *WebSocketClient) Close() {
	ws.lock.Lock()
	defer ws.lock.Unlock()

	if !ws.status {
		return
	}
	ws.status = false
	close(ws.existCh)
	err := ws.conn.Close()
	if err != nil {
		log.Errorf("WebSocketClient close error:%v", err.Error())
	}
}
