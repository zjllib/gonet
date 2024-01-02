package demo

import (
	"demo/handler"
	"demo/proto"
	"fmt"
	"github.com/flylib/gonet"
	transport "github.com/flylib/gonet/transport/tcp"
	"github.com/flylib/goutils/codec/json"
	"github.com/flylib/pkg/log/builtinlog"
	"log"
	"testing"
	"time"
)

func TestTcpServer(t *testing.T) {
	gonet.SetContext(
		gonet.WithEventHandler(handler.EventHandler{}),
		gonet.WithNetPackager(gonet.TcpNetPackager{}),

		gonet.MustWithSessionType(transport.SessionType()),
		gonet.MustWithCodec(&json.Codec{}),
		gonet.MustWithLogger(builtinlog.NewLogger()),
	)
	fmt.Println("server listen on ws://localhost:8088/center/ws")
	if err := transport.NewServer().Listen("ws://localhost:8088/center/ws"); err != nil {
		log.Fatal(err)
	}
}

func TestTcpClient(t *testing.T) {
	gonet.SetContext(
		gonet.WithEventHandler(handler.EventHandler{}),
		gonet.WithNetPackager(gonet.TcpNetPackager{}),

		gonet.MustWithSessionType(transport.SessionType()),
		gonet.MustWithCodec(&json.Codec{}),
		gonet.MustWithLogger(builtinlog.NewLogger()),
	)
	session, err := transport.NewClient(transport.WithHandshakeTimeout(5 * time.Second)).Dial("ws://localhost:8088/center/ws")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connect success")

	tick := time.Tick(time.Second * 1)
	var i int
	for range tick {
		//fmt.Println("send msg", i)
		i++
		err = session.Send(101, &proto.Say{
			fmt.Sprintf("hello server %d", i),
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
