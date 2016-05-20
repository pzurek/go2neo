package bolt

import (
  "fmt"
  "net"
)

var preamble []byte = []byte{0x60, 0x60, 0xB0, 0x17}
var handshakeRequest []byte = []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00,
                                     0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

type client struct {
    readBuffer []byte
}

func ConnectTCP() {
  c := client{readBuffer: make([]byte, 65536)}

  conn, error := net.Dial("tcp", "127.0.0.1:7687")
  if error != nil {
    fmt.Println("Cannot listen to server: ", error)
    return
  }

  // perform handshake
  conn.Write(preamble)
  conn.Write(handshakeRequest)
  size, error := conn.Read(c.readBuffer)
  if error != nil {
    fmt.Println("Cannot receive data: ", error)
    return
  }
  if size != 4 {
    fmt.Println("Not enough data: ", size)
    return
  }

  fmt.Println("Handshake response: ", c.readBuffer[:4])

}
