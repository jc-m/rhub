package fl

import (
	"testing"
)

func TestServer(t *testing.T) {
	s := RPCServer{
		Address: ":1234",
	}

	s.Start()
}