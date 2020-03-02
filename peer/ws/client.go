package ws

import (
	"github.com/gorilla/websocket"
	"goNet"
	. "goNet/log"
	"net/http"
	"time"
)

type client struct {
	goNet.PeerIdentify
	session *session
}

func init() {
	identify := goNet.PeerIdentify{}
	identify.SetType(goNet.PEER_CLIENT)
	c := &client{
		PeerIdentify: identify,
	}
	goNet.RegisterPeer(c)
}

func (c *client) Start() {
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 5 * time.Second,
	}
	conn, _, err := dialer.Dial(c.Addr(), nil)
	if err != nil {
		Log.Errorf("#ws.connect failed(%s) %v", c.Addr(), err.Error())
		return
	}
	Log.Info(conn.RemoteAddr())
	c.session = newSession(conn)
	go c.session.recvLoop()
}

func (c *client) Stop() {
	c.session.conn.SetReadDeadline(time.Now())
}
