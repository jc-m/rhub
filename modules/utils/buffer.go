package utils

import (
	"github.com/jc-m/rhub/modules"
	"fmt"
	"log"
	"bytes"
	"io"
)

type cmdBuffer struct {
	config Config
	queue      *modules.QueuePair
	connected  *modules.QueuePair
	buffer *bytes.Buffer
	name string
}

type Config struct {
	Delimiter byte
}

func (r *cmdBuffer) GetType() int {
	return modules.MODULE
}

func (r *cmdBuffer) GetName() string {
	return r.name
}

func (r *cmdBuffer) CreateQueue() (*modules.QueuePair, error)  {

	if r.queue != nil {
		return nil, fmt.Errorf("Module supports only one queue")
	}
	r.queue = &modules.QueuePair{
		Read:  make(chan modules.Message),
		Write: make(chan modules.Message),
		Ctl:   make(chan bool),

	}
	return r.queue, nil
}

func (r *cmdBuffer) ConnectQueuePair(q *modules.QueuePair) error  {
	if r.connected != nil {
		return fmt.Errorf("Module supports only one connection")
	}
	r.connected = q
	return nil
}

func (r *cmdBuffer) GetQueues() []*modules.QueuePair {
	return []*modules.QueuePair{r.queue}
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
								log.Printf("[ERROR] RadioModel: Error reading buffer: %s", err)
							}
							break
						}
						log.Printf("[DEBUG] RadioModel: Received command: %s", string(l))

						r.connected.Write <- modules.Message{Body: l}
						log.Printf("[DEBUG] RadioModel: Sent command: %s", string(l))
					}
				}
		case <-r.queue.Ctl:
			log.Print("[DEBUG] RadioModel: Terminating upper loop")
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
			log.Printf("[DEBUG] RadioModel: Sent command: %s", string(buff.Body))
		case <-r.connected.Ctl:
			log.Print("[DEBUG] RadioModel: Terminating upper loop")
			return
		}
	}
}

func (r *cmdBuffer) Open()  error {
	// create one command loop per channel pair
	go r.downstreamLoop()
	return nil
}

func (r *cmdBuffer) Close() {
}

func NewCmdBuffer(conf Config) modules.Module {

	return &cmdBuffer{
		buffer: bytes.NewBuffer([]byte{}),
		config: conf,
	}
}
