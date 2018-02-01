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

func TestRMRsp(t *testing.T) {
	x, err := RMRsp("6023")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", x)

}

func TestPCRsp(t *testing.T) {
	x, err := PCRsp("005")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", x)

}

func TestNLRsp(t *testing.T) {
	x, err := NLRsp("0010")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", x)

}

func TestRLRsp(t *testing.T) {
	x, err := RLRsp("006")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", x)

}