package go2neo

import (
	"errors"
	"net"
)

var preamble []byte = []byte{0x60, 0x60, 0xB0, 0x17}
var handshakeRequest []byte = []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

type driver struct {
	BoltVersion uint
	readBuffer  []byte
}

type ProtocolError struct{}

func (err ProtocolError) Error() string {
	return "Bad protocol things have occurred. This is not cool, man."
}

func Driver(address string) (*driver, error) {
	d := new(driver)
	d.readBuffer = make([]byte, 65536)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, errors.New("Cannot listen to server")
	}

	// perform handshake
	conn.Write(preamble)
	conn.Write(handshakeRequest)
	size, err := conn.Read(d.readBuffer)
	if err != nil {
		return nil, errors.New("Cannot receive data")
	}
	if size != 4 {
		return nil, errors.New("Not enough data")
	}

	d.BoltVersion = (uint)(d.readBuffer[0]<<24 |
		d.readBuffer[1]<<16 |
		d.readBuffer[2]<<8 |
		d.readBuffer[3])
	return d, nil
}
