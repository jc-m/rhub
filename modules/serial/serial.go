package serial

import (
	"log"
	"io"
	. "github.com/jc-m/rhub/modules"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"strconv"
)



type SerialConfig struct {
	Port string
	Baud int
	Parity uint
	DataBits uint
	StopBits int
	RtsCts bool
}

type serialPort struct {
	config     SerialConfig
	queue      *QueuePair
	state      int
	port       io.ReadWriteCloser
	uuid       string
	portState  int
}


func init() {
	Register("serial", NewSerial)
}

func (m *serialPort) portOpen() error {

	port, err := openInternal(m.config)

	if err != nil {
		return err
	}
	m.port = port
	m.portState = PORT_OPEN
	return nil
}

func (m *serialPort) serialWrite(buff []byte) int {
	n, err := m.port.Write(buff)
	if err != nil {
		log.Fatalf("port.Write: %v", err)
	}
	log.Printf("[DEBUG] SerialClient: Wrote %d bytes", n)
	return n
}

func (m *serialPort) serialClose() {
	log.Print("[DEBUG] SerialClient: closing port")
	m.portState = PORT_CLOSED
	m.port.Close()
}

func (m *serialPort) sendloop() {
	for {
		select {
		case r := <-m.queue.Read:
			if m.portState == PORT_CLOSED {
				log.Print("[DEBUG] SerialClient: closing send loop")
				return
			}
			log.Printf("[DEBUG] SerialClient: Writing %s", string(r.Body))

			m.serialWrite(r.Body)
		}
	}
}

func (m *serialPort) ctlloop() {
	select {
	case  <-m.queue.Ctl:
		log.Print("[DEBUG] SerialClient: ctlloop close")
		break
	}
	m.Close()
}


func (m *serialPort) receiveloop() {
	buffer := make([]byte, 1024)
	for {
		log.Print("[DEBUG] SerialClient: Reading")

		n, err := m.port.Read(buffer)
		log.Printf("[DEBUG] SerialClient: Received %d bytes", n)

		if n > 0 {
			b := make([]byte, n)
			copy(b, buffer[:n])
			log.Printf("[DEBUG] SerialClient: Received %s", string(b))

			m.queue.Write <- Message{Id:m.config.Port, Body:b}
		}
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				log.Print("[DEBUG] SerialClient: EOF")
			}
			if err != nil {
				log.Printf("[ERROR] SerialClient: Read error : %s", err)
				break
			}
		}

		if  n == 0 && err == nil {
			log.Print("[DEBUG] SerialClient: Null Read")
			break
		}

		if m.portState == PORT_CLOSED {
			break
		}
	}
	log.Print("[DEBUG] SerialClient: Terminating Receiving loop")
	m.serialClose()
	close(m.queue.Ctl)
}

func (m *serialPort) GetName() string {
	return m.config.Port
}

func (m *serialPort) GetUUID() string {
	return m.uuid
}
func (m *serialPort) GetType() int {
	return DRIVER
}


func (m *serialPort) Close() {
	log.Print("[DEBUG] SerialClient: Closing")
	m.state = STATE_STOPPED
	close(m.queue.Read)
	close(m.queue.Write)
}

// Reads from a pair of downstream
// and write to a serial Port

func (m *serialPort) Open()  error {
	if m.state == STATE_STARTED {
		panic("Serial: already started")
	}
	if m.queue == nil {
		return fmt.Errorf("Need to create Add queue first")
	}
	if err := m.portOpen(); err != nil {
		log.Printf("[ERROR] SerialClient: Cannot open port: %s", err)
		return err
	}
	m.state = STATE_STARTED
	go m.receiveloop()
	go m.sendloop()
	go m.ctlloop()

	return nil
}

func (m *serialPort) ConnectQueuePair(q *QueuePair) error  {
	return fmt.Errorf("Not supported")
}

func (m *serialPort) GetQueues() *QueuePair {
	return m.queue
}

func getConfig(conf ModuleConfig) (SerialConfig, error) {
	log.Printf("getConfig %v", conf)
	c := SerialConfig{}
	if port, ok := conf["port"]; ok {
		c.Port = port
	}
	if baud, ok := conf["baud"]; ok {
		if v, err := strconv.Atoi(baud); err == nil {
			c.Baud = v
		} else {
			return c, fmt.Errorf("Invalid Baud value")
		}
	}
	if sb, ok := conf["stop_bits"]; ok {
		if v, err := strconv.Atoi(sb); err == nil {
			c.StopBits = v
		} else {
			return c, fmt.Errorf("Invalid stop_bits value")
		}
	}
	if tap, ok := conf["rts_cts"]; ok {
		switch tap{
		case "true":
			c.RtsCts = true
		case "false":
			c.RtsCts = true
		default:
			return c, fmt.Errorf("invalid rts_cts value")
		}
	}

	return c, nil
}
func NewSerial(conf ModuleConfig) (Module, error) {
	out := &serialPort{}
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
	out.state = STATE_STOPPED
	out.portState = PORT_CLOSED
	if c, err := getConfig( conf ); err != nil {
		return nil, err
	} else {
		out.config = c
	}

	return out, nil
}