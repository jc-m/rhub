package network

import (
	"testing"
)

func TestTcpServer(t *testing.T) {

	srv := NewTCPServ(":7275")


	q, _ := srv.CreateQueue()

	if err:= srv.Open(); err != nil {
		t.Fatalf("Cannot open server %s", err)
	}

	// Read master.
	for {
		select {
		case x := <-q.Write:
			t.Log(x)
		}
	}
}