package stream

import (
	"testing"
	"github.com/jc-m/rhub/modules/serial"
	"github.com/jc-m/rhub/modules/utils"
)

func TestStream(t *testing.T) {

	c := serial.SerialConfig{
		Port: "/dev/ttys012",
		Baud: 38400,
		DataBits:8,
		StopBits:2,
	}

	serPort := serial.NewSerial(c)
	ptyPort := serial.NewPty()
	pipe := utils.NewPipe()

	s := NewStream()

	if err := s.Push(ptyPort, pipe); err != nil {
		t.Fatal(err)
	}

	if err := s.Push(serPort, pipe); err != nil {
		t.Fatal(err)
	}
	for k := range s.NodeMap {
		t.Logf("Module %+v",k)

	}

	// need to open the pty and test
	s.Start()

	<- pipe.GetQueues().Ctl

}
