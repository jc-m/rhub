package utils

import (
	. "github.com/jc-m/rhub/modules"
	"fmt"
	"log"
	"bytes"
	"io"
	"github.com/hashicorp/go-uuid"
)

type cmdBuffer struct {
	config BufferConfig
	queue      *QueuePair
	connected  *QueuePair
	buffer *bytes.Buffer
	name string
	uuid string
	state int
}

type BufferConfig struct {
	Delimiter byte
}

func init() {
	Register("buffer", NewCmdBuffer)
}

func (r *cmdBuffer) GetType() int {
	return MODULE
}

func (r *cmdBuffer) GetName() string {
	return r.name
}

func (r *cmdBuffer) GetUUID() string {
	return r.uuid
}

func (r *cmdBuffer) ConnectQueuePair(q *QueuePair) error  {
	if r.connected != nil {
		return fmt.Errorf("Module supports only one connection")
	}
	r.connected = q
	return nil
}

func (r *cmdBuffer) GetQueues() *QueuePair {
	return r.queue
}


func (r *cmdBuffer) downstreamLoop() {
	for {
		select {
		case buff := <-r.queue.Read:
				r.buffer.Write(buff.Body)
				// if there is any occurrence of the delimiter in the received message
				// start processing the buffer to find all the commands and process them
				if bytes.Contains(buff.Body, []byte{r.config.Delimiter}) {
					for {
						l, err := r.buffer.ReadBytes(r.config.Delimiter)
						if err != nil {
							if err == io.EOF {
								if len(l) > 0 {
									r.buffer.Write(l)
								}
							} else {
								log.Printf("[ERROR] Buffer: Error reading buffer: %s", err)
							}
							break
						}
						log.Printf("[DEBUG] Buffer: Received command: %s", string(l))

						r.connected.Write <- Message{Body: l}
						log.Printf("[DEBUG] Buffer: Sent command: %s", string(l))
					}
				}
		case <-r.queue.Ctl:
			log.Print("[DEBUG] Buffer: Terminating upper loop")
			return
		default:
		}
	}
}

func (r *cmdBuffer) upstreamLoop(pairId int) {
	for {
		select {
		case buff := <-r.connected.Read:
			r.queue.Write <- buff
			log.Printf("[DEBUG] Buffer: Sent command: %s", string(buff.Body))
		case <-r.connected.Ctl:
			log.Print("[DEBUG] Buffer: Terminating upper loop")
			return
		}
	}
}

func (r *cmdBuffer) Open()  error {
	// create one command loop per channel pair
	if r.state == STATE_STARTED {
		panic("Pipe: Already started")
	}
	r.state = STATE_STARTED
	go r.downstreamLoop()
	return nil
}

func (r *cmdBuffer) Close() {
}

func NewCmdBuffer(conf ModuleConfig) (Module, error) {
	q := &QueuePair{
		Read:  make(chan Message),
		Write: make(chan Message),
		Ctl:   make(chan bool),

	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}
	// TODO make this more robust
	c := BufferConfig {
		Delimiter: byte(conf["delimiter"][0]),
	}
	return &cmdBuffer{
		buffer: bytes.NewBuffer([]byte{}),
		queue: q,
		config: c,
		uuid: id,
		state: STATE_STOPPED,
	}, nil
}
