package serial

import (
	"github.com/jacobsa/go-serial/serial"
	"log"
	"io"
	"github.com/jc-m/rhub/modules"
	"fmt"
)




type SerialConfig struct {
	Port string
	Baud uint
	Parity uint
	DataBits uint
	StopBits uint
}

type serialPort struct {
	config SerialConfig
	channels []modules.Channels
	state byte
	port io.ReadWriteCloser
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
		case r := <-m.channels[0].In:
			if m.state == STATE_CLOSED {
				return
			}
			m.serialWrite(r.Body)
		case <-m.channels[0].Ctl:
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

			m.channels[0].Out <- modules.Message{Id:m.config.Port, Body:b}
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

		// TODO address if n == 0

		if m.state == STATE_CLOSED {
			break
		}
	}
	log.Print("[DEBUG] SerialClient: Terminating Receiving loop")
	m.channels[0].Ctl <- true
}

func (m *serialPort) GetName() string {
	return m.config.Port
}

func (m *serialPort) GetType() int {
	return modules.DRIVER
}


func (m *serialPort) Close() {
	log.Print("[DEBUG] SerialClient: Closing")
	m.state = STATE_CLOSED
	close(m.channels[0].Ctl)
	m.port.Close()

}

// Reads from a pair of channels
// and write to a serial Port
// TODO provide the notion of command with a separator
func (m *serialPort) Open()  error {

	if len(m.channels) == 0 {
		return fmt.Errorf("Need to provide one pair of channels")
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
func (m *serialPort) AddChannels(channels modules.Channels) ([]modules.Channels, error)  {
	if len(m.channels) > 0 {
		return nil, fmt.Errorf("Module supports only on pair of channels")
	}
	m.channels = append(m.channels, channels)
	return m.channels, nil
}

func NewSerial(c SerialConfig) modules.Module {
	return &serialPort{
		config: c,
		channels: make([]modules.Channels,0),
		state: STATE_CLOSED,
	}
}