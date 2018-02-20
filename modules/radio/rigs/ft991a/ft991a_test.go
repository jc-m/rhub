package ft991a

import (
	"testing"
	"fmt"
	"github.com/jc-m/rhub/modules/radio/rigs"
)

func TestAIRsp(t *testing.T) {

	x, err := IF("IF001007070000+1000C00000000;", rigs.CAT_DIR_DOWN)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", x)
}

func TestMDRsp(t *testing.T) {
	x, err := MD("0D", rigs.CAT_DIR_DOWN)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", x)

}

func TestRMRsp(t *testing.T) {
	x, err := RM("6023", rigs.CAT_DIR_DOWN)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", x)

}

func TestPCRsp(t *testing.T) {
	x, err := PC("005", rigs.CAT_DIR_DOWN)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", x)

}

func TestNLRsp(t *testing.T) {
	x, err := NL("0010", rigs.CAT_DIR_DOWN)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", x)

}

func TestRLRsp(t *testing.T) {
	x, err := RL("006", rigs.CAT_DIR_DOWN)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", x)

}