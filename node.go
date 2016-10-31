package main

import (
	"errors"
	"strings"
)

type Node struct {
	id           int
	hostname     string
	admin_ip     string
	cassandra_ip string
	ssh          SSHNode
}

func (n *Node) Init() {
	n.ssh = SSHNode{username: "fuel", host: n.admin_ip}
}

func (n *Node) csql(command string) {
	remote_exec := "cqlsh " + n.cassandra_ip + " -e \"" + command + "\""
	n.ssh.Exec(remote_exec)
}

// func (n *Node) set_replication_factor(factor int) {
// 	rf := strconv.Itoa(factor)
// 	log.Println("Set replica_factor to " + rf)
// 	n.csql("ALTER KEYSPACE \"to_bgp_keyspace\" WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : " + rf + " };")
// 	n.csql("ALTER KEYSPACE \"svc_monitor_keyspace\" WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : " + rf + " };")
// 	n.csql("ALTER KEYSPACE \"config_db_uuid\" WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : " + rf + " };")
// }

type FuelNodeJson struct {
	ID       int                 `json:"id"`
	IP       string              `json:"ip"`
	Hostname string              `json:"hostname"`
	Networks []FuelNodeIfaceJson `json:"network_data"`
	Roles    string              `json:"roles"`
}

type FuelNodeIfaceJson struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}

func (fnj FuelNodeJson) is_contrail_db() bool {
	roles := strings.Split(fnj.Roles, ", ")
	for _, k := range roles {
		if strings.Compare(k, "contrail-db") == 0 {
			return true
		}
	}
	return false
}

func (fnj FuelNodeJson) get_contrail_ip() (string, error) {
	for _, nm := range fnj.Networks {
		if strings.Compare(nm.Name, "management") == 0 {
			return nm.IP, nil
		}
	}
	return "", errors.New("IP not found")
}
