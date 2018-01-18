package serial

import (
	"testing"
	"github.com/jc-m/rhub/modules"
)

func TestPTY(t *testing.T) {
	in := make(chan modules.Message)
	out := make(chan modules.Message)


	pty := NewPty()

	pty.AddChannels(modules.Channels{In:in, Out:out})

	if err:= pty.Open(); err != nil {
		t.Fatalf("Cannot open port %s", err)
	}

	// Read master.
	x := <- out
	t.Log( x)

}

