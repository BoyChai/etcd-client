package main

import (
	"etcd-client/etcd"
	"fmt"
)

func main() {
	fmt.Println("========================")
	etcd.KvPUTDemo()
	fmt.Println("========================")
	etcd.KvGETDemo()
	fmt.Println("========================")
	etcd.KvDELDemo()
	defer etcd.CLI.Close()
}
