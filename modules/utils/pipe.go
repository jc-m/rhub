package utils

import (
	. "github.com/jc-m/rhub/modules"
	"fmt"
	"log"
	"bytes"
	"github.com/hashicorp/go-uuid"
)

// pipe with tap

type pipeModule struct {
	config PipeConfig
	queue *QueuePair
	connected []*QueuePair
	buffer *bytes.Buffer
	uuid string
}

type PipeConfig struct {
	Tap bool
}

func init() {
	Register("pipe", NewPipe)
}

func (r *pipeModule) GetType() int {
	return MUX
}

func (r *pipeModule) GetName() string {
	return ""
}

func (r *pipeModule) GetUUID() string {
	return r.uuid
}

func (r *pipeModule) ConnectQueuePair(q *QueuePair) error  {
	if len(r.connected) >1  {
		return fmt.Errorf("Module supports only two connection")
	}
	r.connected = append(r.connected, q)
	return nil
}

func (r *pipeModule) GetQueues() *QueuePair {
	return r.queue
}


func (r *pipeModule) pipeLoop() {
	for {
		select {
		case msg := <-r.connected[0].Write:
			log.Printf("[DEBUG] Pipe: Received on 0 %+v", msg)
			r.connected[1].Read <- msg
			if r.config.Tap {
				r.queue.Write <- msg
			}
		case msg := <-r.connected[1].Write:
			log.Printf("[DEBUG] Pipe: Received on 1 %+v", msg)
			r.connected[0].Read <- msg
			if r.config.Tap {
				r.queue.Write <- msg
			}
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

func getConfig(conf ModuleConfig) (PipeConfig, error) {
	c := PipeConfig{}
	if tap, ok := conf["tap"]; ok {
		switch tap{
		case "true":
			c.Tap = true
		case "false":
			c.Tap = true
		default:
			return c, fmt.Errorf("invalid tap value")
		}
	}
	return c, nil
}

func NewPipe(conf ModuleConfig) (Module, error) {
	out := &pipeModule{
		buffer: bytes.NewBuffer([]byte{}),
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
