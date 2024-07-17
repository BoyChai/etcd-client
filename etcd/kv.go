package etcd

import (
	"context"
	"log"

	"go.etcd.io/etcd/clientv3"
)

func KvPUTDemo() {
	// PUT要上传三个参数 context、key、value、h中的OPTIONS
	// OPTIONS是etcdctl put -h中的OPTIONS
	// 例如--prev-kv就代表clientv3.WithPrevKV()
	// putRes, err := CLI.Put(context.Background(), "/Test", "AAA")
	putRes, err := CLI.Put(context.Background(), "/Test", "AAA", clientv3.WithPrevKV())

	if err != nil {
		log.Fatalln(err)
	}
	if putRes.PrevKv != nil {
		log.Println(putRes.PrevKv)
		log.Println(string(putRes.PrevKv.Key))
		log.Println(string(putRes.PrevKv.Value))
	}
}

func KvGETDemo() {
	getRes, err := CLI.Get(context.Background(), "/Test")
	// getRes, err := CLI.Get(context.Background(), "/Test", clientv3.WithPrefix())

	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(getRes.Kvs[0].Key))
	log.Println(string(getRes.Kvs[0].Value))
}

func KvDELDemo() {
	delRes, err := CLI.Delete(context.Background(), "/Test", clientv3.WithPrevKV())
	if err != nil {
		log.Fatalln(err)
	}
	if delRes.Deleted > 0 {
		for _, v := range delRes.PrevKvs {
			log.Println(string(v.Key))
			log.Println(string(v.Value))
		}
	}
}
