package myrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"time"

	"github.com/Towardsyou/myrpc/internal/compress"
	"github.com/Towardsyou/myrpc/internal/message"
	"github.com/Towardsyou/myrpc/internal/serialize"
)

type Client struct {
	server_ip   string
	server_port string
}

func NewClient(server_ip, server_port string) *Client {
	return &Client{
		server_ip:   server_ip,
		server_port: server_port,
	}
}

func (c *Client) Do() error {
	conn, err := net.DialTimeout("tcp", c.server_ip+c.server_port, time.Second*3)
	if err != nil {
		return err
	}

	req := message.NewRequest([]byte("{\"X\": 1, \"Y\": 2}"), "example-service", "Add", map[string]string{"a/b": "123"}, 0, 0)
	err = message.WriteMsg(conn, message.EncodeReq(req))
	if err != nil {
		return err
	}

	rsp, err := message.ReadMsg(conn)
	if err != nil {
		return err
	}
	resp := message.DecodeResp(rsp)

	fmt.Println("In response:", string(resp.Error), string(resp.Data))
	conn.Close()
	return nil
}

func (c *Client) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	conn, err := net.DialTimeout("tcp", c.server_ip+c.server_port, time.Second*3)
	if err != nil {
		return nil, err
	}

	msg := message.NewRequest(req.Data, req.ServiceName, req.Method, req.Meta, req.Compressor, req.Serializer)
	err = message.WriteMsg(conn, message.EncodeReq(msg))
	if err != nil {
		return nil, err
	}

	rsp, err := message.ReadMsg(conn)
	if err != nil {
		return nil, err
	}
	resp := message.DecodeResp(rsp)

	fmt.Println("In response:", string(resp.Error), string(resp.Data))
	conn.Close()
	return resp, nil
}

func (client *Client) InitService(val Service, c compress.Compressor, serializer serialize.Serializer) error {
	v := reflect.ValueOf(val)
	ele := v.Elem()
	t := ele.Type()
	numField := t.NumField()
	for i := 0; i < numField; i++ {
		field := t.Field(i)
		fieldValue := ele.Field(i)
		if fieldValue.CanSet() {
			fn := func(args []reflect.Value) []reflect.Value {
				in := args[1].Interface()
				out := reflect.Zero(field.Type.Out(0))
				serialized, err := serializer.Serialize(in)
				if err != nil {
					return []reflect.Value{out, reflect.ValueOf(err)}
				}
				compressed, err := c.Compress(serialized)
				if err != nil {
					return []reflect.Value{out, reflect.ValueOf(err)}
				}
				ctx := args[0].Interface().(context.Context)
				meta := make(map[string]string, 2)
				if deadline, ok := ctx.Deadline(); ok {
					meta["deadline"] = strconv.FormatInt(deadline.UnixMilli(), 10)
				}
				req := message.NewRequest(compressed, val.ServiceName(), field.Name, meta, c.Id(), serializer.Id())
				resp, err := client.Invoke(ctx, req)
				if err != nil {
					return []reflect.Value{out, reflect.ValueOf(err)}
				}

				var retErr reflect.Value
				if len(resp.Error) > 0 {
					retErr = reflect.ValueOf(errors.New(string(resp.Error)))
				} else {
					retErr = reflect.Zero(reflect.TypeOf(new(error)).Elem())
				}
				if len(resp.Data) > 0 {
					out = reflect.New(field.Type.Out(0).Elem())
					data, err := c.Extract(resp.Data)
					if err != nil {
						return []reflect.Value{out, reflect.ValueOf(err)}
					}
					err = serializer.Deserialize(data, out.Interface())
					if err != nil {
						return []reflect.Value{out, reflect.ValueOf(err)}
					}
				}
				return []reflect.Value{out, retErr}
			}
			fieldValue.Set(reflect.MakeFunc(field.Type, fn))
		}
	}
	return nil
}
