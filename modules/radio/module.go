package radio

import (
	"fmt"
	. "github.com/jc-m/rhub/modules"
	"bytes"
	"encoding/gob"
	"log"
	"github.com/jc-m/rhub/modules/radio/rigs"
	"github.com/jc-m/rhub/modules/radio/rigs/ft991a"
	"github.com/hashicorp/go-uuid"
)

type rig struct {
	connected  *QueuePair
	queue      *QueuePair
	driver     rigs.Rig
	uuid       string
	state      int
}


func init() {
	Register("radio", New)
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
			if resp, err := m.driver.OnCat(string(buff.Body), rigs.CAT_DIR_UP); err == nil {
				if buff, err := getBytes(resp); err == nil {
					m.queue.Write <- Message{Body:buff}
				} else {
					log.Print("[DEBUG] FT991A: Error Encoding message")
				}
			} else{
				log.Printf("[DEBUG] FT991A: Error processing message: %s", err)
			}
			log.Printf("[DEBUG] RadioModel: Sent command: %s", string(buff.Body))
		case <-m.connected.Ctl:
			log.Print("[DEBUG] RadioModel: Terminating upper loop")
			return
		}
	}
}


func (m *rig) Open()  error {
	if m.state == STATE_STARTED {
		panic("RadioModel: already started")
	}
	m.state = STATE_STARTED
	go m.upstreamLoop()

	return nil

}

func (m *rig) GetName() string {
	return ""
}

func (m *rig) GetUUID() string {
	return m.uuid
}

func (m *rig) GetType() int {
	return MUX
}

func (m *rig) GetQueues() *QueuePair {
	return m.queue
}

func (m *rig) ConnectQueuePair(q *QueuePair) error  {
	if m.connected != nil {
		return fmt.Errorf("Module supports only one connection")
	}
	m.connected = q
	return nil
}

func (m *rig) Close() {
	return
}


func New(conf ModuleConfig) (Module, error) {
	q := &QueuePair{
		Read:  make(chan Message),
		Write: make(chan Message),
		Ctl:   make(chan bool),

	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}

	return &rig {
		queue: q,
		driver: ft991a.New(),
		uuid: id,
		state: STATE_STOPPED,
	}, nil
}