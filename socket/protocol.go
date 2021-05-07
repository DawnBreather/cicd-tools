package network

type Protocol string
const (
	TCP Protocol = "tcp"
	UDP = "udp"
)

func (p Protocol) S() string{
	return string(p)
}