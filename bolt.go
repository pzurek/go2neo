package go2neo

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net"
)

const (
	INIT = 0x01
)

var (
	preamble         = []byte{0x60, 0x60, 0xB0, 0x17}
	handshakeRequest = []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

type Driver struct {
	BoltVersion uint
	connection  net.Conn
	readBuffer  []byte
	writeBuffer *bytes.Buffer
	packer      *Encoder
}

type ProtocolError struct{}

func (err ProtocolError) Error() string {
	return "Bad protocol things have occurred. This is not cool, man."
}

func (d *Driver) sendRaw() error {
	res, err := ioutil.ReadAll(d.writeBuffer)
	if err != nil {
		return err
	}
	d.connection.Write(res)
	return nil
}

func (d *Driver) send() error {
	res, err := ioutil.ReadAll(d.writeBuffer)
	if err != nil {
		return err
	}
	chunkSize := uint16(len(res))
	d.connection.Write([]byte{byte(chunkSize >> 8), byte(chunkSize)})
	d.connection.Write(res)
	d.connection.Write([]byte{0x00, 0x00})
	return nil
}

func (d *Driver) recvRaw(size uint16) error {
	z, err := d.connection.Read(d.readBuffer[0:size])
	if z != int(size) {
		return errors.New("Not enough data")
	}
	return err
}

func (d *Driver) recv() (int, error) {
	err := d.recvRaw(2)
	if err != nil { return -1, err }
	chunkSize := (int)(d.readBuffer[0]<<8 | d.readBuffer[1])
	if chunkSize > 0 {
		err := d.recvRaw(uint16(chunkSize))
		if err != nil { return -1, err }
	}
	return chunkSize, nil
}

func (d *Driver) handshake() error {
	d.writeBuffer.Reset()
	d.writeBuffer.Write(preamble)
	d.writeBuffer.Write(handshakeRequest)
	d.sendRaw()
	d.recvRaw(4)
	d.BoltVersion = (uint)(d.readBuffer[0]<<24 |
		d.readBuffer[1]<<16 |
		d.readBuffer[2]<<8 |
		d.readBuffer[3])
	return nil
}

func (d *Driver) init(userAgent string, user string, password string) error {
	d.writeBuffer.Reset()
	d.writeBuffer.Write([]byte{0xB2, INIT})
	err := d.packer.Encode(userAgent)
	if err != nil { return err }
	d.writeBuffer.Write([]byte{0xA3})
	err = d.packer.Encode("scheme", "basic", "principal", user, "credentials", password)
	if err != nil { return err }
	err = d.send()
	if err != nil { return err }
	more := true
	for more {
		size, err := d.recv()
  	if err != nil { return err }
		if (size == 0) {
			more = false
		} else {
			// TODO: use the return value
  		print(d.readBuffer[0:size])
		}
	}
	return nil
}

func NewDriver(address string) (*Driver, error) {
	d := new(Driver)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, errors.New("Cannot connect to server")
	}
	d.connection = conn
	d.readBuffer = make([]byte, 65536)
	d.writeBuffer = &bytes.Buffer{}
	d.packer, err = NewEncoder(d.writeBuffer)
	if err != nil {
		return nil, err
	}

  err = d.handshake()
	if err != nil {
		return nil, err
	}

	err = d.init("go2neo/0.0.0", "neo4j", "password")
	if err != nil {
		return nil, err
	}

	return d, nil
}
