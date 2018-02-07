package modules

import (
	"fmt"
)

// DRIVER modules only have one pair of channels, and are communicating with the network or a device on the other side.
const (
	DRIVER = iota
	MODULE
	MUX
)
const (
	M_DATA = iota
)

type Message struct {
	Id   string
	Body []byte
	Type int
}
type QueuePair struct {
	Owner string
	Read  chan Message
	Write chan Message
	Ctl chan bool
}


type Module interface {
	Open() error
	GetName() string
	GetUUID() string
	GetType() int
	GetQueues() *QueuePair
	ConnectQueuePair(*QueuePair) error
	Close()
}

type NotImplemented struct {}

func (m *NotImplemented) Open()  error {
	return fmt.Errorf("Not implemented")

}

func (m *NotImplemented) GetName() string {
	return ""
}

func (m *NotImplemented) GetType() int {
	return DRIVER
}

func (m *NotImplemented) CreateQueue() (*QueuePair, error)  {
	return nil, fmt.Errorf("Not implemented")
}

func (m *NotImplemented) GetQueues() *QueuePair {
	return nil
}

func (m *NotImplemented) ConnectQueuePair(q *QueuePair) error  {
	return fmt.Errorf("Not implemented")
}

func (m *NotImplemented) Close() {
	return
}