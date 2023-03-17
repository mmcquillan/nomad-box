package checks

import (
	"os"
	"os/exec"
	"os/user"
	"runtime"

	"github.com/mmcquillan/nomad-box/config"
	"github.com/mmcquillan/nomad-box/network"
	"github.com/mmcquillan/nomad-box/run"
)

func Checks(cfg *config.Config) {

	run.Header("Nomad Box Pre Checks")

	// check we are on a supported OS
	run.Out("Checking OS")
	if runtime.GOOS != "linux" {
		run.Error("nomad-box is made for the linux platform")
		if !cfg.Plan {
			os.Exit(2)
		}
	}

	// check we are running as root
	run.Out("Checking Users")
	user, err := user.Current()
	if err != nil {
		run.Error("Could not read the current user")
		if !cfg.Plan {
			os.Exit(2)
		}
	}
	if user.Username != "root" {
		run.Error("nomad-box must be run with root privledge")
		if !cfg.Plan {
			os.Exit(2)
		}
	}

	// check we have required tools installed
	run.Out("Checking Installed Tools")
	_, err = exec.LookPath("ip")
	if err != nil {
		run.Error("ip is not installed")
		if !cfg.Plan {
			os.Exit(2)
		}
	}

	// check for even server number
	run.Out("Checking Server Count")
	if cfg.Servers%2 == 0 {
		run.Warn("Even number of Servers is weird")
	}

	// check binary exists
	run.Out("Checking Nomad Binary location")
	if _, err := os.Stat(cfg.Binary); err != nil {
		run.Error("Nomad Binary does not exist")
		if !cfg.Plan {
			os.Exit(2)
		}
	}

	// check server config
	if cfg.ServerConfig != "" {
		run.Out("Checking Server Config")
		if _, err := os.Stat(cfg.ServerConfig); err != nil {
			run.Error("Server Config does not exist")
			if !cfg.Plan {
				os.Exit(2)
			}
		}
	}

	// check client config
	if cfg.ClientConfig != "" {
		run.Out("Checking Client Config")
		if _, err := os.Stat(cfg.ClientConfig); err != nil {
			run.Error("Client Config does not exist")
			if !cfg.Plan {
				os.Exit(2)
			}
		}
	}

	// check cidr
	run.Out("Checking Cidr Formatting")
	cfg.Ips, err = network.CidrToIps(cfg.Cidr)
	if err != nil {
		run.Error("Cidr formatting")
		run.Error(err.Error())
		if !cfg.Plan {
			os.Exit(2)
		}
	}

	// check cidr vs server count
	run.Out("Checking Cidr / Server Count")
	if len(cfg.Ips) < cfg.Servers+cfg.Clients {
		run.Error("Cidr does not allow enough IP's")
		if !cfg.Plan {
			os.Exit(2)
		}
	}

	// check server device
	if cfg.BindServer != "" {
		run.Out("Checking Bind Server")
		ip := network.GetIpFromDevice(cfg.BindServer)
		if ip == "" {
			run.Error("Could not find device for Bind Server")
			if !cfg.Plan {
				os.Exit(2)
			}
		} else {
			cfg.Ips[0] = ip
		}
	}

}
