package arc

import (
	"time"
	"math/rand"
)

func createSession(conn *ArcConn, id uint32) (*Session, error) {
	s := new(Session)
	s.conn = conn
	s.Id = id

	s.chanProbe = make(chan bool)
	s.chanKeepAlive = make(chan bool)

	return s, nil
}

type Session struct {
	conn     *ArcConn
	Id       uint32
	Accepted bool

	chanProbe     chan bool
	chanKeepAlive chan bool
}

func (s *Session) loopKeepAlive() {
	for {
		randSeconds := rand.Intn(25) + 30
		time.Sleep(time.Duration(randSeconds) * time.Second)
		s.KeepAlive()
		println("keep aliving")
	}
}

func (s *Session) HasId() bool {
	return s.Id != 0
}

func (s *Session) HandleSessionSegment(b *[]byte, src *ArcAddr) {
	switch ReadType(b) {
	case UDP_SEGMENT_SESSION_PROBE:
		s.HandleSessionProbe(src)
	case TCP_SEGMENT_SESSION_PROBE_ACK:
		s.HandleSessionProbeAck(src)
	case UDP_SEGMENT_SESSION_KEEP_ALIVE:
		s.HandleSessionKeepAlive()
	case TCP_SEGMENT_SESSION_KEEP_ALIVE_ACK:
		s.HandleSessionKeepAliveAck()
	}
}

func (s *Session) HandleSessionProbe(src *ArcAddr) {
	s.conn.remoteAddr.ParseUDP(src.UDP)
	seg := &TCPSegmentSessionCommand{
		Type:      TCP_SEGMENT_SESSION_PROBE_ACK,
		SessionId: s.Id,
	}
	s.conn.outputTCP(seg.Binary())

	if !s.Accepted {
		s.Accepted = true
		s.conn.listener.chanAcceptConn <- s.conn
	}
}

func (s *Session) HandleSessionKeepAlive() {
	seg := &TCPSegmentSessionCommand{
		Type:      TCP_SEGMENT_SESSION_KEEP_ALIVE_ACK,
		SessionId: s.Id,
	}
	s.conn.outputTCP(seg.Binary())
}

func (s *Session) HandleSessionKeepAliveAck() {
	s.chanKeepAlive <- true
}

func (s *Session) HandleSessionProbeAck(src *ArcAddr) {
	s.chanProbe <- true
	if !s.Accepted {
		s.Accepted = true
		s.conn.chanAccepted <- true
		go s.loopKeepAlive()
	}
}

func (s *Session) Make() {
	seg := &TCPSegmentSessionCommand{
		Type: TCP_SEGMENT_SESSION_NEW,
		SessionId: s.Id,
	}
	s.conn.outputTCP(seg.Binary())
}

func (s *Session) Respond() {
	seg := &TCPSegmentSessionCommand{
		Type: TCP_SEGMENT_SESSION_RESPONSE,
		SessionId: s.Id,
	}
	s.conn.outputTCP(seg.Binary())
}

func (s *Session) Probe() {
	sending := true
	go func() {
		for sending {
			seg := &UDPSegmentSessionCommand{
				Type: UDP_SEGMENT_SESSION_PROBE,
				SessionId: s.Id,
			}
			s.conn.outputUDP(seg.Binary())
			time.Sleep(100 * time.Millisecond)
		}
	}()
	<- s.chanProbe
	sending = false
}

func (s *Session) ProbeAsync() {
	go s.Probe()
}

func (s *Session) KeepAlive() {
	sending := true
	go func() {
		for sending {
			seg := &UDPSegmentSessionCommand{
				Type: UDP_SEGMENT_SESSION_KEEP_ALIVE,
				SessionId: s.Id,
			}
			s.conn.outputUDP(seg.Binary())
			time.Sleep(100 * time.Millisecond)
		}
	}()
	<- s.chanKeepAlive
	sending = false
}