package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodeRequest(t *testing.T) {
	testCases := []struct {
		name string
		req  *Request
	}{
		{
			name: "with meta",
			req: &Request{
				MessageId:   123,
				Version:     12,
				Compressor:  25,
				Serializer:  17,
				ServiceName: "user-service",
				Method:      "GetById",
				Meta: map[string]string{
					"trace-id": "123",
					"a/b":      "b",
					"shadow":   "true",
				},
				Data: []byte("hello, world"),
			},
		},
		{
			name: "no meta",
			req: &Request{
				MessageId:   123,
				Version:     12,
				Compressor:  25,
				Serializer:  17,
				ServiceName: "user-service",
				Method:      "GetById",
				Data:        []byte("hello, world"),
			},
		},
		{
			name: "empty value",
			req: &Request{
				MessageId:   123,
				Version:     12,
				Compressor:  25,
				Serializer:  17,
				ServiceName: "user-service",
				Method:      "GetById",
				Meta: map[string]string{
					"trace-id": "123",
					"a/b":      "",
					"shadow":   "true",
				},
				Data: []byte("hello, world"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.req.SetHeadLength()
			tc.req.BodyLength = uint32(len(tc.req.Data))
			bs := EncodeReq(tc.req)
			req := DecodeReq(bs)
			assert.Equal(t, tc.req, req)
		})
	}
}

func TestEncodeDecodeResponse(t *testing.T) {
	testCases := []struct {
		name string
		resp *Response
	}{
		{
			name: "no error",
			resp: &Response{
				MessageId:  123,
				Version:    12,
				Compressor: 25,
				Serializer: 17,
				Data:       []byte("hello, world"),
				Error:      []byte{},
			},
		},
		{
			name: "with error",
			resp: &Response{
				MessageId:  123,
				Version:    12,
				Compressor: 25,
				Serializer: 17,
				Data:       []byte("hello, world"),
				Error:      []byte("error message"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.resp.SetHeadLength()
			tc.resp.BodyLength = uint32(len(tc.resp.Data))
			bs := EncodeResp(tc.resp)
			req := DecodeResp(bs)
			assert.Equal(t, tc.resp, req)
		})
	}
}
