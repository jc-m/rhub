package network

import (
	"testing"
	"github.com/jc-m/rhub/modules"
)

func TestTcpServer(t *testing.T) {


	srv, err := NewTCPServ(modules.ModuleConfig{"address":":7375"})
	if err != nil {
		t.Fatal(err)
	}
	if err := srv.Open(); err != nil {
		t.Fatalf("Cannot open server %s", err)
	}

	q := srv.GetQueues()

	// Read master.
	for {
		select {
		case x := <-q.Write:
			t.Log(x)
		}
	}
}