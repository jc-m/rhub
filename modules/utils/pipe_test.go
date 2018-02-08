package utils

import (
	"testing"
	"github.com/jc-m/rhub/modules"
)

func TestPipe(t *testing.T) {
	leftRD := make(chan modules.Message,1)
	leftWR := make(chan modules.Message,1)
	leftCtl := make(chan bool)

	rightRD := make(chan modules.Message,1)
	rightWR := make(chan modules.Message,1)
	rightCtl := make(chan bool)

	pipe, err := NewPipe(modules.ModuleConfig{"tap":"false"})
	if err != nil {
		t.Fatal(err)
	}
	q := pipe.GetQueues()

	pipe.ConnectQueuePair(&modules.QueuePair{Read:leftRD, Write:leftWR, Ctl:leftCtl})
	pipe.ConnectQueuePair(&modules.QueuePair{Read:rightRD, Write:rightWR, Ctl:rightCtl})

	pipe.Open()


	leftWR <- modules.Message{Body:[]byte("test")}

	m := <- rightRD
	t.Logf("%+v", m)
	q.Ctl <- true
}