package stream

import (
	"github.com/jc-m/rhub/modules"
	"fmt"
	"log"
)


/*

    +--------------------------------+
    |       Generic Radio            |
    +--------------------------------+
       |            |          |
    CMDBuffer    FT991A      FT991A
       |            |          |
      PTY       CMDBuffer    Serial
                    |
                TCPServer

    Nil             <- PTY        -> CMDBuffer
    Nil             <- TCPServer  -> CMDBuffer
    Nil             <- Serial     -> Radio
    PTY             <- CMDBuffer  -> Radio
    TCPServer       <- CMDBuffer  -> Radio

    (CMDBuffer,CMDBuffer) <- Radio      -> Serial

	s := Stream.New()
	r := NewRadio()
	c1 := NewCmdBuffer()
    s.Push(c1, r)
	p := NewPTY()
    s.Push(p, c1)
	c2 := NewCmdBuffer()
	s.Push(c2, r)
	t := NewTcpServer()
	s.Push(t, c2)
	u := NewSerial()
	s.Push(u, r)

	YAML
		- Serial
			module : Serial
			config :
				portname,
				speed,
		- Radio
			module : Radio
			downstream : Serial
        - TCPServer
			module : TCPServer
			downstream : C1
        - PTY
			module : PTY
			downstream : C2
        - C1
			module : CMDBuffer
			downstream : Radio
        - C2
			module : CMDBuffer
			downstream : Radio

*/

type Node struct {
	Id string
	Module modules.Module
	Upstream []string // Id of connected upstream modules
}

type Stream struct {
	NodeMap map[*Node]bool
	Index   map[string]*Node
	vertices map[string][]string
}

func NewStream() *Stream {
	return &Stream {
		NodeMap: make(map[*Node]bool),
		Index:make(map[string]*Node),
		vertices: make(map[string][]string),
	}
}

func NewNode(id string, mod modules.Module) *Node {
	return &Node {
		Id: id,
		Module: mod,
		Upstream: make([]string,0),
	}
}

func (n *Node) Connect(upper *Node) error {
	log.Printf("[DEBUG] Stream: Connecting : %s and %s", upper.Id, n.Id)
	if err := upper.Module.ConnectQueuePair(n.Module.GetQueues()); err != nil {
		return err
	}
	n.Upstream = append(n.Upstream, upper.Id)
	return nil
}

func (n *Node) Start() error {
	log.Printf("[DEBUG] Stream: Starting %s", n.Id)

	return n.Module.Open()
}

func (s *Stream) AddNode(id string, mod modules.Module) *Node {
	if n, ok := s.Index[mod.GetUUID()]; ok {
		return n
	} else {
		log.Printf("[DEBUG] Stream: Adding Node : %s", mod.GetUUID())
		n = NewNode(id, mod)
		s.NodeMap[n] = true
		s.Index[n.Id] = n
		return n
	}
}

func (s *Stream) Push(lower, upper modules.Module) error {

	if lower == nil {
		return fmt.Errorf("Lower module is required")
	}
	lowerNode := s.AddNode(lower.GetUUID(), lower)

	if upper == nil {
		// done
		return nil
	}

	upperNode := s.AddNode(upper.GetUUID(), upper)

	return lowerNode.Connect(upperNode)
}

func (s *Stream) startNode(node *Node) error {

	if err := node.Start(); err != nil {
		return err
	}
	for _, n := range s.vertices[node.Id] {
		s.startNode(s.Index[n])
	}
	return nil
}

func (s *Stream) Start() error {
	// Use graph traversal to open modules in the dependency order.
	var n *Node
	var root *Node

	// find the root

	for n = range s.NodeMap {
		if len(n.Upstream) == 0 {
			root = n
		} else {
			for _, ups := range n.Upstream {
				s.vertices[ups] = append(s.vertices[ups], n.Id)
			}
		}
	}
	if root == nil {
		return fmt.Errorf("No Root")
	}
	log.Printf("[DEBUG] Stream: Root : %s", root.Id)

	return s.startNode(root)
}