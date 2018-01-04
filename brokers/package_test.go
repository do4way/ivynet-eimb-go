package brokers

import (
	"encoding/json"
	"time"

	. "gopkg.in/check.v1"

	"github.com/do4way/ivynet-eimb-go/redisConfig"

	redis "gopkg.in/redis.v5"
)

// func TestPackage(t *testing.T) {
// 	TestingT(t)
// }

type PackageTestSuite struct {
	redis *redis.Client
}

var (
	_ = Suite(&PackageTestSuite{})
)

func (s *PackageTestSuite) SetUpSuite(c *C) {
	config := redisConfig.Read()
	s.redis = redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.DB,
	})
}

func (s *PackageTestSuite) TestGetString(c *C) {
	p := &Package{
		Payload:           "Hello world",
		CreatedAt:         time.Now(),
		workingChannelKey: "TestWorkingChannelKey",
		channelKey:        "TestChannelKey",
	}

	rst := &Package{}
	err := json.Unmarshal([]byte(p.getString()), rst)
	c.Check(err, IsNil)

	c.Check(rst.Payload, Equals, "Hello world")
	c.Check(rst.workingChannelKey, Equals, "")
}

func (s *PackageTestSuite) TestAck(c *C) {
	p := &Package{
		Payload:           "Hello world",
		CreatedAt:         time.Now(),
		workingChannelKey: "TestWorkingChannelKey",
		redis:             s.redis,
	}
	s.redis.LPush(p.workingChannelKey, p.getString())
	p.Ack()
	c.Check(s.redis.LLen(p.workingChannelKey).Val(), Equals, int64(0))
}

func (s *PackageTestSuite) TestFail(c *C) {

	p := &Package{
		Payload:           "Hello world",
		CreatedAt:         time.Now(),
		workingChannelKey: "TestWorkingChannelKey",
		failedChannelKey:  "TestWorkingChannelKey_Failed",
		redis:             s.redis,
	}
	cnt := s.redis.LLen(p.failedChannelKey).Val()
	s.redis.LPush(p.workingChannelKey, p.getString())
	c.Check(p.Fail(), IsNil)
	c.Check(s.redis.LLen(p.failedChannelKey).Val()-cnt, Equals, int64(1))
}
