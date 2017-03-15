package arc

import "encoding/binary"

func WriteLF(b *[]byte, offset int) {
	lf := '\n'
	(*b)[offset] = byte(lf)
}

func WriteUInt16(b *[]byte, offset int, v uint16) {
	ib := make([]byte, 2)
	binary.BigEndian.PutUint16(ib, v)
	for i := 0; i < 2; i += 1 {
		(*b)[offset+i] = ib[i]
	}
}

func WriteUInt32(b *[]byte, offset int, v uint32) {
	ib := make([]byte, 4)
	binary.BigEndian.PutUint32(ib, v)
	for i := 0; i < 4; i += 1 {
		(*b)[offset+i] = ib[i]
	}
}

func WriteUInt64(b *[]byte, offset int, v uint64) {
	ib := make([]byte, 8)
	binary.BigEndian.PutUint64(ib, v)
	for i := 0; i < 8; i += 1 {
		(*b)[offset+i] = ib[i]
	}
}

func CopyBytes(src *[]byte, sOffset int, dst *[]byte, dOffset int, length int) {
	for i := 0; i < length; i++ {
		(*dst)[dOffset + i] = (*src)[sOffset + i]
	}
}

func ReadUInt16(b *[]byte, offset int) uint16 {
	return binary.BigEndian.Uint16((*b)[offset : offset+2])
}

func ReadUInt32(b *[]byte, offset int) uint32 {
	return binary.BigEndian.Uint32((*b)[offset : offset+4])
}

func ReadUInt64(b *[]byte, offset int) uint64 {
	return binary.BigEndian.Uint64((*b)[offset : offset+8])
}