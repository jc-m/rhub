package serial

const (
	STATE_CLOSED = iota
	STATE_OPEN
)

const (
	PARITY_ODD = iota
	PARITY_EVEN
	PARITY_NONE
)

const (
	MAX_QUEUES=0
	MAX_CONNECTED=1
)