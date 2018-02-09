package utils

import (
	. "github.com/jc-m/rhub/modules"
	"fmt"
	"log"
	"github.com/hashicorp/go-uuid"
)


type dummyModule struct {
    echo bool
    queue *QueuePair
	uuid string
	state int
}

func init() {
	Register("dummy", NewDummy)
}

func (r *dummyModule) GetType() int {
	return DRIVER
}

func (r *dummyModule) GetName() string {
	return ""
}

func (r *dummyModule) GetUUID() string {
	return r.uuid
}

func (p *dummyModule) ConnectQueuePair(q *QueuePair) error  {
	return fmt.Errorf("Not supported")
}


func (r *dummyModule) GetQueues() *QueuePair {
	return r.queue
}


func (r *dummyModule) dummyLoop() {
	for {
		select {
		case msg := <-r.queue.Read:
			log.Printf("[DEBUG] Dummy: Received %+v", msg)
			if r.echo {
				r.queue.Write <- msg
			}
		case <-r.queue.Ctl:
			log.Print("[DEBUG] Dummy: Terminating  loop")
			return
		}
	}
}

func (r *dummyModule) Open()  error {
	// create one command loop per channel pair
	if r.state == STATE_STARTED {
		panic("Pipe: Already started")
	}
	r.state = STATE_STARTED
	go r.dummyLoop()
	return nil
}

func (r *dummyModule) Close() {
}

func NewDummy(conf ModuleConfig) (Module, error) {
	out := &dummyModule{
		state: STATE_STOPPED,
	}
	if conf["echo"] == "true" {
		out.echo = true
	}
	out.queue = &QueuePair{
		Read:  make(chan Message),
		Write: make(chan Message),
		Ctl:   make(chan bool),

	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}
	out.uuid = id
	return out, nil
}
