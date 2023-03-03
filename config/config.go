package config

import "flag"

type Config struct {
	Servers      int
	Clients      int
	Binary       string
	Directory    string
	Cidr         string
	BindServer   string
	Prefix       string
	ServerPrefix string
	ClientPrefix string
	Export       bool
	Import       bool
	Persist      bool
	Plan         bool
	Clean        bool
	Ips          []string
}

func MakeConfig() (cfg Config) {
	flag.IntVar(&cfg.Servers, "servers", 3, "Number of Servers")
	flag.IntVar(&cfg.Clients, "clients", 6, "Number of Clients")
	flag.StringVar(&cfg.Binary, "binary", "/usr/bin/nomad", "Location of Nomad Binary")
	flag.StringVar(&cfg.Directory, "directory", "/tmp/nomad-sim", "Working Directory")
	flag.StringVar(&cfg.Cidr, "cidr", "10.10.10.0/24", "CIDR Block for IP Assignment")
	flag.StringVar(&cfg.BindServer, "bind-server", "", "Network device or IP to bind the first server to")
	flag.StringVar(&cfg.Prefix, "prefix", "nmd", "Prefix of Nomad Cluster Members")
	flag.StringVar(&cfg.ServerPrefix, "server-prefix", "s", "Prefix of Nomad Servers")
	flag.StringVar(&cfg.ClientPrefix, "client-prefix", "c", "Prefix of Nomad Clients")
	flag.BoolVar(&cfg.Export, "export", false, "Export Nomad Node Layout")
	flag.BoolVar(&cfg.Import, "import", false, "Import Nomad Node Layout")
	flag.BoolVar(&cfg.Persist, "persist", false, "Persist resources after run")
	flag.BoolVar(&cfg.Plan, "plan", false, "Plan mode stages but does not run")
	flag.BoolVar(&cfg.Clean, "clean", false, "Clean mode to fix up any residual resources")
	flag.Parse()
	return cfg
}
