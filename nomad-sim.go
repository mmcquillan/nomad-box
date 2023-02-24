package main

import (
	"fmt"
	"os"

	"github.com/mmcquillan/nomad-sim/checks"
	"github.com/mmcquillan/nomad-sim/config"
	"github.com/mmcquillan/nomad-sim/node"
)

func main() {

	// configurable variables
	cfg := config.MakeConfig()
	checks.Checks(&cfg)

	// make nodes
	nodes := node.MakeNodes(cfg)

	// clean
	if cfg.Clean {
		node.CleanNodeResources(cfg, nodes)
		os.Exit(0)
	}

	// plan exit
	if cfg.Plan {
		fmt.Println("Plan Mode (quitting)")
		os.Exit(0)
	}

	// start up nodes
	node.BuildNodes(cfg, nodes)

	// wait to quit
	fmt.Print("Cluster Running (enter to quit)")
	fmt.Scanln()

	// clean up
	node.CleanNodes(cfg, nodes)

}
