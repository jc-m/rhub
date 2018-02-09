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

const (
	STATE_STOPPED = iota
	STATE_STARTED
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

type ModuleConfig map[string]string

type InitFunc func(ModuleConfig) (Module, error)

var ModuleMap map[string]InitFunc

type Module interface {
	Open() error
	GetName() string
	GetUUID() string
	GetType() int
	GetQueues() *QueuePair
	ConnectQueuePair(*QueuePair) error
	Close()
}

func init() {
	ModuleMap = make(map[string]InitFunc)
}

func Register(name string, f InitFunc ) error {
	if _, exists := ModuleMap[name]; exists {
		panic("Module registered multiple time")
	}

	ModuleMap[name] = f

	return nil
}

func GetModule(name string, conf ModuleConfig) (Module, error) {
	if f, ok := ModuleMap[name]; ok {
		return f(conf)
	}
	return nil, fmt.Errorf("Unknown module %s", name)
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