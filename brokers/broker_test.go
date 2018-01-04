package brokers

import (
	redis "github.com/do4way/ivynet-eimb-go/redisConfig"
	. "gopkg.in/check.v1"
	redisDB "gopkg.in/redis.v5"
)

type BrokerTestSuite struct {
	broker      *Broker
	redisClient *redisDB.Client
}

var (
	_ = Suite(&BrokerTestSuite{})
)

func (s *BrokerTestSuite) SetUpSuite(c *C) {
	config := redis.Read()
	broker := New("TestBroker")
	s.broker = broker
	s.redisClient = redisDB.NewClient(&redisDB.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.DB,
	})
}

func (s *BrokerTestSuite) TestPing(c *C) {
	rst, err := s.broker.Ping()
	c.Check(err, Equals, nil)
	c.Assert(rst, Equals, "PONG")
	c.Assert(err, Equals, nil)
}
