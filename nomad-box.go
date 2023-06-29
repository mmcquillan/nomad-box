package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/mmcquillan/nomad-box/checks"
	"github.com/mmcquillan/nomad-box/config"
	"github.com/mmcquillan/nomad-box/node"
	"github.com/mmcquillan/nomad-box/run"
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
		run.Out("Plan Mode (quitting)")
		os.Exit(0)
	}

	// start up nodes
	node.BuildNodes(cfg, nodes)

	// setup to catch sigint
	q := make(chan os.Signal, 1)
	signal.Notify(q, os.Interrupt)
	go func() {
		<-q
		node.CleanNodes(cfg, nodes)
	}()

	// wait to quit
	run.Out("Cluster Running (enter to quit)")
	fmt.Scanln()

	// clean up
	node.CleanNodes(cfg, nodes)

}
