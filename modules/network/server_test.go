package network

import (
	"testing"
	"github.com/jc-m/rhub/modules"
)

func TestTcpServer(t *testing.T) {
	in := make(chan modules.Message)
	out := make(chan modules.Message)


	srv := NewTCPServ(":7275")


	srv.AddChannels(modules.Channels{In:in, Out:out})

	if err:= srv.Open(); err != nil {
		t.Fatalf("Cannot open server %s", err)
	}

	// Read master.
	for {
		select {
		case x := <-out:
			t.Log(x)
		}
	}
}