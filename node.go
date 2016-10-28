package main

import (
	"errors"
	"log"
	"strconv"
)

type Node struct {
	hostname     string
	admin_ip     string
	cassandra_ip string
	ssh          SSHNode
}

func (n *Node) run(command string) string {
	log.Println("At " + n.hostname + " run: " + command)
	n.ssh = SSHNode{username: "fuel", host: "ip"}
	return n.ssh.Exec(command)
}

func (n *Node) csql(command string) {
	remote_exec := "cqlsh " + n.cassandra_ip + " -e \"" + command + "\""
	n.run(remote_exec)
}

func (n *Node) set_replication_factor(factor int) {
	rf := strconv.Itoa(factor)
	log.Println("Set replica_factor to " + rf)
	n.csql("ALTER KEYSPACE \"to_bgp_keyspace\" WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : " + rf + " };")
	n.csql("ALTER KEYSPACE \"svc_monitor_keyspace\" WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : " + rf + " };")
	n.csql("ALTER KEYSPACE \"config_db_uuid\" WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : " + rf + " };")
}

type FuelNodeJson struct {
	ip           string
	hostname     string
	network_data []FuelNodeIfaceJson
	roles        []string
}
type FuelNodeIfaceJson struct {
	name string
	ip   string
}

func (fnj FuelNodeJson) is_contrail_db() bool {
	for _, k := range fnj.roles {
		if k == "contrail-db" {
			return true
		}
	}
	return false
}

func (fnj FuelNodeJson) get_contrail_ip() (string, error) {
	for _, nm := range fnj.network_data {
		if nm.name == "management" {
			return nm.ip, nil
		}
	}
	return "", errors.New("IP not found")
}
