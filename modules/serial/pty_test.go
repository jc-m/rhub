package serial

import (
	"testing"
)

func TestPTY(t *testing.T) {

	pty := NewPty()
	q, err := pty.CreateQueue()
	if err != nil {
		t.Fatalf("Unexpected number of queue")
	}

	if err:= pty.Open(); err != nil {
		t.Fatalf("Cannot open port %s", err)
	}

	// Read master.
	x := <- q.Read
	t.Log( x)

}

