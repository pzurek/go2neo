package graphdatabase

import (
  "fmt"
  "net"
)

var preamble []byte = []byte{0x60, 0x60, 0xB0, 0x17}
var handshakeRequest []byte = []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
                                     0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

type driver struct {
  BoltVersion int
  readBuffer []byte
}

func Driver(address string) *driver {
  d := new(driver)
  d.readBuffer = make([]byte, 65536)

  conn, error := net.Dial("tcp", address)
  if error != nil {
    fmt.Println("Cannot listen to server: ", error)
    return nil
  }

  // perform handshake
  conn.Write(preamble)
  conn.Write(handshakeRequest)
  size, error := conn.Read(d.readBuffer)
  if error != nil {
    fmt.Println("Cannot receive data: ", error)
    return nil
  }
  if size != 4 {
    fmt.Println("Not enough data: ", size)
    return nil
  }

  d.BoltVersion = (int) (d.readBuffer[0] << 24 |
                         d.readBuffer[1] << 16 |
                         d.readBuffer[2] << 8 |
                         d.readBuffer[3])
  return d
}
