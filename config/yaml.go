package config

type Module struct {
	Name string
	Module string
	Upstream string
	Config map[string]string
}

type Stream struct {
	Name string
	Modules []Module
}

type Config struct {
	Streams []Stream
}
