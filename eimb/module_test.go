package eimb

import (
	"testing"

	. "gopkg.in/check.v1"
)

func TestApis(t *testing.T) {
	TestingT(t)
}

//MasterKey : a set that contains all the brokers' name
func masterKey() string {
	return "eimbs::master"
}

//InputChannelKey : broker input channel's key name.
func inputChannelKey(bName string) string {
	return "eimbs::" + bName
}

func workingPrefix(bName string) string {
	return inputChannelKey(bName) + "::working"
}

//WorkingChannelKey : working channel's key name
func workingChannelKey(qName string, cName string) string {
	return workingPrefix(qName) + "::" + cName
}

func failChannelKey(qName string, cName string) string {
	return workingChannelKey(qName, cName) + "::failed"
}

//ResponseChannelKey : response channel's key name
func responseChannelKey(inputChannelKey string, requestID string) string {
	return inputChannelKey + "::" + requestID
}
