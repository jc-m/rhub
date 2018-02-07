package utils

import (
	"github.com/jc-m/rhub/modules"
	"fmt"
	"log"
	"bytes"
	"github.com/hashicorp/go-uuid"
)

// pipe with tap

type pipeModule struct {
	queue *modules.QueuePair
	connected  []*modules.QueuePair
	buffer *bytes.Buffer
	uuid string
}

func (r *pipeModule) GetType() int {
	return modules.MUX
}

func (r *pipeModule) GetName() string {
	return ""
}

func (r *pipeModule) GetUUID() string {
	return r.uuid
}

func (r *pipeModule) ConnectQueuePair(q *modules.QueuePair) error  {
	if len(r.connected) >1  {
		return fmt.Errorf("Module supports only two connection")
	}
	r.connected = append(r.connected, q)
	return nil
}

func (r *pipeModule) GetQueues() *modules.QueuePair {
	return r.queue
}


func (r *pipeModule) pipeLoop() {
	for {
		select {
		case msg := <-r.connected[0].Write:
			log.Printf("[DEBUG] Pipe: Received on 0 %+v", msg)
			r.connected[1].Read <- msg
		case msg := <-r.connected[1].Write:
			log.Printf("[DEBUG] Pipe: Received on 1 %+v", msg)
			r.connected[0].Read <- msg
		case <-r.queue.Ctl:
			log.Print("[DEBUG] Pipe: Terminating  loop")
			return
		}
	}
}

func (r *pipeModule) Open()  error {
	// create one command loop per channel pair
	go r.pipeLoop()
	return nil
}

func (r *pipeModule) Close() {
}

func NewPipe() modules.Module {

	q := &modules.QueuePair{
		Read:  make(chan modules.Message),
		Write: make(chan modules.Message),
		Ctl:   make(chan bool),

	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}
	return &pipeModule{
		buffer: bytes.NewBuffer([]byte{}),
		queue: q,
		uuid: id,
	}
}
