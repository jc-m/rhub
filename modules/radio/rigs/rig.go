package rigs


type Rig interface {
	Open() error
	OnCat(cmd string, dir int) (*RigCommand, error)
	OnRig(cmd *RigCommand, dir int) (string, error)

}

type RigCommand struct {
	Id string `json:"type"`
	Params  map[string]string `json:"params"`
}
