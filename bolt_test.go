package go2neo_test

import (
	"testing"

	. "github.com/nigelsmall/go2neo"
  . "gopkg.in/check.v1"
)

type driverSuite struct{}

var _ = Suite(&driverSuite{})

// Hook up gocheck into the "go test" runner. (only once per package)
func Test(t *testing.T) {
	TestingT(t)
}

func (s *driverSuite) TestBoltHandshake(c *C) {
	driver, err := Driver("127.0.0.1:7687")
	c.Assert(err, IsNil)
	c.Assert(driver.BoltVersion, Equals, uint(1))
}
