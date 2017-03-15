package arc

import (
	"errors"
	"net"
)

func NewArcAddr(s string) (*ArcAddr, error) {
	a := new(ArcAddr)

	udpAddr, err := net.ResolveUDPAddr("udp", s)
	if err != nil {
		return nil, err
	}
	a.UDP = udpAddr

	tcpAddr, err := net.ResolveTCPAddr("tcp", s)
	if err != nil {
		return nil, err
	}
	a.TCP = tcpAddr

	a.IP = a.TCP.IP
	a.Port = a.TCP.Port

	return a, nil
}

func ArcAddrParse(addr interface{}) (*ArcAddr, error) {
	switch addr.(type) {
	case net.UDPAddr, *net.UDPAddr:
		return NewArcAddr(addr.(*net.UDPAddr).String())
	case net.TCPAddr, *net.TCPAddr:
		return NewArcAddr(addr.(*net.TCPAddr).String())
	default:
		return nil, errors.New("Unknown address type")
	}
}

type ArcAddr struct {
	UDP  *net.UDPAddr
	TCP  *net.TCPAddr
	IP   net.IP
	Port int
}

func (a *ArcAddr) String() string {
	return a.TCP.String()
}