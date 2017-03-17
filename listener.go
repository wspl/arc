package arc

import (
	"math/rand"
	"net"
	"time"
)

func ListenArc(addr string) (*ArcListener, error) {
	l := new(ArcListener)
	l.listenAddr = NewArcAddr().Set(addr)
	l.socketTCP, _ = net.ListenTCP("tcp", l.listenAddr.TCP)
	l.socketUDP, _ = net.ListenUDP("udp", l.listenAddr.UDP)
	l.sessionSet, _ = createSessionSet()

	l.chanAcceptConn = make(chan *ArcConn)
	l.chanAcceptConnError = make(chan error)

	l.taskLoopTCPAccept = NewLoopTask(l.loopTCPAccept).Start()
	l.taskLoopUDPRead = NewLoopTask(l.loopUDPRead).Start()

	return l, nil
}

type ArcListener struct {
	listenAddr *ArcAddr

	socketTCP *net.TCPListener
	socketUDP *net.UDPConn

	sessionSet *SessionSet

	taskLoopTCPAccept *LoopTask
	taskLoopUDPRead *LoopTask

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

func (l *ArcListener) loopTCPAccept() (bool, error) {
	conn, _ := l.socketTCP.AcceptTCP()
	l.tcpAccept(conn)
	return true, nil
}

func (l *ArcListener) loopUDPRead() (bool, error)  {
	buf := make([]byte, 1480)
	size, addr, _ := l.socketUDP.ReadFromUDP(buf)
	l.inputUDP(buf[:size], NewArcAddr().ParseUDP(addr))
	return true, nil
}

func (l *ArcListener) tcpAccept(conn *net.TCPConn) {
	c := new(ArcConn)
	c.remoteAddr = NewArcAddr().ParseTCP(conn.RemoteAddr())
	c.socketTCP = conn
	c.socketUDP = c.socketUDP
	c.listener = l

	c.session, _ = createSession(c, 0)

	c.taskTCPRead = NewLoopTask(c.loopTCPRead).Start()
	c.taskTCPRead.Start()

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