package arc

import (
	"net"
	"time"
	"math/rand"
)

func ListenArc(addr string) (*ArcListener, error) {
	l := new(ArcListener)
	l.listenAddr = NewArcAddr().Set(addr)
	l.socketTCP, _ = net.ListenTCP("tcp", l.listenAddr.TCP)
	l.socketUDP, _ = net.ListenUDP("udp", l.listenAddr.UDP)
	l.sessionSet, _ = createSessionSet()

	l.chanAcceptConn = make(chan *ArcConn)
	l.chanAcceptConnError = make(chan error)

	l.enabledLoopTCPAccept = true
	l.enabledLoopUDPRead = true
	go l.loopTCPAccept()
	go l.loopUDPRead()

	return l, nil
}

type ArcListener struct {
	listenAddr *ArcAddr

	socketTCP *net.TCPListener
	socketUDP *net.UDPConn

	sessionSet *SessionSet

	enabledLoopTCPAccept bool
	enabledLoopUDPRead bool

	chanAcceptConn      chan *ArcConn
	chanAcceptConnError chan error
}

func (l *ArcListener) AcceptArc() (*ArcConn, error) {
	select {
	case conn := <- l.chanAcceptConn:
		return conn, nil
	case err := <- l.chanAcceptConnError:
		return nil, err
	}
}

func (l *ArcListener) loopTCPAccept() {
	for l.enabledLoopTCPAccept {
		conn, _ := l.socketTCP.AcceptTCP()
		l.tcpAccept(conn)
	}
}

func (l *ArcListener) loopUDPRead() {
	for l.enabledLoopUDPRead {
		buf := make([]byte, 1480)
		size, addr, _ := l.socketUDP.ReadFromUDP(buf)
		l.inputUDP(buf[:size], NewArcAddr().ParseUDP(addr))
	}
}

func (l *ArcListener) tcpAccept(conn *net.TCPConn) {
	c := new(ArcConn)
	c.remoteAddr = NewArcAddr().ParseTCP(conn.RemoteAddr())
	c.socketTCP = conn
	c.socketUDP = c.socketUDP
	c.listener = l

	c.session, _ = createSession(c, 0)

	c.enabledLoopTCPRead = true
	c.startLoopTCPRead()

	if ARC_DEBUG_SIMULATION_MODE {
		go func() {
			// Automatic Close TCP Connection
			for {
				randSeconds := rand.Intn(5) + 5
				time.Sleep(time.Duration(randSeconds) * time.Second)
				c.socketTCP.Close()
			}
		}()
	}
}

func (l *ArcListener) inputUDP(b []byte, src *ArcAddr) {
	l.handleUDPSegment(&b, src)
}

func (l *ArcListener) handleUDPSegment(b *[]byte, src *ArcAddr) {
	switch SegmentClass(ReadType(b)) {
	case SEGMENT_CLASS_SESSION:
		s, _ := ParseTCPSegmentSessionCommand(b)
		session := l.sessionSet.Get(s.SessionId)
		if session == nil { return }
		session.HandleSessionSegment(b, src)
	}
}