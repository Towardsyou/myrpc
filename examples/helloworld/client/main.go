package main

import (
	"context"
	"fmt"

	"github.com/Towardsyou/myrpc/examples/helloworld"
	myrpc "github.com/Towardsyou/myrpc/internal"
	"github.com/Towardsyou/myrpc/internal/compress/do_nothing_compressor"
	"github.com/Towardsyou/myrpc/internal/serialize/json"
)

const (
	server_ip   = "localhost"
	server_port = ":48521"
)

func main() {
	c := myrpc.NewClient(
		server_ip,
		server_port,
	)
	srv := helloworld.ExampleService{}
	srv2 := helloworld.StrAddService{}
	c.InitService(&srv, &do_nothing_compressor.DoNothingCompressor{}, &json.JsonSerializer{})
	c.InitService(&srv2, &do_nothing_compressor.DoNothingCompressor{}, &json.JsonSerializer{})

	in := helloworld.AddArgs{X: 2, Y: 4}
	ctx := context.Background()
	out, err := srv.Add(ctx, &in)
	if err != nil {
		fmt.Println("return err:", err)
	}
	fmt.Printf("return value: %#v\n", out)

	in2 := helloworld.StrAddArgs{X: "Hello ", Y: "world!"}
	ctx2 := context.Background()
	out2, err := srv2.Add(ctx2, &in2)
	if err != nil {
		fmt.Println("return err:", err)
	}
	fmt.Printf("return value: %#v\n", out2)

	fmt.Println("exited")
}
