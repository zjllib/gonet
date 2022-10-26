package ws

import (
	"github.com/gorilla/websocket"
	. "github.com/zjllib/gonet/v3"
	"github.com/zjllib/gonet/v3/transport"
	"log"
	"net"
)

// webSocket conn
type session struct {
	SessionIdentify
	SessionStore
	conn *websocket.Conn
}

//新会话
func newSession(conn *websocket.Conn) *session {
	ses := CreateSession()
	newSession, _ := ses.(*session)
	newSession.conn = conn
	CacheMsg(&Message{
		Session: newSession,
		ID:      SessionConnect,
	})
	return newSession
}

func (s *session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *session) Close() error {
	return s.conn.Close()
}

//websocket does not support sending messages concurrently
func (s *session) Send(msg interface{}) error {
	return transport.SendWSPacket(s.conn, msg)
}

//循环读取消息
func (s *session) recvLoop() {
	for {
		_, pkt, err := s.conn.ReadMessage()
		if err != nil {
			RecycleSession(s, err)
			return
		}
		msg, err := transport.ParserWSPacket(pkt)
		if err != nil {
			log.Printf("session_%v msg parser error,reason is %v \n", s.ID(), err)
			continue
		}
		msg.Session = s
		CacheMsg(msg)
	}
}
