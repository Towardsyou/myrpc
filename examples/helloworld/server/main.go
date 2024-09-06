package main

import (
	"github.com/Towardsyou/myrpc/examples/helloworld"
	myrpc "github.com/Towardsyou/myrpc/internal"
)

const (
	server_ip   = "localhost"
	server_port = ":48521"
)

func main() {
	s := myrpc.NewServer(server_ip, server_port)
	srv := helloworld.ExampleService{}
	srv2 := helloworld.StrAddService{}
	s.RegisterService(srv)
	s.RegisterService(srv2)
	
	s.Start()
}
