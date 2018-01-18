package server

import (
	"github.com/jc-m/rhub/modules"
)


func Push(lower, upper modules.Module) error {

	in := make(chan modules.Message)
	out := make(chan modules.Message)
	ctl := make(chan bool)

	if _, err := lower.AddChannels(modules.Channels{In:in, Out:out, Ctl:ctl}); err != nil {
		return err
	}
	if _, err := upper.AddChannels(modules.Channels{In:out, Out:in, Ctl:ctl}); err != nil {
		return err
	}

	if err := lower.Open(); err != nil {
		return err
	}
	if err := upper.Open(); err != nil {
		return err
	}

	return nil
}