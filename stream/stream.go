package stream

import "github.com/jc-m/rhub/modules"


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
	Module *modules.Module
	Upstream []string // Id of upstream modules
	Downstream []string // Id of downstream modules
}

type Stream struct {
	NodeMap map[*Node]bool
}

func NewStream() *Stream {
	return &Stream {
		NodeMap: make(map[*Node]bool),
	}
}

func NewNode(id string, mod *modules.Module) *Node {
	return &Node {
		Id: id,
		Module: mod,
		Upstream: make([]string,0),
		Downstream: make([]string,0),
	}
}

func (s *Stream) Push(lower, upper modules.Module) error {


	if _, err := lower.AddChannels(modules.QueuePair{Read:in, Write:out, Ctl:ctl}); err != nil {
		return err
	}
	if _, err := upper.AddChannels(modules.QueuePair{Read:out, Write:in, Ctl:ctl}); err != nil {
		return err
	}

	if err := lower.Open(); err != nil {
		return err
	}
	if err := upper.Open(); err != nil {
		return err
	}

	return nil
}