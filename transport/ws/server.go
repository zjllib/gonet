package ws

import (
	"github.com/gorilla/websocket"
	. "github.com/zjllib/gonet/v3"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

var _ IServer = new(server)

// 接收端
type server struct {
	PeerIdentify
	//指定将HTTP连接升级到WebSocket连接的参数。
	upGrader websocket.Upgrader
	//响应头
	//respHeader http.Header
}

func NewServer(ctx *Context) IServer {
	s := &server{}
	s.WithContext(ctx)
	ctx.InitSessionMgr(reflect.TypeOf(session{}))
	return s
}

func (s *server) Listen(addr string) error {
	s.SetAddr(addr)
	url, err := url.Parse(s.Addr())
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.HandleFunc(url.Path, s.newConn)
	return http.ListenAndServe(url.Host, mux)
}

func (s *server) Stop() error {
	s.upGrader.HandshakeTimeout = time.Nanosecond
	return nil
}

func (s *server) newConn(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("Connection", "upgrade") //升级
	r.Header.Add("Upgrade", "websocket")  //websocket

	conn, err := s.upGrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	go newSession(s.Context, conn).readLoop()
}
