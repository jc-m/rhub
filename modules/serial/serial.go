package serial

import (
	"github.com/jacobsa/go-serial/serial"
	"log"
	"io"
	. "github.com/jc-m/rhub/modules"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"strconv"
)



type SerialConfig struct {
	Port string
	Baud uint
	Parity uint
	DataBits uint
	StopBits uint
}

type serialPort struct {
	config     SerialConfig
	queue      *QueuePair
	state      byte
	port       io.ReadWriteCloser
	uuid       string
}

func init() {
	Register("serial", NewSerial)
}

func (m *serialPort) serialOpen() error {
	var parity serial.ParityMode
	c := m.config

	switch c.Parity {
	case PARITY_EVEN:
		parity = serial.PARITY_EVEN
	case PARITY_ODD:
		parity = serial.PARITY_ODD
	case PARITY_NONE:
		parity = serial.PARITY_NONE
	default:
		parity = serial.PARITY_NONE
	}


	conf := serial.OpenOptions{
		PortName: c.Port,
		BaudRate:c.Baud,
		ParityMode:parity,
		StopBits:c.StopBits,
		DataBits:c.DataBits,
		MinimumReadSize: 1,
		}

	port, err := serial.Open(conf)
	if err != nil {
		log.Printf("[ERROR] SerialClient: Cannot open port %s",err)
		return err
	}
	m.port = port
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

func (m *serialPort) sendloop() {
	for {
		select {
		case r := <-m.queue.Read:
			if m.state == STATE_CLOSED {
				return
			}
			m.serialWrite(r.Body)
		case <-m.queue.Ctl:
			log.Print("[DEBUG] SerialClient: Terminating Sending loop")
			return
		}
	}
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
			log.Printf("[DEBUG] SerialClient: Sending %+v", b)

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

		if err == nil {
			log.Print("[DEBUG] Serial  Null Read")
		}

		if m.state == STATE_CLOSED {
			break
		}
	}
	log.Print("[DEBUG] SerialClient: Terminating Receiving loop")
	m.queue.Ctl <- true
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
	m.state = STATE_CLOSED
	close(m.queue.Ctl)
	m.port.Close()

}

// Reads from a pair of downstream
// and write to a serial Port
// TODO provide the notion of command with a separator
func (m *serialPort) Open()  error {
	if m.queue == nil {
		return fmt.Errorf("Need to create Add queue first")
	}
	if err := m.serialOpen(); err != nil {
		log.Printf("[ERROR] SerialClient: Cannot open port: %s", err)
		return err
	}
	m.state = STATE_OPEN
	go m.receiveloop()
	go m.sendloop()

	return nil
}

func (m *serialPort) ConnectQueuePair(q *QueuePair) error  {
	return fmt.Errorf("Not supported")
}

func (m *serialPort) GetQueues() *QueuePair {
	return m.queue
}

func getConfig(conf ModuleConfig) (SerialConfig, error) {
	c := SerialConfig{}
	if port, ok := conf["port"]; ok {
		c.Port = port
	}
	if baud, ok := conf["baud"]; ok {
		if v, err := strconv.Atoi(baud); err == nil {
			c.Baud = uint(v)
		} else {
			return c, fmt.Errorf("Invalid Baud value")
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
	out.state = STATE_CLOSED

	if c, err := getConfig( conf ); err != nil {
		return nil, err
	} else {
		out.config = c
	}

	return out, nil
}