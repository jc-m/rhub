package network

import (
	. "github.com/jc-m/rhub/modules"
	"fmt"
	"net"
	"log"
	"github.com/hashicorp/go-uuid"
)
const (
	CNX_OPEN = iota
	CNX_CLOSE
)
type tcpServ struct {
	queue   *QueuePair
	address string
	listener net.Listener
	state byte
	uuid string
	connected bool
}

type connection struct {
	conn net.Conn
	server *tcpServ
	state int
}

func init() {
	log.Printf("Init tcp_server")
	Register("tcp_server", NewTCPServ)
}

func (c *connection) receiveloop(){
	buffer := make([]byte, 1024)
	log.Print("[DEBUG] TCPServer: Starting Receive loop")

	for {
		if c.state == CNX_CLOSE {
			break
		}
		n, err := c.conn.Read(buffer)
		if err != nil {
			log.Printf("[DEBUG] TCPServer: Error %+v", err)
			break
		}
		if n > 0 {
			b := make([]byte, n)
			copy(b, buffer[:n])

			log.Printf("[DEBUG] TCPServer: Received %s", string(b))
			c.server.queue.Write <- Message{Id:c.conn.RemoteAddr().String(), Body:b}

		} else {
			log.Printf("[DEBUG] TCPServer: 0 read")
		}
	}
	log.Print("[DEBUG] TCPServer: Closing receiveloop")

	if c.state == CNX_OPEN {
		c.state = CNX_CLOSE
		c.conn.Close()
		c.server.connected = false
	}
}

func (c *connection) sendloop() {
	log.Print("[DEBUG] TCPServer: Starting Sending loop")

	for {
		select {
		case r := <-c.server.queue.Read:
			// TODO server address does not have the right address
			if r.Id == c.server.address {
			}
			if c.state == CNX_CLOSE {
				goto close
			}
			n, err := c.conn.Write(r.Body)
			if err != nil {
				panic(err)
			}
			log.Printf("[DEBUG] TCPServer: Sent %d bytes", n)
		case <-c.server.queue.Ctl:
			log.Print("[DEBUG] TCPServer: Terminating Sending loop")
			return
		}
	}
close:
	if c.state == CNX_OPEN {
		c.state = CNX_CLOSE
		c.conn.Close()
		c.server.connected = false
	}
}
func (t *tcpServ) GetType() int {
	return DRIVER
}

func (t *tcpServ) GetName() string {
	return t.address
}

func (t *tcpServ) GetUUID() string {
	return t.uuid
}

func (t *tcpServ) Close() {

}

func (t *tcpServ) Open()  error {
	if t.state == STATE_STARTED {
		panic("TCPServer: already started")
	}
	t.state = STATE_STARTED
	go func(server *tcpServ) error {

		log.Printf("[DEBUG] TCPServer: Listening on %s", server.address)

		l, err := net.Listen("tcp", server.address)
		if err != nil {
			log.Printf("[ERROR] TCPServer: listening: %s", err)
			return err
		}
		server.listener = l

		for {
			log.Print("[DEBUG] TCPServer: Accepting")
			conn, err := l.Accept()
			if err != nil {
				log.Fatal(err)
			}
			if t.connected {
				log.Print("[DEBUG] TCPServer: Already connected")
				conn.Close()
				continue
			}
			client := &connection{
				conn:   conn,
				server: server,
			}
			go client.receiveloop()
			go client.sendloop()
			t.connected = true
		}
		return nil
	}(t)

	return nil
}

func (t *tcpServ) ConnectQueuePair(q *QueuePair) error  {
	return fmt.Errorf("Not supported")
}

func (t *tcpServ) GetQueues() *QueuePair {
	return t.queue
}


func NewTCPServ(conf ModuleConfig) (Module, error) {

	q := &QueuePair{
		Read:  make(chan Message),
		Write: make(chan Message),
		Ctl:   make(chan bool),

	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		panic(err)
	}

	return &tcpServ {
		address : conf["address"],
		uuid: id,
		queue: q,
		state: STATE_STOPPED,
		connected: false,
	}, nil
}
