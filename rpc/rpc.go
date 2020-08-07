package rpc

import (
	"github.com/fanghongbo/dlog"
	"github.com/fanghongbo/ops-hbs/common/g"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

type Hbs int
type Agent int

func Start() {
	var (
		addr     string
		server   *rpc.Server
		listener net.Listener
		err      error
	)

	addr = g.Conf().Rpc.Listen
	server = rpc.NewServer()

	if err = server.Register(new(Agent)); err != nil {
		dlog.Errorf("register rpc err: %s", err)
	}
	if err = server.Register(new(Hbs)); err != nil {
		dlog.Errorf("register hbs err: %s", err)
	}

	listener, err = net.Listen("tcp", addr)
	if err != nil {
		dlog.Errorf("listen error: %s", err)
	} else {
		dlog.Infof("listening %s", addr)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			dlog.Errorf("listener accept err:", err)
			time.Sleep(time.Duration(100) * time.Millisecond)
			continue
		}

		go server.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
