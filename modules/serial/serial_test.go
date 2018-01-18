package serial

import (
	"testing"
	"bytes"
	"fmt"
	"os"
	"github.com/jc-m/rhub/modules"
)

func TestSerial(t *testing.T) {
		c := SerialConfig{
			Port: "/dev/ttys005",
			Baud: 38400,
			DataBits:8,
			StopBits:2,
		}
		channels := modules.Channels{
			In : make(chan modules.Message),
			Out : make(chan modules.Message),
			Ctl :  make(chan bool),
		}


		client := NewSerial(c)
		client.AddChannels(channels)

		err := client.Open()
		if err != nil {
			t.Fatal(err)
		}
		var buffer bytes.Buffer
		buffer.Write([]byte("Test"))
		channels.In <- modules.Message{ Id:"test", Body:buffer.Bytes()}

		for {
			select {
			case r := <-channels.Out:
				fmt.Printf("XX %s\n", string(r.Body))
			case <- channels.Ctl: // TODO this gets executed twice if the port is closed
				fmt.Print("Got an exit")
				client.Close()
				break
			}
		}

	    os.Exit(0)

}