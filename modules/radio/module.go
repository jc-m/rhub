package radio

import (
	"fmt"
	"github.com/jc-m/rhub/modules"
	"bytes"
	"encoding/gob"
	"log"
	"github.com/jc-m/rhub/modules/radio/rigs"
	"github.com/jc-m/rhub/modules/radio/rigs/ft991a"
)

type rig struct {
	connected  *modules.QueuePair
	queue      *modules.QueuePair
	driver     rigs.Rig
}



func getBytes(msg interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(msg)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *rig) upstreamLoop() {
	for {
		select {
		case buff := <-m.connected.Read:
			if resp, err := m.driver.OnCatUpStream(string(buff.Body)); err == nil {
				if buff, err := getBytes(resp); err == nil {
					m.queue.Write <- modules.Message{Body:buff}
				} else {
					log.Print("[DEBUG] FT991A: Error Encoding message")
				}
			} else{
				log.Printf("[DEBUG] FT991A: Error processing message: %s", string(buff.Body))
			}
			log.Printf("[DEBUG] RadioModel: Sent command: %s", string(buff.Body))
		case <-m.connected.Ctl:
			log.Print("[DEBUG] RadioModel: Terminating upper loop")
			return
		}
	}
}


func (m *rig) Open()  error {

	go m.upstreamLoop()

	return nil

}

func (m *rig) GetName() string {
	return ""
}

func (m *rig) GetType() int {
	return modules.MUX
}

func (m *rig) CreateQueue() (*modules.QueuePair, error)  {

	if m.queue != nil {
		return nil, fmt.Errorf("Module supports only one queue")
	}
	m.queue = &modules.QueuePair{
		Read:  make(chan modules.Message),
		Write: make(chan modules.Message),
		Ctl:   make(chan bool),

	}
	return m.queue, nil
}

func (m *rig) GetQueues() []*modules.QueuePair {
	return nil
}

func (m *rig) ConnectQueuePair(q *modules.QueuePair) error  {
	if m.connected != nil {
		return fmt.Errorf("Module supports only one connection")
	}
	m.connected = q
	return nil
}

func (m *rig) Close() {
	return
}


func New() modules.Module {

	return &rig {
		driver: ft991a.New(),
	}
}