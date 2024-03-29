package config

import "flag"

type Config struct {
	Servers      int
	Clients      int
	Binary       string
	Directory    string
	Cidr         string
	BindServer   string
	Log          bool
	LogLevel     string
	Prefix       string
	ServerPrefix string
	ClientPrefix string
	ServerConfig string
	ClientConfig string
	ServerParams string
	ClientParams string
	Export       bool
	Import       bool
	Persist      bool
	Plan         bool
	Clean        bool
	UI           bool
	Ips          []string
}

func MakeConfig() (cfg Config) {
	flag.IntVar(&cfg.Servers, "servers", 3, "Number of Servers")
	flag.IntVar(&cfg.Clients, "clients", 6, "Number of Clients")
	flag.StringVar(&cfg.Binary, "binary", "/usr/bin/nomad", "Location of Nomad Binary")
	flag.StringVar(&cfg.Directory, "directory", "/tmp/nomad-box", "Working Directory")
	flag.StringVar(&cfg.Cidr, "cidr", "10.10.10.0/24", "CIDR Block for IP Assignment")
	flag.StringVar(&cfg.BindServer, "bind-server", "", "Network device or IP to bind the first server to")
	flag.BoolVar(&cfg.Log, "log", false, "Show Nomad Logs in the console")
	flag.StringVar(&cfg.LogLevel, "log-level", "INFO", "Prefix of Nomad Cluster Members")
	flag.StringVar(&cfg.Prefix, "prefix", "nmd", "Prefix of Nomad Cluster Members")
	flag.StringVar(&cfg.ServerPrefix, "server-prefix", "s", "Prefix of Nomad Servers")
	flag.StringVar(&cfg.ClientPrefix, "client-prefix", "c", "Prefix of Nomad Clients")
	flag.StringVar(&cfg.ServerConfig, "server-config", "", "Path to a Server Config")
	flag.StringVar(&cfg.ClientConfig, "client-config", "", "Path to a Client Config")
	flag.StringVar(&cfg.ServerParams, "server-params", "", "Path to a Server Params")
	flag.StringVar(&cfg.ClientParams, "client-params", "", "Path to a Client Params")
	flag.BoolVar(&cfg.Export, "export", false, "Export Nomad Node Layout")
	flag.BoolVar(&cfg.Import, "import", false, "Import Nomad Node Layout")
	flag.BoolVar(&cfg.Persist, "persist", false, "Persist resources after run")
	flag.BoolVar(&cfg.Plan, "plan", false, "Plan mode stages but does not run")
	flag.BoolVar(&cfg.Clean, "clean", false, "Clean mode to fix up any residual resources")
	flag.BoolVar(&cfg.UI, "ui", false, "Adds a UI label")
	flag.Parse()
	return cfg
}
