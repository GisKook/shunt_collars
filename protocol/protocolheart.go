package protocol

import (
	"fmt"
	"strconv"
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
	termID0 := buffer[5]
	termID1 := buffer[6] - 0x80
	termID2 := buffer[7] - 0x80
	termID3 := buffer[8]

	strTermID0 := fmt.Sprintf("%02d", termID0)
	strTermID1 := fmt.Sprintf("%02d", termID1)
	strTermID2 := fmt.Sprintf("%02d", termID2)
	strTermID3 := fmt.Sprintf("%02d", termID3)

	strTermID := strTermID0 + strTermID1 + strTermID2 + strTermID3
	tmp, _ := strconv.Atoi(strTermID)
	imei := uint32(tmp)

	sum := buffer[10]

	return NewHeartPacket(imei, protocolid, sum), NewDasHeartPacket(fmt.Sprint(imei), "0")

}
