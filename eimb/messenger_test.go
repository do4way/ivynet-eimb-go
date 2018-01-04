package eimb

import (
	"errors"
	"fmt"

	"github.com/do4way/ivynet-eimb-go/brokers"
	"github.com/do4way/ivynet-eimb-go/redisConfig"
	. "gopkg.in/check.v1"
	redis "gopkg.in/redis.v5"
)

type MessengerTestSuite struct {
	redis *redis.Client
}

var (
	_ = Suite(&MessengerTestSuite{})
)

func (s *MessengerTestSuite) SetUpSuite(c *C) {

	config := redisConfig.Read()
	s.redis = redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.DB,
	})
}

func (s *MessengerTestSuite) TestSend(c *C) {

	channelKey := "TestTopic_TestSend"
	cnt := s.redis.LLen(inputChannelKey(channelKey)).Val()
	messenger := NewMessenger(channelKey)
	err := messenger.Send("message")
	c.Check(err, Equals, nil)
	nowCnt := s.redis.LLen(inputChannelKey(channelKey)).Val()
	c.Assert(nowCnt-cnt, Equals, int64(1))
	// s.redis.LPop(inputChannelKey("TestTopic_TestEmit"))
}

func (s *MessengerTestSuite) TestRequest(c *C) {

	messenger := NewMessenger("TestRequest_Test")

	rst := make(chan bool)
	messenger.Request("Message", 10).
		Then(func(data interface{}) (interface{}, error) {
			c.Check(data, Not(Equals), nil)
			p, ok := data.(*brokers.Package)
			if !ok {
				rst <- false
				return nil, errors.New("Received an unexpected data, expected type package, but get ")
			}
			c.Check(p.Payload, Equals, "OK!")
			rst <- true
			return p, nil
		}, func(err error) error {
			rst <- false
			c.Fail()
			return err
		})
	messenger.broker.GetInputChannel("RequestTesterInputReader").AsyncRead().
		Then(func(data interface{}) (interface{}, error) {
			p, ok := data.(*brokers.Package)
			if !ok {
				return nil, errors.New("Received unexpected data")
			}
			err := p.Response("OK!")
			if err != nil {
				return nil, err
			}
			p.Ack()
			return data, nil
		}, func(err error) error {
			c.Fail()
			return err
		})
	c.Assert(<-rst, Equals, true)
}

type Handler struct {
}

func (h *Handler) OnMessage(pck *brokers.Package) {
	fmt.Println(pck.Payload)
	pck.Ack()
}

func (s *MessengerTestSuite) TestRegisterHandler(c *C) {
	msger := NewMessenger("TestRegisterHandler")
	msger.RegisterHandler(&Handler{})
	msger.Send("Test Register message handler")
}

func (s *MessengerTestSuite) TestOnMessageWith(c *C) {
	msger := NewMessenger("TestOnMessageWith")
	crt := make(chan bool)
	msger.RegisterHandler(OnMessageWith(func(pck *brokers.Package) {
		c.Check(pck.Payload, Equals, "Test on message with")
		crt <- true
		pck.Ack()
	}))
	msgerClient := NewMessenger("TestOnMessageWith")
	msgerClient.Send("Test on message with")
	msgerClient.Send("Test on message with")
	<-crt
	<-crt
}

func (s *MessengerTestSuite) TestHandlerResponse(c *C) {

	msger := NewMessenger("TestHandlerResponse")
	crt := make(chan bool)
	msger.RegisterHandler(OnMessageWith(func(pck *brokers.Package) {
		c.Check(pck.Payload, Equals, "Test handler response")
		pck.Response("OK")
		pck.Ack()
	}))
	msgerClient := NewMessenger("TestHandlerResponse")

	msgerClient.Request("Test handler response", 10).
		Then(func(data interface{}) (interface{}, error) {
			p, ok := data.(*brokers.Package)
			if !ok {
				return nil, errors.New("Not a suitable package")
			}
			c.Check(p.Payload, Equals, "OK")
			crt <- true
			return data, nil
		}, func(err error) error {
			c.Fail()
			crt <- false
			return err
		})

	c.Assert(<-crt, Equals, true)

}

func (s *MessengerTestSuite) TestHandlerFailed(c *C) {
	msger := NewMessenger("TestHandlerFailed")
	crt := make(chan bool)
	msger.RegisterHandler(OnMessageWith(func(pck *brokers.Package) {
		pck.Fail()
		crt <- false
	}))
	msger.Send("TestHandler failed")
	c.Assert(<-crt, Equals, false)
}
