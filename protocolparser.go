package shunt_collars

import (
	"bytes"
	"encoding/binary"
	//	"errors"
)

var (
	Flag          byte   = 0x29
	DasLeftFlag   byte   = '$'
	DasRightFlag1 byte   = '\r'
	DasRightFlag2 byte   = '\n'
	Comma         []byte = []byte{','}
	Asterisk      []byte = []byte{'*'}
	Colon         []byte = []byte{':'}

	protocolID = map[uint8]uint16{
		0xB1: HeartBeat,
		0xB2: PosUp,
	}

	protocolIDDas = map[string]uint16{
		"$LOGRT":  Login,
		"$HCHECK": HeartBeat,
	}
)

//main protocol format 0x29 0x29 protocolid(byte) len(word) terminalid(dwords) content(nbytes) checksum(byte) endflag
func CheckProtocol(buffer *bytes.Buffer) (uint16, uint16) {
	bufferlen := buffer.Len()
	if bufferlen == 0 {
		return Illegal, 0
	}
	if bufferlen > 4 {
		if buffer.Bytes()[0] != Flag && buffer.Bytes()[1] != Flag {
			buffer.ReadByte()
			CheckProtocol(buffer)
		}
		protocolid := protocolID[buffer.Bytes()[2]]
		protocollen := binary.BigEndian.Uint16(buffer.Bytes()[3:5])
		if int(protocollen+5) < bufferlen {
			return HalfPack, 0
		}

		return protocolid, protocollen + 5

	}

	return HalfPack, 0
}

func parseDasCmdID(buf *bytes.Buffer) uint16 {
	split := bytes.Split(buf.Bytes(), Colon)

	return protocolIDDas[string(split[0])]

}

func CheckProtocolDas(buffer *bytes.Buffer) (uint16, uint16) {
	bufferlen := buffer.Len()
	if bufferlen == 0 {
		return Illegal, 0
	}
	if buffer.Bytes()[0] != DasLeftFlag {
		buffer.ReadByte()
		CheckProtocol(buffer)
	}
	for i := 1; i < bufferlen; i++ {
		if buffer.Bytes()[i] == DasRightFlag2 && buffer.Bytes()[i-1] == DasRightFlag1 {
			return parseDasCmdID(buffer), uint16(i)
		}
	}

	return HalfPack, 0

}
