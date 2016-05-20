package go2neo

import (
	"testing"
)

func TestBoltHandshake(t *testing.T) {
	driver, err := Driver("127.0.0.1:7687")
	if err != nil {
		t.Error(err)
	}
	if driver.BoltVersion != 1 {
		t.Error("incorrect bolt version")
	}
}
