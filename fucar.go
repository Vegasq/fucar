package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func run(command string) string {
	cmd := strings.Split(command, " ")
	bin := cmd[0]

	out, err := exec.Command(bin, cmd[1:]...).Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func get_contral_db_nodes() []*Node {
	// Contains our represent of nodes
	var nodes []*Node
	// Contains parsed JSON
	var fuelNodes []*FuelNodeJson

	// Parse JSON
	output := run("fuel nodes --json")

	non_parsed_nodes := []byte(output)

	err := json.Unmarshal(non_parsed_nodes, &fuelNodes)
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	// Build Node object
	for _, n := range fuelNodes {

		// Filter only contrail-db nodes
		if !n.is_contrail_db() {
			continue
		}

		// cs_ip, err := n.get_contrail_ip()
		// if err != nil {
		// 	log.Fatal("Can't find contrail IP.")
		// 	os.Exit(1)
		// }

		log.Println("Contrail DB node " + n.Hostname + " found.")

		node := Node{
			id:       n.ID,
			admin_ip: n.IP,
			hostname: n.Hostname}
		node.Init()

		nodes = append(nodes, &node)
	}
	return nodes
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {
	log.Println("Starting.")

	oldNodesParam := flag.String("old-nodes", "", "Lis of node ids separated by koma.")
	flag.Parse()
	oldNodes := strings.Split(*oldNodesParam, ",")

	nodes := get_contral_db_nodes()
	for _, node := range nodes {
		for _, od := range oldNodes {
			id_as_str, err := strconv.Atoi(od)
			if err != nil {
				panic("Can't parse ID: " + od)
			}

			if node.id == id_as_str {
				log.Println("Removing node ", node.hostname)
				node.ssh.Exec("sudo nodetool repair")
				result := node.ssh.Exec("sudo nodetool decommission")
				log.Println("Decommission output:", result)
			}
		}
	}

	for _, node := range nodes {
		if !stringInSlice(strconv.Itoa(node.id), oldNodes) {
			node.ssh.Exec("sudo nodetool repair")
		}
	}
}
