package modules

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
	Read  chan Message
	Write chan Message
	Ctl chan bool
}


type Module interface {
	Open() error
	GetName() string
	GetType() int
	CreateQueue() (*QueuePair, error)
	GetQueues() []*QueuePair
	ConnectQueuePair(*QueuePair) error
	Close()
}

