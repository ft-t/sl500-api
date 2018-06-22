package sl500_api

import (
	"github.com/tarm/serial"
	"log"
	"encoding/binary"
	"bytes"
	)

var Baud = baudRegistry()

func baudRegistry() *baudList {
	return &baudList{
		Baud19200: baud{3, 19200},
	}
}

const (
	Type_A     = byte(0xA)
	Type_B     = byte(0xB)
	Type_ISO   = byte(0x1)
	AntennaOn  = byte(0x1)
	AntennaOff = byte(0x0)
)

type baudList struct {
	Baud19200 baud
}

type baud struct {
	ByteValue byte
	IntValue  int
}

type Sl500 struct {
	config *serial.Config
	port   *serial.Port
}

func NewConnection(path string, baud baud) Sl500 {
	c := &serial.Config{Name: path, Baud: baud.IntValue} // TODO
	o, err := serial.OpenPort(c)

	if err != nil {
		log.Fatal(err)
	}

	res := Sl500{config: c, port: o}
	log.Println("RfInitCom", res.RfInitCom(baud.ByteValue))
	res.RfInitType(Type_ISO)
	log.Println("AntennaSta", res.RfAntennaSta(AntennaOff))
	log.Println("Device ID", res.RfGetDeviceNumber())

	return res
}

func (s *Sl500) RfInitCom(baud byte) []byte {
	sendRequest(s.port, 0x0101, []byte{baud})
	return readResponse(s.port)
}

func (s *Sl500) RfGetDeviceNumber() []byte {
	sendRequest(s.port, 0x0301, []byte{})
	return readResponse(s.port)
}
func (s *Sl500) RfAntennaSta(antennaState byte) []byte {
	sendRequest(s.port, 0x0C01, []byte{antennaState})
	return readResponse(s.port)
}
func (s *Sl500) RfInitType(workType byte) {
	sendRequest(s.port, 0x0801, []byte{workType})
	log.Println("RfInitType", readResponse(s.port))
}
func (s *Sl500) RfBeep(durationMs byte) {
	sendRequest(s.port, 0x0601, []byte{durationMs})
	log.Println("RfBeep", readResponse(s.port))
}
func readResponse(port *serial.Port) []byte {
	var buf []byte
	totalRead := 0
	readTriesCount := 0
	maxReadCount := 50
	for ; ; {
		readTriesCount += 1

		if readTriesCount >= maxReadCount {
			break
		}

		innerBuf := make([]byte, 128)
		n, err := port.Read(innerBuf)

		if err != nil {
			log.Println(err)
		}

		totalRead += n
		buf = append(buf, innerBuf[:n]...)

		if totalRead < 3 {
			continue
		}
		if int(buf[2]) != len(buf)-4 {
			continue
		}
		break

	}
	if buf[0] != 0xAA || buf[1] != 0xBB {
		log.Println("shit happens")
	}

	return buf
}
func sendRequest(port *serial.Port, commandCode int16, bytesData []byte) {
	buf := new(bytes.Buffer)

	ver := byte(0x00)
	length := len(bytesData) + 5

	binary.Write(buf, binary.BigEndian, byte(0xAA))
	binary.Write(buf, binary.BigEndian, byte(0xBB))
	binary.Write(buf, binary.BigEndian, byte(length))
	binary.Write(buf, binary.BigEndian, byte(0x00))
	binary.Write(buf, binary.BigEndian, int16(0)) // device id
	binary.Write(buf, binary.BigEndian, commandCode)
	binary.Write(buf, binary.BigEndian, bytesData)

	for _, k := range buf.Bytes()[3:] {
		ver = ver ^ k
	}
	binary.Write(buf, binary.BigEndian, ver)

	log.Println(buf.Bytes())

	port.Write(buf.Bytes())
}
