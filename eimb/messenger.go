package eimb

import (
	"fmt"
	"sync"
	"time"

	"github.com/do4way/ivynet-eimb-go/brokers"
	"github.com/do4way/ivynet-eimb-go/promises"
	"github.com/google/uuid"
)

const (
	handlersDefaultLength = 5
)

//Messenger ...
type Messenger struct {
	sync.Mutex
	broker       *brokers.Broker
	handlers     []MessageHandler
	onMonitoring bool
	signal       byte
}

//MessageHandler ..
type MessageHandler interface {
	OnMessage(pck *brokers.Package)
}

//OnMessageWith ..
type OnMessageWith func(pck *brokers.Package)

//OnMessage ...
func (messageWith OnMessageWith) OnMessage(pck *brokers.Package) {
	messageWith(pck)
}

//NewMessenger ...
func NewMessenger(topic string) *Messenger {
	return &Messenger{
		broker:       brokers.New(topic),
		handlers:     make([]MessageHandler, 0, handlersDefaultLength),
		onMonitoring: false,
		signal:       0x00,
	}
}

//Send ...
func (m *Messenger) Send(msg interface{}) error {

	pck, err := brokers.BuildPackage(msg)
	if err != nil {
		return err
	}

	return m.send(pck)
}

func (m *Messenger) Read() (*brokers.Package, error) {
	return m.broker.GetInputChannel(m.botName()).Read(0)
}

func (m *Messenger) send(pck *brokers.Package) error {
	return m.broker.GetOutputChannel().Write(pck)
}

//Request ...
func (m *Messenger) Request(msg interface{}, t time.Duration) *promises.Promise {

	promise := promises.New()
	pck, bErr := brokers.BuildPackage(msg)
	if bErr != nil {
		promise.Reject(bErr)
		return promise
	}
	sErr := m.send(pck)
	if sErr != nil {
		promise.Reject(sErr)
		return promise
	}
	responseChannel := m.broker.GetResponseChannel(pck.ID)
	go func() {
		data, err := responseChannel.Read(t)
		if err != nil {
			promise.Reject(err)
			return
		}
		promise.Resolve(data)
	}()
	return promise
}

//RegisterHandler ...
func (m *Messenger) RegisterHandler(h MessageHandler) {

	if m.isHandlerRegistered(h) {
		return
	}
	m.addHandler(h)
	if m.onMonitoring {
		return
	}
	m.Lock()
	if !m.onMonitoring {
		m.startMonitor()
		m.onMonitoring = true
	}
	m.Unlock()
}

func (m *Messenger) isHandlerRegistered(h MessageHandler) bool {
	for _, v := range m.handlers {
		if v == h {
			return true
		}
	}
	return false
}

func (m *Messenger) addHandler(h MessageHandler) {
	m.Lock()
	m.handlers = append(m.handlers, h)
	m.Unlock()
}

func (m *Messenger) startMonitor() {
	go func() {
		for m.signal == 0x00 {
			pck, err := m.Read()
			if err != nil {
				//TODO: need an smart way to handle the error
				fmt.Println(err)
				continue
			}
			for _, h := range m.handlers {
				h.OnMessage(pck)
			}
		}
	}()
}

func (m *Messenger) botName() string {
	return m.broker.Name + "::bot::" + uuid.New().String()
}
