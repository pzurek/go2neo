package main

import (
  "fmt"

  "github.com/nigelsmall/go2neo/graphdatabase"
)

func main() {
  driver, _ := graphdatabase.Driver("127.0.0.1:7687")
  fmt.Println("Handshake complete, using Bolt version", driver.BoltVersion)
}
