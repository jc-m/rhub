package radio

import (
	"github.com/jc-m/rhub/modules"
	"fmt"
)

type model struct {
	config Config
	channels []modules.Channels
	name string
}

type Config struct {
	Delimiter string
}

func (r *model) GetType() int {
	return modules.MUX
}

func (r *model) GetName() string {
	return r.name
}

func (r *model) AddChannels(channels modules.Channels) ([]modules.Channels, error)  {
	if len(r.channels) > 2 {
		return nil, fmt.Errorf("Module supports only 2 pairs of channels")
	}
	r.channels = append(r.channels, channels)
	return r.channels, nil
}

func (r *model) Open()  error {
}

func (r *model) Close() {
}

func NewPty() modules.Module {

	return &model {
		channels: make([]modules.Channels,0),
	}
}
