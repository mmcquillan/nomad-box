package node

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/mmcquillan/nomad-box/config"
	"github.com/mmcquillan/nomad-box/network"
	"github.com/mmcquillan/nomad-box/run"
)

type Node struct {
	Server bool
	Binary string
	Name   string
	Region string
	Dc     string
	Pool   string
	Ip     string
	Device string
	Dir    string
	Config string
	Params string
	Pid    int `json:"-"`
}

func MakeNodes(cfg config.Config) (nodes []Node) {

	// node slice
	nodes = make([]Node, cfg.Servers+cfg.Clients)

	// import
	if cfg.Import {
		nodes = importNodes()
		if len(nodes) > 0 {
			run.Header("Importing Nodes")
			for i := 0; i < len(nodes); i++ {
				printNode(nodes[i])
			}
			return nodes
		}
	}

	// node marker
	marker := 0

	// start feedback
	run.Header("Mapping Nodes")

	// make servers
	for s := 0; s < cfg.Servers; s++ {
		nodes[marker].Server = true
		nodes[marker].Binary = cfg.Binary
		nodes[marker].Name = cfg.Prefix + cfg.ServerPrefix + strconv.Itoa(s)
		nodes[marker].Region = "global"
		nodes[marker].Dc = "dc1"
		nodes[marker].Pool = "default"
		nodes[marker].Ip = cfg.Ips[marker]
		nodes[marker].Device = cfg.Prefix + "eth" + strconv.Itoa(marker)
		if s == 0 && cfg.BindServer != "" {
			nodes[marker].Device = cfg.BindServer
		}
		nodes[marker].Dir = cfg.Directory + "/" + nodes[marker].Name
		if cfg.ServerConfig != "" {
			nodes[marker].Config = cfg.ServerConfig
		}
		if cfg.ServerParams != "" {
			nodes[marker].Params = cfg.ServerParams
		}
		nodes[marker].Pid = 0
		printNode(nodes[marker])
		marker++
	}

	// make clients
	for c := 0; c < cfg.Clients; c++ {
		nodes[marker].Server = false
		nodes[marker].Binary = cfg.Binary
		nodes[marker].Name = cfg.Prefix + cfg.ClientPrefix + strconv.Itoa(c)
		nodes[marker].Region = "global"
		nodes[marker].Dc = "dc1"
		nodes[marker].Pool = "default"
		nodes[marker].Ip = cfg.Ips[marker]
		nodes[marker].Device = cfg.Prefix + "eth" + strconv.Itoa(marker)
		nodes[marker].Pid = 0
		nodes[marker].Dir = cfg.Directory + "/" + nodes[marker].Name
		if cfg.ClientConfig != "" {
			nodes[marker].Config = cfg.ClientConfig
		}
		if cfg.ClientParams != "" {
			nodes[marker].Params = cfg.ClientParams
		}
		printNode(nodes[marker])
		marker++
	}

	// export
	if cfg.Export {
		exportNodes(nodes)
	}

	return nodes
}

func BuildNodes(cfg config.Config, nodes []Node) {

	run.Header("Building Nodes")

	// check the nodes
	for i := 0; i < len(nodes); i++ {

		printNode(nodes[i])

		// node networking and directory space
		makeNodeResources(cfg, nodes[i])

		// run nomad process
		if nodes[i].Server {

			// run server nomad process
			nomad := nodes[i].Binary + " agent "
			nomad += " -node=" + nodes[i].Name
			nomad += " -bind=" + nodes[i].Ip
			nomad += " -bootstrap-expect=" + strconv.Itoa(cfg.Servers)
			nomad += " -data-dir=" + nodes[i].Dir
			nomad += " -dc=" + nodes[i].Dc
			if nodes[i].Config != "" {
				nomad += " -config=" + nodes[i].Config
			}
			if cfg.UI {
				nomad += " -config=" + cfg.Directory + "/ui-config.hcl"
			}
			for j := 0; j < len(nodes); j++ {
				if nodes[j].Server {
					nomad += " -join=" + nodes[j].Ip
				}
			}
			if cfg.Log {
				nomad += " -log-level=" + cfg.LogLevel
			}
			nomad += " -network-interface=" + nodes[i].Device
			nomad += " -region=" + nodes[i].Region
			nomad += " -server"
			if nodes[i].Params != "" {
				nomad += " " + nodes[i].Params
			}
			nodes[i].Pid = run.Process(nomad, nodes[i].Name, cfg.Log)
			time.Sleep(3 * time.Second)

		} else {

			// run client nomad process
			nomad := nodes[i].Binary + " agent "
			nomad += " -node=" + nodes[i].Name
			nomad += " -bind=" + nodes[i].Ip
			nomad += " -client"
			nomad += " -data-dir=" + nodes[i].Dir
			nomad += " -dc=" + nodes[i].Dc
			nomad += " -node-pool=" + nodes[i].Pool
			if nodes[i].Config != "" {
				nomad += " -config=" + nodes[i].Config
			}
			if cfg.Log {
				nomad += " -log-level=" + cfg.LogLevel
			}
			for j := 0; j < len(nodes); j++ {
				if nodes[j].Server {
					nomad += " -servers=" + nodes[j].Ip + ":4647"
				}
			}
			nomad += " -network-interface=" + nodes[i].Device
			nomad += " -region=" + nodes[i].Region
			if nodes[i].Params != "" {
				nomad += " " + nodes[i].Params
			}
			nodes[i].Pid = run.Process(nomad, nodes[i].Name, cfg.Log)

		}

	}

	run.Out("export NOMAD_ADDR=\"http://" + nodes[1].Ip + ":4646")

}

func CleanNodes(cfg config.Config, nodes []Node) {
	run.Header("Cleaning Nodes")
	for i := 0; i < len(nodes); i++ {
		if !nodes[i].Server {
			printNode(nodes[i])
			cleanNodeProcess(cfg, nodes[i])
			if !cfg.Persist {
				cleanNodeResources(cfg, nodes[i])
			}
		}
	}
	for i := 0; i < len(nodes); i++ {
		if nodes[i].Server {
			printNode(nodes[i])
			cleanNodeProcess(cfg, nodes[i])
			if !cfg.Persist {
				cleanNodeResources(cfg, nodes[i])
			}
		}
	}
}

func CleanNodeResources(cfg config.Config, nodes []Node) {
	run.Header("Cleaning Node Resources")
	for i := 0; i < len(nodes); i++ {
		printNode(nodes[i])
		cleanNodeResources(cfg, nodes[i])
	}
}

func makeNodeResources(cfg config.Config, node Node) {

	// network check if exists
	if run.CommandContains("ip a", node.Ip) {
		if !cfg.Persist {
			cleanNodeResources(cfg, node)
			makeNodeResourcesNetwork(cfg, node)
		}
	} else {
		makeNodeResourcesNetwork(cfg, node)
	}

	// make server directory
	run.Command("mkdir -p " + node.Dir)

	// write ui config
	if cfg.UI {
		config := []byte(`ui {
  enabled = true
  label {
    text             = "NMD-BOX"
    background_color = "red"
    text_color       = "black"
  }
}
`)
		err := os.WriteFile(cfg.Directory+"/ui-config.hcl", config, 0644)
		if err != nil {
			run.Error("Cannot Write UI Config")
			run.Error(err.Error())
		}
	}

}

func makeNodeResourcesNetwork(cfg config.Config, node Node) {

	if cfg.BindServer != node.Device {

		// setup network device
		run.Command("ip link add " + node.Device + " type dummy")

		// set mac address
		run.Command("ip link set dev " + node.Device + " address " + network.GenerateMac())

		// set IP address
		run.Command("ip addr add " + node.Ip + "/24 brd + dev " + node.Device + " label " + node.Device + ":0")

		// bring up device
		run.Command("ip link set dev " + node.Device + " up")

	}

}

func cleanNodeResources(cfg config.Config, node Node) {

	if cfg.BindServer != node.Device {

		// delete address from device
		run.Command("ip addr del " + node.Ip + "/24 brd + dev " + node.Device + " label " + node.Device + ":0")

		// delete network device
		run.Command("ip link delete " + node.Device + " type dummy")

	}

	// delete server directory
	time.Sleep(3 * time.Second)
	run.Command("rm -rf " + node.Dir)

}

func cleanNodeProcess(cfg config.Config, node Node) {

	// kill process
	run.Command("kill -2 " + strconv.Itoa(node.Pid))
	for run.CheckProcess(node.Pid) {
		time.Sleep(3 * time.Second)
	}

}

func printNode(node Node) {
	n := fmt.Sprintf("%s.%s.%s [ %s : %s : %s ]", node.Region, node.Dc, node.Name, node.Ip, node.Device, node.Dir)
	run.Out(n)
}

func importNodes() (nodes []Node) {
	mydir, _ := os.Getwd()
	file, err := os.ReadFile(mydir + "/nodes.json")
	if err != nil {
		return nodes
	}
	err = json.Unmarshal([]byte(file), &nodes)
	if err != nil {
		run.Error("Cannot Import Nodes")
		run.Error(err.Error())
	}
	return nodes
}

func exportNodes(nodes []Node) {
	mydir, _ := os.Getwd()
	imp_json, err := json.MarshalIndent(nodes, "", "   ")
	if err != nil {
		run.Error("Cannot Export Nodes")
		run.Error(err.Error())
	}
	err = os.WriteFile(mydir+"/nodes.json", imp_json, 0644)
	if err != nil {
		run.Error("Cannot Export Nodes")
		run.Error(err.Error())
	}
}
