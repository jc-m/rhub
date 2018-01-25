package radio

import (
	"testing"
	"github.com/jc-m/rhub/modules"
	"encoding/gob"
	"bytes"
	"github.com/jc-m/rhub/modules/radio/rigs"
)

func TestRadio(t *testing.T) {

	upstreamRD := make(chan modules.Message,1)
	upstreamWR := make(chan modules.Message,1)
	ctl := make(chan bool)


	r := New()
	q, _ := r.CreateQueue()

	r.ConnectQueuePair(&modules.QueuePair{Read: upstreamWR, Write: upstreamRD, Ctl:ctl})

	r.Open()

	upstreamWR <- modules.Message{Body:[]byte("IF001007070000+1000C00000000;")}
	x:= <- q.Write

	enc := gob.NewDecoder(bytes.NewReader(x.Body))
	var v rigs.RigCommand
	if err := enc.Decode(&v); err != nil {
		t.Fatalf("Failed to decode %s", err)
	} else {
		t.Logf("%+v", v)
	}

}