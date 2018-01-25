package ft991a

import (
	"testing"
	"fmt"
)

func TestAIRsp(t *testing.T) {

	x, err := IFRsp("IF001007070000+1000C00000000;")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", x)
}

func TestMDRsp(t *testing.T) {
	x, err := MDRsp("0D")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", x)

}