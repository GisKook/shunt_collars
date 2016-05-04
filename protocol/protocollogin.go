package protocol

import (
	"encoding/binary"
	"fmt"
)

func CheckSum(buf []byte, length uint8) byte {
	temp := buf[0]
	for i := uint8(1); i < length; i++ {
		temp ^= buf[i]
	}

	return temp
}

type LoginPacket struct {
	imei       uint32
	sum        uint8
	protocolid uint8
}

func (p *LoginPacket) Serialize() []byte {
	var buf []byte
	buf = append(buf, 0x29)
	buf = append(buf, 0x29)
	buf = append(buf, 0x21)
	buf = append(buf, 0x00)
	buf = append(buf, 0x05)
	buf = append(buf, p.sum)
	buf = append(buf, p.protocolid)
	buf = append(buf, 0x00)
	buf = append(buf, CheckSum(buf, uint8(len(buf))))
	buf = append(buf, 0x0D)

	return buf
}

func NewLoginPakcet(imei uint32, protocolid uint8, sum uint8) *LoginPacket {
	return &LoginPacket{
		imei:       imei,
		protocolid: protocolid,
		sum:        sum,
	}
}

func ParseLogin(buffer []byte) (*LoginPacket, *DasLoginPacket, uint32) {
	protocolid := buffer[2]
	imei := binary.BigEndian.Uint32(buffer[5:9])
	sum := buffer[10]

	return NewLoginPakcet(imei, protocolid, sum), NewDasLoginPacket(fmt.Sprint(imei)), imei
}
