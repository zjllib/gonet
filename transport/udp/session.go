package udp

import (
	"github.com/flylib/gonet"
	"net"
	"reflect"
	"time"
)

// Socket会话
type session struct {
	gonet.SessionCommon

	remoteAddr             *net.UDPAddr
	serverConn, remoteConn *net.UDPConn
	uuid                   string
	heartbeatTime          time.Time //最近心跳时间点
	nexCheckTime           time.Time //下次检查时间点
}

// 新会话
func newSession(conn *net.UDPConn, remote *net.UDPAddr) *session {
	is := gonet.GetSessionManager().GetIdleSession()
	ns := is.(*session)
	ns.serverConn = conn
	ns.remoteAddr = remote
	return ns
}

func (s *session) RemoteAddr() net.Addr {
	return s.remoteAddr
}

// 发送封包
func (s *session) Send(msgID uint32, msg any) error {
	data, err := gonet.GetNetPackager().Package(msgID, msg)
	if err != nil {
		return err
	}
	if s.remoteConn != nil {
		_, err = s.remoteConn.Write(data)
	} else {
		_, err = s.serverConn.WriteToUDP(data, s.remoteAddr)
	}

	return err
}

func (s *session) Close() error {
	return s.serverConn.Close()
}

// Loop to read messages
func (s *session) readLoop() {
	var buf = make([]byte, 1024)
	for {
		n, err := s.serverConn.Read(buf)
		if err != nil {
			gonet.GetSessionManager().RecycleSession(s)
			return
		}
		msg, err := gonet.GetNetPackager().UnPackage(s, buf[:n])
		if err != nil {
			gonet.GetEventHandler().OnError(s, err)
			continue
		}
		gonet.GetAsyncRuntime().PushMessage(msg)
	}
}

func SessionType() reflect.Type {
	return reflect.TypeOf(session{})
}
