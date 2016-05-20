package bolt

import (
  "fmt"
  "net"
)

func ConnectTCP() {
  _, tcperr := net.Dial("tcp", "127.0.0.1:7687") // replace _ with variable, e.g. "conn" once used

  if tcperr != nil {
    fmt.Println("Cannot listen to server: ", tcperr)
  }
}

// introduce handshake here
