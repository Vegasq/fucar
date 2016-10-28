package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
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
	var fuelNodes []FuelNodeJson

	// Parse JSON
	output := run("fuel nodes --json")
	non_parsed_nodes := []byte(output)
	json.Unmarshal(non_parsed_nodes, &fuelNodes)

	// Build Node object
	for _, n := range fuelNodes {

		// Filter only contrail-db nodes
		if !n.is_contrail_db() {
			continue
		}

		cs_ip, err := n.get_contrail_ip()
		if err != nil {
			log.Fatal("Can't find contrail IP.")
			os.Exit(1)
		}

		log.Println("Node " + n.hostname + " found.")

		node := Node{
			admin_ip:     n.ip,
			hostname:     n.hostname,
			cassandra_ip: cs_ip}

		nodes = append(nodes, &node)
	}
	return nodes
}

func main() {
	set_rf_to := flag.Int("replica-factor", 0, "Set replica factor to value.")
	flag.Parse()

	if *set_rf_to == 0 {
		log.Fatalln("Set --replica-factor X.")
	}

	nodes := get_contral_db_nodes()
	nodes[0].set_replication_factor(*set_rf_to)
	for _, node := range nodes {
		node.run("sudo nodetool repair")
	}

}
