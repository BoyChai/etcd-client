package etcd

import (
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var CLI *clientv3.Client

func init() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   GetEtcdEndpoints(),
		DialTimeout: 5 * time.Second,
		Username:    "root",
		Password:    "Qwer1234",
	})
	if err != nil {
		log.Fatalln(err)
	}
	CLI = cli
}

func GetEtcdEndpoints() []string {
	return []string{"192.168.29.130:12379", "192.168.29.130:22379", "192.168.29.130:32379"}
}
