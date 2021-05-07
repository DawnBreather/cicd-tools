package socket

import (
	"fmt"
	"net"
	"time"
)

type Socket struct {
	host string
	port int
	protocol Protocol
}

func (n *Socket) SetHost(host string) *Socket {
	n.host = host
	return n
}

func (n *Socket) SetProtocol(protocol Protocol) *Socket {
	n.protocol = protocol
	return n
}

func (n *Socket) SetPort(port int) *Socket {
	n.port = port
	return n
}

func (n *Socket) GetSocket() string {
	return fmt.Sprintf("%s:%d", n.host, n.port)
}

func (n *Socket) IsPortOpen() bool {
	timeout := time.Second
	conn, err := net.DialTimeout(n.protocol.S(), n.GetSocket(), timeout)
	if err != nil {
		//fmt.Println("Connecting error:", err)
		return false
	}
	if conn != nil {
		defer conn.Close()
		//fmt.Println("Opened", net.JoinHostPort(host, port))
		return true
	}
	return false
}