package main

import (
	"context"
	"etcd-client/discovery"
	"etcd-client/discovery/proto"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getServerAdd(svcName string) string {
	s := discovery.ServiceDiscovery(svcName)
	// if s == nil {
	// 	return ""
	// }
	if s.IP == "" && s.Port == "" {
		return ""
	}
	return s.IP + ":" + s.Port
}

func SayHello() {
	addr := getServerAdd("Hello.Greeter")
	if addr == "" {
		log.Println("未发现可用服务")
		return
	}
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		log.Println(err)
		return
	}

	c := proto.NewGreeterClient(conn)
	r, err := c.SayHello(context.Background(), &proto.HelloRequest{
		Msg: "Hello",
	})

	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(r.Msg)
}
func main() {
	log.SetFlags(log.Llongfile)
	go discovery.WatchServiceName("Hello.Greeter")
	for {
		SayHello()
		time.Sleep(time.Second * 2)
	}
}
