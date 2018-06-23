package sl500_api

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

var Baud = baudRegistry()

func baudRegistry() *baudList {
	return &baudList{
		Baud4800:   baud{0, 4800},
		Baud9600:   baud{1, 9600},
		Baud14400:  baud{2, 14400},
		Baud19200:  baud{3, 19200},
		Baud28800:  baud{4, 28800},
		Baud38400:  baud{5, 38400},
		Baud57600:  baud{6, 57600},
		Baud115200: baud{7, 115200},
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
	Baud4800   baud
	Baud9600   baud
	Baud14400  baud
	Baud19200  baud
	Baud28800  baud
	Baud38400  baud
	Baud57600  baud
	Baud115200 baud
}

type baud struct {
	ByteValue byte
	IntValue  int
}

type Sl500 struct {
	config *serial.Config
	port   *serial.Port
}

func NewConnection(path string, baud baud) (Sl500, error) {
	c := &serial.Config{Name: path, Baud: baud.IntValue, ReadTimeout: 5 * time.Second} // TODO
	o, err := serial.OpenPort(c)

	res := Sl500{}

	if err != nil {
		return res, err
	}

	res.config = c
	res.port = o

	_, err = res.RfInitCom(baud.ByteValue)
	if err != nil {
		return res, err
	}

	_, err = res.RfInitType(Type_ISO)
	if err != nil {
		return res, err
	}

	_, err = res.RfAntennaSta(AntennaOff)
	if err != nil {
		return res, err
	}

	_, err = res.RfGetDeviceNumber()
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *Sl500) RfInitCom(baud byte) ([]byte, error) {
	sendRequest(s.port, 0x0101, []byte{baud})
	return readResponse(s.port)
}

func (s *Sl500) RfGetDeviceNumber() ([]byte, error) {
	sendRequest(s.port, 0x0301, []byte{})
	return readResponse(s.port)
}

func (s *Sl500) RfAntennaSta(antennaState byte) ([]byte, error) {
	sendRequest(s.port, 0x0C01, []byte{antennaState})
	return readResponse(s.port)
}

func (s *Sl500) RfInitType(workType byte) ([]byte, error) {
	sendRequest(s.port, 0x0801, []byte{workType})
	return readResponse(s.port)
}

func (s *Sl500) RfBeep(durationMs byte) ([]byte, error) {
	sendRequest(s.port, 0x0601, []byte{durationMs})
	return readResponse(s.port)
}

func readResponse(port *serial.Port) ([]byte, error) {
	var buf []byte
	innerBuf := make([]byte, 128)

	totalRead := 0
	readTriesCount := 0
	maxReadCount := 50

	for ; ; {
		readTriesCount += 1

		if readTriesCount >= maxReadCount {
			return nil, fmt.Errorf("Reads tries exceeded")
		}

		n, err := port.Read(innerBuf)

		if err != nil {
			return nil, err
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
		return nil, fmt.Errorf("Response format invalid")
	}

	buf = buf[3:]
	ver := byte(0x00)

	for _, k := range buf[:len(buf)-1] {
		ver = ver ^ k
	}

	if ver != buf[len(buf)-1] {
		return nil, fmt.Errorf("Response verification failed")
	}

	return buf[3 : len(buf)-1], nil
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
