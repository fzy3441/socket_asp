package socket_asp

import (
	"net"
)

type Listener struct {
	lis net.Listener
}

// 监听端口
func Listen(network, address string) (*Listener, error) {
	listener, err := net.Listen(network, address)
	return &Listener{listener}, err
}

// 等待客户连接
func (obj *Listener) Accept(inter INetIO) *NetIO {
RE:
	conn, err := obj.lis.Accept()
	if err != nil {
		goto RE
	}
	return NewConn(inter, conn)
}
