package main

import (
	"etcd-client/etcd"
)

func main() {
	etcd.KvDemo()
	defer etcd.CLI.Close()
}
