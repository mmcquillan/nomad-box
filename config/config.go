package config

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"

	"github.com/mmcquillan/nomad-box/run"
)

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
	Export       bool `json:"-"`
	Import       bool `json:"-"`
	Persist      bool
	Plan         bool
	Clean        bool
	UI           bool
	Ips          []string `json:"-"`
}

func MakeConfig() (cfg Config) {

	// defaults
	cfg.Servers = 3
	cfg.Clients = 6
	cfg.Binary = "/usr/bin/nomad"
	cfg.Directory = "/tmp/nomad-box"
	cfg.Cidr = "10.10.10.0/24"
	cfg.BindServer = ""
	cfg.Log = false
	cfg.LogLevel = "INFO"
	cfg.Prefix = "nmd"
	cfg.ServerPrefix = "s"
	cfg.ClientPrefix = "c"
	cfg.ServerConfig = ""
	cfg.ClientConfig = ""
	cfg.ServerParams = ""
	cfg.ClientParams = ""
	cfg.Export = false
	cfg.Import = false
	cfg.Persist = false
	cfg.Plan = false
	cfg.Clean = false
	cfg.UI = false

	// env vars
	if val, err := strconv.Atoi(os.Getenv("NOMAD_BOX_SERVERS")); err == nil {
		cfg.Servers = val
	}
	if val, err := strconv.Atoi(os.Getenv("NOMAD_BOX_CLIENTS")); err == nil {
		cfg.Clients = val
	}
	if val := os.Getenv("NOMAD_BOX_BINARY"); val != "" {
		cfg.Binary = val
	}
	if val := os.Getenv("NOMAD_BOX_DIRECTORY"); val != "" {
		cfg.Directory = val
	}
	if val := os.Getenv("NOMAD_BOX_CIDR"); val != "" {
		cfg.Cidr = val
	}
	if val := os.Getenv("NOMAD_BOX_BIND_SERVER"); val != "" {
		cfg.BindServer = val
	}
	if val, err := strconv.ParseBool(os.Getenv("NOMAD_BOX_LOG")); err == nil {
		cfg.Log = val
	}
	if val := os.Getenv("NOMAD_BOX_LOG_LEVEL"); val != "" {
		cfg.LogLevel = val
	}
	if val := os.Getenv("NOMAD_BOX_PREFIX"); val != "" {
		cfg.Prefix = val
	}
	if val := os.Getenv("NOMAD_BOX_SERVER_PREFIX"); val != "" {
		cfg.ServerPrefix = val
	}
	if val := os.Getenv("NOMAD_BOX_CLIENT_PREFIX"); val != "" {
		cfg.ClientPrefix = val
	}
	if val := os.Getenv("NOMAD_BOX_SERVER_CONFIG"); val != "" {
		cfg.ServerConfig = val
	}
	if val := os.Getenv("NOMAD_BOX_CLIENT_CONFIG"); val != "" {
		cfg.ClientConfig = val
	}
	if val := os.Getenv("NOMAD_BOX_SERVER_PARAMS"); val != "" {
		cfg.ServerParams = val
	}
	if val := os.Getenv("NOMAD_BOX_CLIENT_PARAMS"); val != "" {
		cfg.ClientParams = val
	}
	if val, err := strconv.ParseBool(os.Getenv("NOMAD_BOX_EXPORT")); err == nil {
		cfg.Export = val
	}
	if val, err := strconv.ParseBool(os.Getenv("NOMAD_BOX_IMPORT")); err == nil {
		cfg.Import = val
	}
	if val, err := strconv.ParseBool(os.Getenv("NOMAD_BOX_PERSIST")); err == nil {
		cfg.Persist = val
	}
	if val, err := strconv.ParseBool(os.Getenv("NOMAD_BOX_PLAN")); err == nil {
		cfg.Plan = val
	}
	if val, err := strconv.ParseBool(os.Getenv("NOMAD_BOX_CLEAN")); err == nil {
		cfg.Clean = val
	}
	if val, err := strconv.ParseBool(os.Getenv("NOMAD_BOX_UI")); err == nil {
		cfg.UI = val
	}

	// flags
	flag.IntVar(&cfg.Servers, "servers", cfg.Servers, "Number of Servers")
	flag.IntVar(&cfg.Clients, "clients", cfg.Clients, "Number of Clients")
	flag.StringVar(&cfg.Binary, "binary", cfg.Binary, "Location of Nomad Binary")
	flag.StringVar(&cfg.Directory, "directory", cfg.Directory, "Working Directory")
	flag.StringVar(&cfg.Cidr, "cidr", cfg.Cidr, "CIDR Block for IP Assignment")
	flag.StringVar(&cfg.BindServer, "bind-server", cfg.BindServer, "Network device or IP to bind the first server to")
	flag.BoolVar(&cfg.Log, "log", cfg.Log, "Show Nomad Logs in the console")
	flag.StringVar(&cfg.LogLevel, "log-level", cfg.LogLevel, "Prefix of Nomad Cluster Members")
	flag.StringVar(&cfg.Prefix, "prefix", cfg.Prefix, "Prefix of Nomad Cluster Members")
	flag.StringVar(&cfg.ServerPrefix, "server-prefix", cfg.ServerPrefix, "Prefix of Nomad Servers")
	flag.StringVar(&cfg.ClientPrefix, "client-prefix", cfg.ClientPrefix, "Prefix of Nomad Clients")
	flag.StringVar(&cfg.ServerConfig, "server-config", cfg.ServerConfig, "Path to a Server Config")
	flag.StringVar(&cfg.ClientConfig, "client-config", cfg.ClientConfig, "Path to a Client Config")
	flag.StringVar(&cfg.ServerParams, "server-params", cfg.ServerParams, "Path to a Server Params")
	flag.StringVar(&cfg.ClientParams, "client-params", cfg.ClientParams, "Path to a Client Params")
	flag.BoolVar(&cfg.Export, "export", cfg.Export, "Export Nomad Node Layout")
	flag.BoolVar(&cfg.Import, "import", cfg.Import, "Import Nomad Node Layout")
	flag.BoolVar(&cfg.Persist, "persist", cfg.Persist, "Persist resources after run")
	flag.BoolVar(&cfg.Plan, "plan", cfg.Plan, "Plan mode stages but does not run")
	flag.BoolVar(&cfg.Clean, "clean", cfg.Clean, "Clean mode to fix up any residual resources")
	flag.BoolVar(&cfg.UI, "ui", cfg.UI, "Adds a UI label")
	flag.Parse()

	// import config
	if cfg.Import {
		cfg = importConfig()
	}

	// export config
	if cfg.Export {
		exportConfig(cfg)
	}

	// return the config
	return cfg

}

func importConfig() (cfg Config) {
	mydir, _ := os.Getwd()
	file, err := os.ReadFile(mydir + "/config.json")
	if err != nil {
		run.Error("Cannot Import Config")
		run.Error(err.Error())
	}
	err = json.Unmarshal([]byte(file), &cfg)
	if err != nil {
		run.Error("Cannot Import Config")
		run.Error(err.Error())
	}
	return cfg
}

func exportConfig(cfg Config) {
	mydir, _ := os.Getwd()
	cfg_json, err := json.MarshalIndent(cfg, "", "   ")
	if err != nil {
		run.Error("Cannot Export Config")
		run.Error(err.Error())
	}
	err = os.WriteFile(mydir+"/config.json", cfg_json, 0644)
	if err != nil {
		run.Error("Cannot Export Config")
		run.Error(err.Error())
	}
}
