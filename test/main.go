package main

import (
	"fmt"

	"github.com/nigelsmall/go2neo"
)

func main() {
	driver, err := go2neo.Driver("127.0.0.1:7687")
	if err != nil {
	  panic(err)
	}
	fmt.Println("Handshake complete, using Bolt version", driver.BoltVersion)
}
