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
		c["port"] = "/dev/ttys005"
		c["baud"] = "38400"
		c["data_bits"] = "8"
	    c["stop_bits"] = "1"

		client, err := NewSerial(c)
		if err != nil {
			t.Fatal(err)
		}
		channels:= client.GetQueues()

		if err := client.Open(); err != nil {
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
				os.Exit(0)
			}
		}



}