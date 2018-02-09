package network

import (
	. "github.com/jc-m/rhub/modules"
	"fmt"
	"net"
	"log"
	"github.com/hashicorp/go-uuid"
)

type tcpServ struct {
	queue   *QueuePair
	address string
	listener net.Listener
	state byte
	uuid string
}

type connection struct {
	conn net.Conn
	server *tcpServ
}

func init() {
	log.Printf("Init tcp_server")
	Register("tcp_server", NewTCPServ)
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
			c.server.queue.Write <- Message{Id:c.conn.RemoteAddr().String(), Body:b}
		} else {
			log.Printf("[DEBUG] TCPServer: 0 read")
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
			// TODO server address does not have the right address
			if r.Id == c.server.address {
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
			client := &connection{
				conn:   conn,
				server: server,
			}
			go client.receiveloop()
			go client.sendloop()

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
	}, nil
}
