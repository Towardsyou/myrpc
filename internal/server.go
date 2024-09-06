package myrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/Towardsyou/myrpc/internal/compress"
	"github.com/Towardsyou/myrpc/internal/compress/do_nothing_compressor"
	"github.com/Towardsyou/myrpc/internal/message"
	"github.com/Towardsyou/myrpc/internal/serialize"
	"github.com/Towardsyou/myrpc/internal/serialize/json"
)

type Server struct {
	server_ip   string
	server_port string
	services    map[string]*ServiceStub
	serializers []serialize.Serializer
	compressors []compress.Compressor
}

func NewServer(server_ip, server_port string) *Server {
	s := Server{
		server_ip:   server_ip,
		server_port: server_port,
		serializers: make([]serialize.Serializer, 128),
		compressors: make([]compress.Compressor, 128),
		services:    make(map[string]*ServiceStub),
	}
	s.RegisterSerializer(&json.JsonSerializer{})
	s.RegisterCompressor(&do_nothing_compressor.DoNothingCompressor{})
	return &s
}

func (s *Server) RegisterSerializer(serializer serialize.Serializer) {
	s.serializers[serializer.Id()] = serializer
}

func (s *Server) RegisterCompressor(c compress.Compressor) {
	s.compressors[c.Id()] = c
}

type Service interface {
	ServiceName() string
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.server_ip+s.server_port)
	if err != nil {
		return err
	}
	fmt.Println("Register services")

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		content, err := message.ReadMsg(conn)
		if err != nil {
			return err
		}
		req := message.DecodeReq(content)

		fmt.Println("Input: ", req.ServiceName, req.Method, req.Meta, string(req.Data))

		service := s.services[req.ServiceName]
		resp, err := service.Call(req)
		if err != nil {
			resp = message.NewResponse([]byte("execution error"), []byte(err.Error()), 0, 0)
		}
		message.WriteMsg(conn, message.EncodeResp(resp))
	}
}

type ServiceStub struct {
	Methods     map[string]reflect.Value
	Serializers []serialize.Serializer
	Compressors []compress.Compressor
}

func (s *Server) RegisterService(service Service) {
	typ := reflect.TypeOf(service)
	val := reflect.ValueOf(service)
	methods := make(map[string]reflect.Value, val.NumMethod())
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		methods[method.Name] = val.MethodByName(method.Name)
	}
	s.services[service.ServiceName()] = &ServiceStub{
		Methods:     methods,
		Serializers: s.serializers,
		Compressors: s.compressors,
	}
}

func (s *ServiceStub) Call(req *message.Request) (*message.Response, error) {
	method_fn, ok := s.Methods[req.Method]
	if !ok {
		return nil, errors.New("method with name " + req.Method + " not found")
	}
	inType := method_fn.Type().In(1)
	arg := reflect.New(inType.Elem()).Interface()
	serializer := s.Serializers[req.Serializer]
	if err := serializer.Deserialize(req.Data, arg); err != nil {
		return nil, err
	}

	ctx := context.Background()
	ret := method_fn.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(arg)})

	body, err := serializer.Serialize(ret[0].Interface())
	if err != nil {
		return nil, err
	}
	error_msg := []byte{}
	if !ret[1].IsZero() {
		error_msg = []byte(ret[1].Interface().(error).Error())
	}
	return message.NewResponse(body, error_msg, 0, 0), nil
}
