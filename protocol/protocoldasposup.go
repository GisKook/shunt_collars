package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"strconv"
)

var (
	BCD2 = map[uint8]string{
		0x0: "0",
		0x1: "1",
		0x2: "2",
		0x3: "3",
		0x4: "4",
		0x5: "5",
		0x6: "6",
		0x7: "7",
		0x8: "8",
		0x9: "9",
	}
	BCD = map[uint8]string{
		0x0:  "00",
		0x1:  "01",
		0x2:  "02",
		0x3:  "03",
		0x4:  "04",
		0x5:  "05",
		0x6:  "06",
		0x7:  "07",
		0x8:  "08",
		0x9:  "09",
		0x10: "10",
		0x11: "11",
		0x12: "12",
		0x13: "13",
		0x14: "14",
		0x15: "15",
		0x16: "16",
		0x17: "17",
		0x18: "18",
		0x19: "19",
		0x20: "20",
		0x21: "21",
		0x22: "22",
		0x23: "23",
		0x24: "24",
		0x25: "25",
		0x26: "26",
		0x27: "27",
		0x28: "28",
		0x29: "29",
		0x30: "30",
		0x31: "31",
		0x32: "32",
		0x33: "33",
		0x34: "34",
		0x35: "35",
		0x36: "36",
		0x37: "37",
		0x38: "38",
		0x39: "39",
		0x40: "40",
		0x41: "41",
		0x42: "42",
		0x43: "43",
		0x44: "44",
		0x45: "45",
		0x46: "46",
		0x47: "47",
		0x48: "48",
		0x49: "49",
		0x50: "50",
		0x51: "51",
		0x52: "52",
		0x53: "53",
		0x54: "54",
		0x55: "55",
		0x56: "56",
		0x57: "57",
		0x58: "58",
		0x59: "59",
		0x60: "60",
		0x61: "61",
		0x62: "62",
		0x63: "63",
		0x64: "64",
		0x65: "65",
		0x66: "66",
		0x67: "67",
		0x68: "68",
		0x69: "69",
		0x70: "70",
		0x71: "71",
		0x72: "72",
		0x73: "73",
		0x74: "74",
		0x75: "75",
		0x76: "76",
		0x77: "77",
		0x78: "78",
		0x79: "79",
		0x80: "80",
		0x81: "81",
		0x82: "82",
		0x83: "83",
		0x84: "84",
		0x85: "85",
		0x86: "86",
		0x87: "87",
		0x88: "88",
		0x89: "89",
		0x90: "90",
		0x91: "91",
		0x92: "92",
		0x93: "93",
		0x94: "94",
		0x95: "95",
		0x96: "96",
		0x97: "97",
		0x98: "98",
		0x99: "99",
	}
)

type DasPosUpPacket struct {
	IMEI      string
	Time      string
	Batt      string
	Speed     string
	Parse     string
	PosReason string
	Postype   string

	Longitude string
	Latitude  string
	Lac       string
	Cid       string
	Lac1      string
	Cid1      string
	Dbm1      string
	Lac2      string
	Cid2      string
	Dbm2      string
	Lac3      string
	Cid3      string
	Dbm3      string
}

func (p *DasPosUpPacket) Serialize() []byte {
	cmd := "$POSUP:" + p.IMEI + ":" + p.Time + ":" + p.Batt + ":" + p.Speed + ":" + p.Parse + ":" + p.PosReason + ":" + p.Postype + ":"
	if p.Postype == "0" {
		cmd += p.Longitude + "," + p.Latitude + "\r\n"
	} else {
		cmd += p.Lac + "," + p.Cid + ":" + p.Lac1 + "," + p.Cid1 + "," + p.Dbm1 + ":"
		cmd += p.Lac2 + "," + p.Cid2 + "," + p.Dbm2 + ":"
		cmd += p.Lac3 + "," + p.Cid3 + "," + p.Dbm3 + "\r\n"
	}

	return []byte(cmd)
}

type PosUpPacket struct {
	SerialID uint16
	Checksum uint8
	SubCmdID uint8
}

func (p *PosUpPacket) Serialize() []byte {
	var buf []byte
	buf = append(buf, 0x29)
	buf = append(buf, 0x29)
	buf = append(buf, 0x22)
	buf = append(buf, 0x00)
	buf = append(buf, 0x07)
	serialid := make([]byte, 2)
	binary.BigEndian.PutUint16(serialid, p.SerialID)
	buf = append(buf, serialid...)
	buf = append(buf, p.Checksum)
	buf = append(buf, p.SubCmdID)
	buf = append(buf, CheckSum(buf, uint8(len(buf))))
	buf = append(buf, 0x0D)

	return buf
}

func ParsePosUp(buffer []byte) (*DasPosUpPacket, *PosUpPacket, string) {
	batt := BCD[buffer[40]]
	checksum := buffer[len(buffer)-2]

	if buffer[32] == 'A' {
		log.Printf("A %x", buffer[0:len(buffer)])
		buf := bytes.NewReader(buffer)
		buf.Seek(5, 0)
		//	imei_bytes := make([]byte, 4)
		//	buf.Read(imei_bytes)
		//	imei := binary.BigEndian.Uint32(imei_bytes)
		termID0, _ := buf.ReadByte()
		termID1, _ := buf.ReadByte()
		termID1 -= 0x80
		termID2, _ := buf.ReadByte()
		termID2 -= 0x80
		termID3, _ := buf.ReadByte()

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

		//	buf.Seek(3, 1)
		serialid_byte := make([]byte, 2)
		buf.Read(serialid_byte)
		serialid := binary.BigEndian.Uint16(serialid_byte)
		subcmdid, _ := buf.ReadByte()

		year, _ := buf.ReadByte()
		timepos := BCD[year]
		month, _ := buf.ReadByte()
		timepos += BCD[month]
		day, _ := buf.ReadByte()
		timepos += BCD[day]
		timepos += "-"
		hour, _ := buf.ReadByte()
		timepos += BCD[hour]
		minute, _ := buf.ReadByte()
		timepos += BCD[minute]
		second, _ := buf.ReadByte()
		timepos += BCD[second]

		lat0, _ := buf.ReadByte()
		lat1, _ := buf.ReadByte()
		lat2, _ := buf.ReadByte()
		lat3, _ := buf.ReadByte()

		lat_degree_string := BCD2[lat0] + BCD2[((lat1&0xF0)>>4)]
		log.Println(lat_degree_string)
		lat_minute_string := BCD2[(lat1&0x0F)] + BCD2[((lat2&0xF0)>>4)] + "." + BCD2[(lat2&0x0F)] + BCD2[(lat3&0xF0)>>4] + BCD2[lat3&0x0F]
		lat_minute, _ := strconv.ParseFloat(lat_minute_string, 64)
		lat_minute = lat_minute / 60
		lat_minute_in_degree := strconv.FormatFloat(lat_minute, 'f', 6, 64)
		lat := lat_degree_string + lat_minute_in_degree[1:]

		long0, _ := buf.ReadByte()
		long1, _ := buf.ReadByte()
		long2, _ := buf.ReadByte()
		long3, _ := buf.ReadByte()
		long_degree_string := BCD[long0] + BCD2[(long1&0xF0)>>4]
		long_minute_string := BCD2[(long1&0x0F)] + BCD2[(long2&0xF0)>>4] + "." + BCD2[(long2&0x0F)] + BCD2[(long3&0xF0)>>4] + BCD2[long3&0x0F]
		long_minute, _ := strconv.ParseFloat(long_minute_string, 64)
		long_minute = long_minute / 60
		long_minute_in_degree := strconv.FormatFloat(long_minute, 'f', 6, 64)
		long := long_degree_string + long_minute_in_degree[1:]

		speed_bytes := make([]byte, 2)
		buf.Read(speed_bytes)
		speed := binary.BigEndian.Uint16(speed_bytes)
		parse := "0"
		posreason := "2"
		postype := "0"

		return &DasPosUpPacket{
				IMEI:      fmt.Sprint(imei),
				Time:      timepos,
				Batt:      batt,
				Speed:     fmt.Sprint(speed),
				Parse:     parse,
				PosReason: posreason,
				Postype:   postype,
				Longitude: long,
				Latitude:  lat,
			}, &PosUpPacket{
				SerialID: serialid,
				SubCmdID: subcmdid,
				Checksum: checksum,
			}, batt
	} else {
		buf := bytes.NewReader(buffer)
		buf.Seek(9, 0)
		serialid_byte := make([]byte, 2)
		buf.Read(serialid_byte)
		serialid := binary.BigEndian.Uint16(serialid_byte)
		subcmdid, _ := buf.ReadByte()

		return nil, &PosUpPacket{
			SerialID: serialid,
			SubCmdID: subcmdid,
			Checksum: checksum,
		}, batt
	}

	return nil, nil, batt

}
