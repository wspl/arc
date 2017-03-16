package arc

import (
	"net"
)

func NewArcAddr() *ArcAddr {
	return new(ArcAddr)
}

type ArcAddr struct {
	UDP  *net.UDPAddr
	TCP  *net.TCPAddr
}

func (a *ArcAddr) SetUDP(s string) {
	udpAddr, _ := net.ResolveUDPAddr("udp", s)
	a.UDP = udpAddr
}

func (a *ArcAddr) SetTCP(s string) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", s)
	a.TCP = tcpAddr
}

func (a *ArcAddr) Set(s string) *ArcAddr {
	a.SetTCP(s)
	a.SetUDP(s)
	return a
}

func (a *ArcAddr) Parse(addr interface{}) *ArcAddr {
	switch addr.(type) {
	case net.UDPAddr, *net.UDPAddr:
		a.Set(addr.(*net.UDPAddr).String())
	case net.TCPAddr, *net.TCPAddr:
		a.Set(addr.(*net.TCPAddr).String())
	}
	return a
}

func (a *ArcAddr) ParseTCP(addr interface{}) *ArcAddr {
	switch addr.(type) {
	case net.UDPAddr, *net.UDPAddr:
		a.SetTCP(addr.(*net.UDPAddr).String())
	case net.TCPAddr, *net.TCPAddr:
		a.SetTCP(addr.(*net.TCPAddr).String())
	}
	return a
}

func (a *ArcAddr) ParseUDP(addr interface{}) *ArcAddr {
	switch addr.(type) {
	case net.UDPAddr, *net.UDPAddr:
		a.SetUDP(addr.(*net.UDPAddr).String())
	case net.TCPAddr, *net.TCPAddr:
		a.SetUDP(addr.(*net.TCPAddr).String())
	}
	return a
}

func (a *ArcAddr) String() string {
	return "TCP: " + a.TCP.String() + ", UDP: " + a.UDP.String()
}