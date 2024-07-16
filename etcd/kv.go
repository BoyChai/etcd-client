package etcd

import (
	"context"
	"fmt"
	"log"
)

func KvDemo() {
	// PUT要上传三个参数 context、key、value、h中的OPTIONS
	// OPTIONS是etcdctl put -h中的OPTIONS
	// 例如--prev-kv就代表clientv3.WithPrevKV()
	putRes, err := CLI.Put(context.Background(), "/Test", "AAA")
	// putRes, err := cli.Put(context.Background(), "/Test", "AAA", clientv3.WithPrevKV())

	if err != nil {
		log.Fatalln(err)
	}
	if putRes.PrevKv != nil {
		log.Println(putRes.PrevKv)
		log.Println(string(putRes.PrevKv.Key))
		log.Println(string(putRes.PrevKv.Value))
	}
	fmt.Println("OK")
}
