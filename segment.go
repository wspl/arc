package arc

const ARC_PROTOCOL_VERSION uint16 = 0x01

const(
	SEGMENT_CLASS_SESSION        uint16 = 0x1

	TCP_SEGMENT_SESSION_NEW      uint16 = 0x11
	TCP_SEGMENT_SESSION_RESPONSE uint16 = 0x12
	TCP_SEGMENT_SESSION_ACCEPT   uint16 = 0x13

	UDP_SEGMENT_SESSION_PROBE uint16 = 0x1a
	UDP_SEGMENT_DATA          uint16 = 0x2a
)

func SegmentClass(t uint16) uint16 {
	return t >> 4
}

func SegmentSubId(t uint16) uint16 {
	return t & 0x0f
}

func SegmentProtocolTCP(t uint16) bool {
	tt := SegmentSubId(t)
	return 0x1 <= tt && tt <= 0x9
}

func SegmentProtocolUDP(t uint16) bool {
	tt := SegmentSubId(t)
	return 0xa <= tt && tt <= 0xf
}

func ReadVer(b *[]byte) uint16 { return ReadUInt16(b, 0) }
func ReadType(b *[]byte) uint16 { return ReadUInt16(b, 2) }

func WriteVer(b *[]byte, v uint16) { WriteUInt16(b, 0, v) }
func WriteType(b *[]byte, v uint16) { WriteUInt16(b, 2, v) }

type TCPSegmentSessionCommand struct {
	Type uint16
	SessionId uint32
}
func (s *TCPSegmentSessionCommand) Binary() []byte {
	b := make([]byte, 9)
	WriteVer(&b, ARC_PROTOCOL_VERSION)
	WriteType(&b, s.Type)
	WriteUInt32(&b, 4, s.SessionId)
	WriteLF(&b, 8)
	return b
}
func ParseTCPSegmentSessionCommand(b *[]byte) (*TCPSegmentSessionCommand, error) {
	return &TCPSegmentSessionCommand{
		Type: ReadUInt16(b, 2),
		SessionId: ReadUInt32(b, 4),
	}, nil
}

type UDPSegmentSessionCommand struct {
	Type uint16
	SessionId uint32
}
func (s *UDPSegmentSessionCommand) Binary() []byte {
	b := make([]byte, 8)
	WriteVer(&b, ARC_PROTOCOL_VERSION)
	WriteType(&b, s.Type)
	WriteUInt32(&b, 4, s.SessionId)
	return b
}
func ParseUDPSegmentSessionCommand(b *[]byte) (*UDPSegmentSessionCommand, error) {
	return &UDPSegmentSessionCommand{
		Type: ReadUInt16(b, 2),
		SessionId: ReadUInt32(b, 4),
	}, nil
}
