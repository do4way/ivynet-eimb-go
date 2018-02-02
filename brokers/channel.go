package brokers

import (
	"encoding/json"
	"time"

	redis "gopkg.in/redis.v5"

	"github.com/do4way/ivynet-eimb-go/promises"
)

type channel struct {
	channelKey string
	redis      *redis.Client
}

//InputChannel ..
type InputChannel struct {
	channel
	workingChannelKey string
	failedChannelKey  string
}

//AsyncRead ...
func (input *InputChannel) AsyncRead() *promises.Promise {
	p := promises.New()
	go func() {
		pack, error := input.Read(10)
		if error != nil {
			p.Reject(error)
			return
		}
		p.Resolve(pack)
	}()
	return p
}

func (input *InputChannel) Read(to time.Duration) (*Package, error) {

	answer := input.redis.BRPopLPush(
		input.channelKey,
		input.workingChannelKey,
		to*time.Second,
	)
	if answer.Err() != nil {
		return nil, answer.Err()
	}
	if answer.Val() == "" {
		return nil, nil
	}

	pack, err := parseRedisAnswer(answer)
	if err != nil {
		return nil, err
	}

	pack.channelKey = input.channelKey
	pack.workingChannelKey = input.workingChannelKey
	pack.failedChannelKey = input.failedChannelKey
	pack.redis = input.redis
	return pack, nil
}

//OutputChannel ...
type OutputChannel struct {
	channel
}

//Write ...
func (output *OutputChannel) Write(pack *Package) error {
	lpush := output.redis.LPush(output.channelKey, pack.getString())
	return lpush.Err()
}

//ResponseChannel ...
type ResponseChannel struct {
	channel
	responseChannelKey string
}

func (response *ResponseChannel) Read(to time.Duration) (*Package, error) {
	answer := response.redis.BRPop(
		to*time.Second,
		response.responseChannelKey,
	)

	return parseRedisAnswers(answer)
}

func parseRedisAnswer(answer *redis.StringCmd) (*Package, error) {

	if answer.Err() != nil {
		return nil, answer.Err()
	}
	p, err := unmarshalPackage(answer.Val())
	if err != nil {
		return nil, err
	}

	return p, nil
}

func unmarshalPackage(input string) (*Package, error) {

	p := &Package{}
	err := json.Unmarshal([]byte(input), p)
	if err != nil {
		return nil, err
	}
	p.origStr = input
	return p, nil
}

func parseRedisAnswers(answer *redis.StringSliceCmd) (*Package, error) {

	if answer.Err() != nil {
		return nil, answer.Err()
	}
	val := answer.Val()
	pack, err := unmarshalPackage(val[1])
	if err != nil {
		return nil, err
	}
	return pack, nil
}
