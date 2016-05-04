package protocol

import (
	"encoding/binary"
	"fmt"
)

type HeartPacket struct {
	imei       uint32
	sum        uint8
	protocolid uint8
}

func (p *HeartPacket) Serialize() []byte {
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

func NewHeartPacket(imei uint32, protocolid uint8, sum uint8) *HeartPacket {
	return &HeartPacket{
		imei:       imei,
		protocolid: protocolid,
		sum:        sum,
	}
}

func ParseHeart(buffer []byte) (*HeartPacket, *DasHeartPacket) {
	protocolid := buffer[2]
	imei := binary.BigEndian.Uint32(buffer[5:9])
	sum := buffer[10]

	return NewHeartPacket(imei, protocolid, sum), NewDasHeartPacket(fmt.Sprint(imei), "0")

}
