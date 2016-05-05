package protocol

import (
	"fmt"
	"log"
	"strconv"
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
	termID0 := buffer[5]
	termID1 := buffer[6] - 0x80
	termID2 := buffer[7] - 0x80
	termID3 := buffer[8]

	strTermID0 := fmt.Sprintf("%02d", termID0)
	strTermID1 := fmt.Sprintf("%02d", termID1)
	strTermID2 := fmt.Sprintf("%02d", termID2)
	strTermID3 := fmt.Sprintf("%02d", termID3)
	log.Println(strTermID0)
	log.Println(strTermID1)
	log.Println(strTermID2)
	log.Println(strTermID3)

	strTermID := strTermID0 + strTermID1 + strTermID2 + strTermID3

	tmp, _ := strconv.Atoi(strTermID)
	imei := uint32(tmp)

	sum := buffer[10]

	return NewLoginPakcet(imei, protocolid, sum), NewDasLoginPacket(fmt.Sprint(imei)), imei
}
