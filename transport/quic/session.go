package quic

import (
	"context"
	"github.com/quic-go/quic-go"
	. "github.com/zjllib/gonet/v3"
	"log"
	"net"
)

// conn
type session struct {
	SessionIdentify
	SessionStore
	conn   quic.Connection
	stream quic.Stream
}

// 新会话
func newSession(c *Context, conn quic.Connection) *session {
	ses := c.CreateSession()
	s, _ := ses.(*session)
	s.conn = conn
	s.WithContext(c)
	return s
}

func (s *session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *session) Send(msg any) error {
	data, err := s.Context.Package(msg)
	if err != nil {
		return err
	}
	_, err = s.stream.Write(data)
	return err
}

func (s *session) Close() error {
	err := s.conn.CloseWithError(0, "EOF")
	s.conn = nil
	return err
}

// 循环读取消息
func (s *session) recvLoop() {
	var err error
	s.stream, err = s.conn.AcceptStream(context.Background())
	if err != nil {
		log.Printf("session_%v AcceptStream error,reason is %v \n", s.ID(), err)
		s.Context.RecycleSession(s, err)
		return
	}

	for {
		var n int
		buf := make([]byte, 1024)
		n, err = s.stream.Read(buf)
		if err != nil {
			log.Printf("session_%v reading error,reason is %v \n", s.ID(), err)
			err = s.stream.Close()
			if err != nil {
				log.Printf("session_%v close error,reason is %v \n", s.ID(), err)
			}
			s.Context.RecycleSession(s, err)
			return
		}
		msg, _, err := s.Context.UnPackage(buf[:n])
		if err != nil {
			log.Printf("session_%v msg parser error,reason is %v \n", s.ID(), err)
			continue
		}
		s.Context.PushGlobalMessageQueue(s, msg)
	}
}