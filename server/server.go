package server

import (
	"context"
	"log"
	"os"
	"syscall"
	"os/signal"
	"github.com/jc-m/rhub/config"
	"github.com/jc-m/rhub/stream"
)

type Server struct {
	params *Params
}

type Params struct {
	ConfigPath string
}

func NewParams() *Params  {
	return &Params{}
}

func New(p *Params) *Server {
	return &Server{
		params: p,
	}
}

func (s *Server) Run(ctx context.Context) error {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	conf := config.GetConfig(s.params.ConfigPath)

	log.Printf("%+v", conf)
	for _, streamConfig := range conf.Streams {
		if st := stream.NewStream(streamConfig); st != nil {
			st.Start()
		}
	}

	for {
		select {
		case sig := <-sigChan:
			log.Printf("[INFO] Caught signal %s; shutting down", sig)
			os.Exit(0)
		}
	}
	return nil
}