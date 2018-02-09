package serial

const (
	PARITY_ODD = iota
	PARITY_EVEN
	PARITY_NONE
)

const (
	MAX_QUEUES=0
	MAX_CONNECTED=1
)

const (
	PORT_CLOSED = iota
	PORT_OPEN
)