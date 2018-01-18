package server

import (
	"testing"
	"github.com/jc-m/rhub/modules/serial"

)

func TestStreamConnect(t *testing.T) {

	c := serial.SerialConfig{
		Port: "/dev/ttys005",
		Baud: 38400,
		DataBits:8,
		StopBits:2,
	}

	serPort := serial.NewSerial(c)
	ptyPort := serial.NewPty()

	if err := Push(serPort,ptyPort); err != nil {
		t.Fatal(err)

	}

	// need to open the pty and test

}
