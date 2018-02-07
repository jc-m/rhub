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



		client := NewSerial(c)
		channels:= client.GetQueues()

		err := client.Open()
		if err != nil {
			t.Fatal(err)
		}
		var buffer bytes.Buffer
		buffer.Write([]byte("Test"))
		channels.Read <- modules.Message{ Id:"test", Body:buffer.Bytes()}

		for {
			select {
			case r := <-channels.Write:
				fmt.Printf("XX %s\n", string(r.Body))
			case <- channels.Ctl: // TODO this gets executed twice if the port is closed
				fmt.Print("Got an exit")
				client.Close()
				break
			}
		}

	    os.Exit(0)

}