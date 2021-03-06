package serial

import (
	"testing"
	"bytes"
	"fmt"
	"os"
	"github.com/jc-m/rhub/modules"
)

func TestSerial(t *testing.T) {

	c := make(map[string]string)
	c["port"] = "/dev/tty.SLAB_USBtoUART"
	c["baud"] = "38400"
	c["stop_bits"] = "2"
	c["rts_cts"] = "true"

	client, err := NewSerial(c)
	if err != nil {
		t.Fatal(err)
	}
	channels:= client.GetQueues()

	if err := client.Open(); err != nil {
		t.Fatal(err)
	}

	var buffer bytes.Buffer
	buffer.Write([]byte("PS1;"))
	channels.Read <- modules.Message{ Id:"test", Body:buffer.Bytes()}

	for {
		select {
		case r := <-channels.Write:
			fmt.Printf("XX %s\n", string(r.Body))
		case <- channels.Ctl: // TODO this gets executed twice if the port is closed
			fmt.Print("Got an exit")
			os.Exit(0)
		}
	}
}

func TestOpen(t *testing.T) {
	c := SerialConfig{
		Port: "/dev/tty.Bluetooth-Incoming-Port",
		Baud: 9600,
		StopBits: 2,
	}
	if _, err := openInternal(c); err != nil {
		t.Fatal(err)
	}
}