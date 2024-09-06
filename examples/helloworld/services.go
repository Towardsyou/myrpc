package helloworld

import "context"

type AddArgs struct {
	X int
	Y int
}

type AddResp struct {
	Val int
}

type ExampleService struct{}

func (s ExampleService) ServiceName() string {
	return "example-service"
}

func (s ExampleService) Add(ctx context.Context, args *AddArgs) (*AddResp, error) {
	return &AddResp{Val: args.X + args.Y}, nil
}

type StrAddArgs struct {
	X string
	Y string
}

type StrAddResp struct {
	Val string
}

type StrAddService struct{}

func (s StrAddService) ServiceName() string {
	return "str-add-service"
}

func (s StrAddService) Add(ctx context.Context, args *StrAddArgs) (*StrAddResp, error) {
	return &StrAddResp{Val: args.X + args.Y}, nil
}