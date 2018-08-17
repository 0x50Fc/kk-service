package kk

import (
	"errors"

	"github.com/gorilla/websocket"
)

const (
	WSChannelReadClosed  = 1
	WSChannelWriteClosed = 2
)

type WSChannel struct {
	conn *websocket.Conn
}

func NewWSChannel(conn *websocket.Conn) *WSChannel {
	v := WSChannel{}
	v.conn = conn
	return &v
}

/**
 * 发送消息
 */
func (C *WSChannel) Send(data []byte) error {
	if C.conn != nil {
		return C.conn.WriteMessage(websocket.TextMessage, data)
	}
	return errors.New("Not Found Connection")
}

func (C *WSChannel) Close() {
	if C.conn != nil {
		conn := C.conn
		C.conn = nil
		conn.Close()
	}
}
