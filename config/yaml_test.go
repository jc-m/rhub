package config

import (
	"log"
	"testing"
	"github.com/go-yaml/yaml"
)
var sampleConf = `
streams:
  - name : ser2net
    modules :
    - name     : pipe
      module   : pipe
    - name     : serial
      module   : serial
      upstream : pipe
      config:
         port  : 7375
         bind  : 0.0.0.0/0
`
func TestYaml(t *testing.T) {

	v := Config{}
	err := yaml.Unmarshal([]byte(sampleConf), &v)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("%+v", v)
	if len(v.Streams) != 1 {
		t.Errorf("Unexpected value")
	}
	if v.Streams[0].Name != "ser2net" {
		t.Errorf("Unexpected value")
	}
	if len(v.Streams[0].Modules) != 2 {
		t.Errorf("Unexpected value")
	}
	if v.Streams[0].Modules[1].Config["port"] != "7375" {
		t.Errorf("Unexpected value")
	}
}