package utils

import (
	"testing"
	"github.com/jc-m/rhub/modules"
	"bytes"
)

func TestBuffer(t *testing.T) {
	downstreamRD := make(chan modules.Message,1)
	downstreamWR := make(chan modules.Message,1)

	ctl := make(chan bool)

	conf := BufferConfig{
		Delimiter: ';',
	}
	c := NewCmdBuffer(conf)
	q := c.GetQueues()
	c.ConnectQueuePair(&modules.QueuePair{Read:downstreamRD, Write:downstreamWR, Ctl:ctl})


	c.Open()
	q.Read <- modules.Message{Body:[]byte("AF;G")}
	cmd := <- downstreamWR
	if bytes.Compare(cmd.Body, []byte("AF;")) != 0 {
		t.Fatalf("Unexpected result")
	} else {
		t.Logf("received %s", cmd.Body)
	}
	q.Read <- modules.Message{Body:[]byte("H;FE;")}
	cmd = <- downstreamWR
	if bytes.Compare(cmd.Body, []byte("GH;")) != 0 {
		t.Fatalf("Unexpected result")
	} else {
		t.Logf("received %s", cmd.Body)
	}
	cmd = <- downstreamWR
	if bytes.Compare(cmd.Body, []byte("FE;")) != 0 {
		t.Fatalf("Unexpected result")
	} else {
		t.Logf("received %s", cmd.Body)
	}

	q.Ctl <- true
}