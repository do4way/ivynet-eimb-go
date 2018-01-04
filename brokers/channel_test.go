package brokers

import (
	"errors"
	"time"

	"github.com/do4way/ivynet-eimb-go/redisConfig"

	. "gopkg.in/check.v1"
	redis "gopkg.in/redis.v5"
)

type ChannelTestSuite struct {
	broker      *Broker
	redisClient *redis.Client
}

var (
	_ = Suite(&ChannelTestSuite{})
)

func (s *ChannelTestSuite) SetUpSuite(c *C) {
	config := redisConfig.Read()
	broker := New("TestChannel")
	s.broker = broker
	s.redisClient = redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.DB,
	})

}

func (s *ChannelTestSuite) TestWrite(c *C) {
	oldLen := s.redisClient.LLen(inputChannelKey("TestChannel"))
	s.write("This is a test message")
	len := s.redisClient.LLen(inputChannelKey("TestChannel")).Val() - oldLen.Val()
	s.redisClient.LPop(inputChannelKey("TestChannel"))
	c.Assert(len, Equals, int64(1))
}

func (s *ChannelTestSuite) TestWriteRead(c *C) {

	pck := s.write("Test WriteRead")

	input := s.broker.GetInputChannel("Consumer1")
	rst, err := input.Read(1)
	c.Assert(err, Equals, nil)
	c.Assert(rst.getString(), Equals, pck.getString())

	c.Assert(s.redisClient.LLen(workingChannelKey(s.broker.Name, "Consumer1")).Val(), Equals, int64(1))
	s.redisClient.LRem(workingChannelKey(s.broker.Name, "Consumer1"), 0, rst.getString())
}

func (s *ChannelTestSuite) TestWriteReadAndAck(c *C) {

	s.write("Test WriteReadAndAck")

	cnt := s.redisClient.LLen(workingChannelKey(s.broker.Name, "Consumer2")).Val()
	input := s.broker.GetInputChannel("Consumer2")
	rst, err := input.Read(10)
	c.Assert(err, Equals, nil)

	c.Assert(s.redisClient.LLen(workingChannelKey(s.broker.Name, "Consumer2")).Val()-cnt, Equals, int64(1))
	err = rst.Ack()
	c.Assert(err, Equals, nil)
	c.Assert(s.redisClient.LLen(workingChannelKey(s.broker.Name, "Consumer2")).Val()-cnt, Equals, int64(0))
}

func (s *ChannelTestSuite) TestWriteReadAndFail(c *C) {

	cName := "TestWriteReadAndFail"
	s.write("Test WriteReadAndFail")

	cnt := s.redisClient.LLen(failChannelKey(s.broker.Name, cName)).Val()
	input := s.broker.GetInputChannel(cName)
	rst, err := input.Read(10)
	c.Assert(err, Equals, nil)

	c.Assert(s.redisClient.LLen(failChannelKey(s.broker.Name, cName)).Val()-cnt, Equals, int64(0))
	err = rst.Fail()
	c.Assert(err, Equals, nil)
	c.Assert(s.redisClient.LLen(failChannelKey(s.broker.Name, cName)).Val()-cnt, Equals, int64(1))
}

func (s *ChannelTestSuite) TestWriteAsyncRead(c *C) {

	input := s.broker.GetInputChannel("Consumer3")
	rst := make(chan bool)
	input.AsyncRead().Then(func(d interface{}) (interface{}, error) {
		if d == nil {
			rst <- false
			return nil, errors.New("Unexpected null data")
		}
		p, ok := d.(*Package)
		if !ok {
			rst <- false
			return nil, errors.New("Received an unexpected package")
		}
		p.Ack()
		if !c.Check(p.Payload, Equals, "Test WriteAsyncRead") {
			rst <- false
			return nil, errors.New("Not expected package")
		}
		rst <- true
		return d, nil
	}, func(err error) error {
		c.Fail()
		rst <- false
		return err
	})
	s.write("Test WriteAsyncRead")
	c.Assert(<-rst, Equals, true)
}

func (s *ChannelTestSuite) write(msg string) *Package {
	output := s.broker.GetOutputChannel()
	pck := &Package{
		Payload:   "Test WriteAsyncRead",
		CreatedAt: time.Now(),
	}
	output.Write(pck)
	return pck
}
