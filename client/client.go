package client

import (
	"context"
	"net"
	"zombiezen.com/go/capnproto2/rpc"
)

type Client struct {
	ctx     context.Context
	conn    *rpc.Conn
	rawConn net.Conn
	api     capnp.API
}
