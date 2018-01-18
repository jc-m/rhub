package modules

// DRIVER modules only have one pair of channels, and are communicating with the network or a device on the other side.
const (
	DRIVER = iota
	MUX
)
type Message struct {
	Id   string
	Body []byte
}
type Channels struct {
	In  chan Message
	Out chan Message
	Ctl chan bool
}


type Module interface {
	Open() error
	GetName() string
	GetType() int
	AddChannels(Channels) ([]Channels, error)
	Close()
}

