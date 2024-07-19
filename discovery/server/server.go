package main

import (
	"context"
	"etcd-client/discovery"
	"etcd-client/discovery/proto"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "")
)

type server struct{}

func (server) SayHello(context context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	fmt.Println(req.Msg)
	return &proto.HelloReply{
		Msg: "ok",
	}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalln(err)
	}
	s := grpc.NewServer()
	// proto.RegisterGreeterServer(s, &server{})
	serverRegister(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}

func serverRegister(s *grpc.Server, srv proto.GreeterServer) {
	proto.RegisterGreeterServer(s, srv)
	s1 := &discovery.Service{
		Name:     "Hello.Greeter",
		Port:     strconv.Itoa(*port),
		IP:       "127.0.0.1",
		Protocol: "grpc",
	}
	go discovery.ServiceRegister(s1)
}
