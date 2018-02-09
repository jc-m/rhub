package serial

import (
	"testing"
)

func TestPTY(t *testing.T) {
	c := make(map[string]string)

	pty, err := NewPty(c)
	if err != nil {
		t.Fatal(err)
	}
	q := pty.GetQueues()

	if err:= pty.Open(); err != nil {
		t.Fatalf("Cannot open port %s", err)
	}

	// Read master.
	x := <- q.Read
	t.Log( x)

}

