package arc

import (
	"bufio"
	"net"
)

func DialArc(addr string) (*ArcConn, error) {
	c := new(ArcConn)
	c.remoteAddr = NewArcAddr().Set(addr)
	c.socketTCP, _ = net.DialTCP("tcp", nil, c.remoteAddr.TCP)
	c.socketUDP, _ = net.ListenUDP("udp", nil)

	c.session, _ = createSession(c, 0)

	c.chanAccepted = make(chan bool)

	c.taskTCPRead = NewLoopTask(c.loopTCPRead).Start()
	c.taskUDPRead = NewLoopTask(c.loopUDPRead).Start()

	c.session.Make()

	<- c.chanAccepted

	return c, nil
}

type ArcConn struct {
	remoteAddr *ArcAddr

	listener *ArcListener

	socketTCP *net.TCPConn
	socketUDP *net.UDPConn

	session *Session

	taskTCPRead *LoopTask
	taskUDPRead *LoopTask

	chanAccepted chan bool
}

func (c *ArcConn) IsServer() bool {
	return c.listener != nil
}

func (c *ArcConn) loopTCPRead() (bool, error) {
	line, _, err := bufio.NewReader(c.socketTCP).ReadLine()
	if err != nil {
		if err.Error() == "EOF" {
			// Client Error
			c.session.RequestRestore()
			return true, err
		} else {
			// Server Error
			return false, err
		}
	}
	c.inputTCP(line)
	return true, nil
}


func (c *ArcConn) loopUDPRead() (bool, error) {
	buf := make([]byte, 1480)
	size, _, _ := c.socketUDP.ReadFromUDP(buf)
	c.inputUDP(buf[:size])
	return true, nil
}

func (c *ArcConn) outputTCP(b []byte) {
	c.socketTCP.Write(b)
}

func (c *ArcConn) outputUDP(b []byte) {
	c.socketUDP.WriteToUDP(b, c.remoteAddr.UDP)
}

func (c *ArcConn) outputTCPAsync(b []byte) {
	go c.outputTCP(b)
}

func (c *ArcConn) inputTCP(b []byte) {
	c.handleSegment(&b)
}

func (c *ArcConn) inputUDP(b []byte) {
	c.handleSegment(&b)
}

func (c *ArcConn) handleSegment(b *[]byte) {
	switch SegmentClass(ReadType(b)) {
	case SEGMENT_CLASS_SESSION:
		s, _ := ParseTCPSegmentSessionCommand(b)
		switch s.Type {
		case TCP_SEGMENT_SESSION_NEW:
			c.handleSegmentNewSession(s)
		case TCP_SEGMENT_SESSION_RESPONSE:
			c.handleSegmentResponseSession(s)
		default:
			c.session.HandleSessionSegment(b, c.remoteAddr)
		}
	}
}

func (c *ArcConn) handleSegmentNewSession(s *TCPSegmentSessionCommand) {
	if s.SessionId == 0 {
		_, session := c.listener.sessionSet.New(c)
		session.Respond()
		c.session = session
	} else {
		// Restore connection
		session := c.listener.sessionSet.Get(s.SessionId)
		session.Restore(c)
	}
}

func (c *ArcConn) handleSegmentResponseSession(s *TCPSegmentSessionCommand) {
	c.session.Id = s.SessionId
	c.session.ProbeAsync()
}

func (c *ArcConn) DumpRemoteAddr() string {
	return c.remoteAddr.String()
}