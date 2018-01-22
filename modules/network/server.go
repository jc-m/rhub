package network

import (
	"github.com/jc-m/rhub/modules"
	"fmt"
	"net"
	"log"
)

type tcpServ struct {
	queue   *modules.QueuePair
	address string
	listener net.Listener
	state byte
}

type connection struct {
	conn net.Conn
	server *tcpServ
}

func (c *connection) receiveloop(){
	buffer := make([]byte, 1024)
	log.Print("[DEBUG] TCPServer: Starting Receive loop")

	for {
		n, err := c.conn.Read(buffer)
		if err != nil {
			log.Printf("[DEBUG] TCPServer: Error %+v", err)
			break
		}
		if n > 0 {
			b := make([]byte, n)
			copy(b, buffer[:n])

			log.Printf("[DEBUG] TCPServer: Sending %+v", b)
			// needs to have the client address instead
			c.server.queue.Write <- modules.Message{Id:c.conn.RemoteAddr().String(), Body:b}
		}
	}
	log.Print("[DEBUG] TCPServer: Closing receiveloop")

	c.conn.Close()
	c.server.queue.Ctl <- true
}

func (c *connection) sendloop() {
	log.Print("[DEBUG] TCPServer: Starting Sending loop")

	for {
		select {
		case r := <-c.server.queue.Read:
			if r.Id == c.server.address {
				n, err := c.conn.Write(r.Body)
				if err != nil {
					panic(err)
				}
				log.Printf("[DEBUG] TCPServer: Sent %d bytes", n)
			}
		case <-c.server.queue.Ctl:
			log.Print("[DEBUG] TCPServer: Terminating Sending loop")
			return
		}
	}
}
func (t *tcpServ) GetType() int {
	return modules.DRIVER
}

func (t *tcpServ) GetName() string {
	return t.address
}

func (t *tcpServ) Close() {

}

func (t *tcpServ) Open()  error {

	go func(server *tcpServ) error {
		log.Printf("[DEBUG] TCPServer: Listening on %s", server.address)

		l, err := net.Listen("tcp", server.address)
		if err != nil {
			log.Printf("[ERROR] TCPServer: listening: %s", err)
			return err
		}
		server.listener = l

		for {
			conn, _ := l.Accept()
			client := &connection{
				conn:   conn,
				server: server,
			}
			go client.receiveloop()
			go client.sendloop()

		}
	}(t)

	return nil
}

func (t *tcpServ) CreateQueue() (*modules.QueuePair, error)  {
	if t.queue != nil {
		return nil, fmt.Errorf("Module supports only one queue")
	}
	t.queue = &modules.QueuePair{
		Read:  make(chan modules.Message),
		Write: make(chan modules.Message),
		Ctl:   make(chan bool),

	}
	return t.queue, nil
}

func (t *tcpServ) ConnectQueuePair(q *modules.QueuePair) error  {
	return fmt.Errorf("Not supported")
}

func (t *tcpServ) GetQueues() []*modules.QueuePair {
	return []*modules.QueuePair{t.queue}
}


// c is a string in the form of what Go Net package accept for dial
func NewTCPServ(address string) modules.Module {

	return &tcpServ {
		address : address,
	}
}
