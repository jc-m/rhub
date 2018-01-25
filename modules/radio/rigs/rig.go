package rigs


type Rig interface {
	Open() error
	OnCatUpStream(cmd string) (*RigCommand, error)
}

type RigCommand struct {
	Id string `json:"type"`
	Params  map[string]string `json:"params"`
}
