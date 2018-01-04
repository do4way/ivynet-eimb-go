package brokers

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"

	redis "gopkg.in/redis.v5"
)

//Package a data object used for data transfer.
type Package struct {
	ID                string
	Payload           string
	CreatedAt         time.Time
	channelKey        string
	workingChannelKey string
	failedChannelKey  string
	failedAt          time.Time
	redis             *redis.Client
}

func (pack *Package) getString() string {
	json, err := json.Marshal(pack)
	if err != nil {
		log.Printf("Queue failed to marshal content %v [%s]", pack, err.Error())
		return ""
	}
	return string(json)
}

func (pack *Package) getFailedPackage() string {

	p := struct {
		Payload   string
		CreatedAt time.Time
		FailedAt  time.Time
	}{
		Payload:   pack.Payload,
		CreatedAt: pack.CreatedAt,
		FailedAt:  pack.failedAt,
	}
	json, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(json)
}

//Ack to complete the package.
func (pack *Package) Ack() error {
	if pack.redis == nil {
		return errors.New("Package not cantained")
	}
	return pack.redis.LRem(pack.workingChannelKey, 0, pack.getString()).Err()
}

//Fail to abandon the package, and log it into the failed channel.
func (pack *Package) Fail() error {

	if pack.redis == nil {
		return errors.New("Package not cantained")
	}
	rst := pack.redis.LRem(pack.workingChannelKey, 0, pack.getString())
	if rst.Err() != nil {
		return rst.Err()
	}
	pack.failedAt = time.Now()
	return pack.redis.RPush(pack.failedChannelKey, pack.getFailedPackage()).Err()
}

//Response ...
func (pack *Package) Response(msg interface{}) error {
	response, err := BuildPackage(msg)
	if err != nil {
		return err
	}
	return pack.redis.LPush(pack.channelKey+"::"+pack.ID, response.getString()).Err()

}

//BuildPackage ...
func BuildPackage(payload interface{}) (*Package, error) {
	msg, err := marshal(payload)
	if err != nil {
		return nil, err
	}
	pck := &Package{
		ID:        uuid.New().String(),
		Payload:   msg,
		CreatedAt: time.Now(),
	}
	return pck, nil

}

func marshal(payload interface{}) (string, error) {
	msg := ""
	switch payload.(type) {
	case string:
		msg = payload.(string)
	default:
		bytes, err := json.Marshal(payload)
		if err != nil {
			return "", err
		}
		msg = string(bytes)
	}
	return msg, nil
}
