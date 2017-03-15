package arc

import "time"

func createSession(conn *ArcConn, id uint32) (*Session, error) {
	s := new(Session)
	s.conn = conn
	s.Id = id
	return s, nil
}

type Session struct {
	conn     *ArcConn
	Id       uint32
	Accepted bool
}

func (s *Session) HasId() bool {
	return s.Id != 0
}

func (s *Session) HandleSessionSegment(b *[]byte) {
	switch ReadType(b) {
	case UDP_SEGMENT_SESSION_PROBE:
		seg, _ := ParseUDPSegmentSessionCommand(b)
		s.HandleSessionProbeSegment(seg)
	case TCP_SEGMENT_SESSION_ACCEPT:
		s.Accepted = true
	}
}

func (s *Session) HandleSessionProbeSegment(seg *UDPSegmentSessionCommand) {
	if !s.Accepted {
		seg := &TCPSegmentSessionCommand{
			Type: TCP_SEGMENT_SESSION_ACCEPT,
			SessionId: s.Id,
		}
		s.conn.outputTCP(seg.Binary())
		s.Accepted = true
		s.conn.listener.chanAcceptConn <- s.conn
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
	go func() {
		for !s.Accepted {
			seg := &UDPSegmentSessionCommand{
				Type: UDP_SEGMENT_SESSION_PROBE,
				SessionId: s.Id,
			}
			s.conn.outputUDP(seg.Binary())
			time.Sleep(100 * time.Millisecond)
		}
	}()
}