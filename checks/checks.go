package checks

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"

	"github.com/mmcquillan/nomad-sim/config"
	"github.com/mmcquillan/nomad-sim/network"
)

func Checks(cfg *config.Config) {

	fmt.Println("Sim Checks...")

	// check we are on a supported OS
	fmt.Println(" - Checking OS")
	if runtime.GOOS != "linux" {
		fmt.Println("   ERROR: nomad-sim is made for the linux platform")
		if !cfg.Plan {
			os.Exit(2)
		}
	}

	// check we are running as root
	fmt.Println(" - Checking User")
	user, err := user.Current()
	if err != nil {
		fmt.Println("   ERROR: Could not read the current user")
		if !cfg.Plan {
			os.Exit(2)
		}
	}
	if user.Username != "root" {
		fmt.Println("   ERROR: nomad-sim must be run with root privledge")
		if !cfg.Plan {
			os.Exit(2)
		}
	}

	// check we have required tools installed
	fmt.Println(" - Checking Installed Tools")
	_, err = exec.LookPath("ip")
	if err != nil {
		fmt.Println("   ERROR: ip is not installed")
		if !cfg.Plan {
			os.Exit(2)
		}
	}

	// check for even server number
	fmt.Println(" - Checking Server Count")
	if cfg.Servers%2 == 0 {
		fmt.Println("   WARN: Even number of Servers is weird")
	}

	// check binary exists
	fmt.Println(" - Checking Nomad Binary location")
	if _, err := os.Stat(cfg.Binary); err != nil {
		fmt.Println("   ERROR: Nomad Binary does not exist")
		if !cfg.Plan {
			os.Exit(2)
		}
	}

	// check cidr
	fmt.Println(" - Checking Cidr Formatting")
	cfg.Ips, err = network.CidrToIps(cfg.Cidr)
	if err != nil {
		fmt.Println("   ERROR: Cidr formatting")
		fmt.Print(err)
		if !cfg.Plan {
			os.Exit(2)
		}
	}

	// check cidr vs server count
	fmt.Println(" - Checking Cidr / Server count")
	if len(cfg.Ips) < cfg.Servers+cfg.Clients {
		fmt.Println("   ERROR: Cidr does not allow enough IP's")
		if !cfg.Plan {
			os.Exit(2)
		}
	}

}
