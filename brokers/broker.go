package brokers

import (
	rediscfg "github.com/do4way/ivynet-eimb-go/redisConfig"
	redis "gopkg.in/redis.v5"
)

//Broker an EIM broker object.
type Broker struct {
	redisClient *redis.Client
	Name        string
}

//New to create a new EIB broker
func New(name string) *Broker {
	cfg := rediscfg.Read()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return createInstance(name, redisClient)
}

func createInstance(name string, redisClient *redis.Client) *Broker {
	b := &Broker{
		redisClient: redisClient,
		Name:        name,
	}
	b.redisClient.SAdd(masterKey(), name)
	return b
}

//Ping ping to redis server.
func (broker *Broker) Ping() (string, error) {

	return broker.redisClient.Ping().Result()
}

//GetOutputChannel ...
func (broker *Broker) GetOutputChannel() *OutputChannel {
	return &OutputChannel{
		channel: channel{
			channelKey: inputChannelKey(broker.Name),
			redis:      broker.redisClient,
		},
	}
}

//GetInputChannel ...
func (broker *Broker) GetInputChannel(cName string) *InputChannel {
	return &InputChannel{
		channel: channel{
			channelKey: inputChannelKey(broker.Name),
			redis:      broker.redisClient,
		},
		workingChannelKey: workingChannelKey(broker.Name, cName),
		failedChannelKey:  failChannelKey(broker.Name, cName),
	}
}

//GetResponseChannel ...
func (broker *Broker) GetResponseChannel(rid string) *ResponseChannel {
	return &ResponseChannel{
		channel: channel{
			channelKey: inputChannelKey(broker.Name),
			redis:      broker.redisClient,
		},
		responseChannelKey: responseChannelKey(broker.Name, rid),
	}
}
