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


	r, err := New(modules.ModuleConfig{})
	if err != nil {
		t.Fatal(err)

	}
	q := r.GetQueues()

	r.ConnectQueuePair(&modules.QueuePair{Read: upstreamWR, Write: upstreamRD, Ctl:ctl})

	r.Open()

	upstreamWR <- modules.Message{Body:[]byte("IF001007070000+1000C00000000;")}
	x:= <- q.Write

	enc := gob.NewDecoder(bytes.NewReader(x.Body))
	var v rigs.RigCommand
	if err := enc.Decode(&v); err != nil {
		t.Fatalf("Failed to decode %s", err)
	}
	if v, ok := v.Params["VFOA"]; ok {
		if v != "7070000" {
			t.Fatalf("Unexpected value %s", v)

		}
	} else {
		t.Fatalf("Missing value %s", "VFOA")
	}
}