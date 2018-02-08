package stream

import (
	"testing"
	"github.com/jc-m/rhub/config"
)

func TestNew(t *testing.T) {


	mods := []config.Module{
		{"serial", "serial", "pipe", map[string]string{"port":"/dev/ttys012"} },
		{"pty", "pty", "pipe", map[string]string{} },
		{"pipe", "pipe", "", map[string]string{"tap":"false"} },
	}

	conf := config.Stream{
		Name: "test",
		Modules: mods,
	}
	s := NewStream(conf)

	if len(s.NodeMap) != len(mods) {
		t.Fatalf("Missing modules")
	}

	for k,v := range s.Index {
		t.Logf("%+v", v)

		if k == "pty" || k == "serial" {
			if v.Upstream[0] != "pipe" {
				t.Logf("Invalid upstream %+v", v)
			}
		}
		if k == "pipe" {
			if len(v.Upstream) >0 {
				t.Logf("Invalid upstream %+v", v)
			}
		}
	}
}
