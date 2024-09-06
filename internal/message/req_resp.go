package message

import (
	"bytes"
	"encoding/binary"
	"sync/atomic"
)

const (
	kvSeparator  = ':'
	rowSeparator = ','
)

var messageId uint32 = 0

type Request struct {
	HeadLength  uint32
	BodyLength  uint32
	MessageId   uint32
	Version     uint8
	Compressor  uint8
	Serializer  uint8
	ServiceName string
	Method      string
	Meta        map[string]string
	Data        []byte
}

func (req *Request) SetHeadLength() {
	// fixed length
	res := 15
	res += len(req.ServiceName)
	// row separator
	res++
	res += len(req.Method)
	// row separator
	res++
	for key, value := range req.Meta {
		res += len(key)
		// kv separator
		res++
		res += len(value)
		// row separator
		res++
	}
	req.HeadLength = uint32(res)
}

func EncodeReq(req *Request) []byte {
	bs := make([]byte, req.HeadLength+req.BodyLength)
	cur := bs
	// HeaderLength
	binary.BigEndian.PutUint32(cur[:4], req.HeadLength)
	cur = cur[4:]
	// BodyLength
	binary.BigEndian.PutUint32(cur[:4], req.BodyLength)
	cur = cur[4:]
	// MessageId
	binary.BigEndian.PutUint32(cur[:4], req.MessageId)
	cur = cur[4:]
	// Version
	cur[0] = req.Version
	cur = cur[1:]
	// Compresser
	cur[0] = req.Compressor
	cur = cur[1:]
	// Serializer
	cur[0] = req.Serializer
	cur = cur[1:]
	// ServiceName
	copy(cur, []byte(req.ServiceName))
	cur[len(req.ServiceName)] = rowSeparator
	cur = cur[len(req.ServiceName)+1:]
	// Method
	copy(cur, []byte(req.Method))
	cur[len(req.Method)] = rowSeparator
	cur = cur[len(req.Method)+1:]
	// Meta
	for k, v := range req.Meta {
		t := []byte(k)
		copy(cur, t)
		cur[len(t)] = kvSeparator
		cur = cur[len(t)+1:]
		t = []byte(v)
		copy(cur, t)
		cur[len(t)] = rowSeparator
		cur = cur[len(t)+1:]
	}
	// Data
	copy(cur, req.Data)
	return bs
}

func DecodeReq(bs []byte) *Request {
	req := Request{}
	req.HeadLength = binary.BigEndian.Uint32(bs[:4])
	req.BodyLength = binary.BigEndian.Uint32(bs[4:8])
	req.MessageId = binary.BigEndian.Uint32(bs[8:12])
	req.Version = bs[12]
	req.Compressor = bs[13]
	req.Serializer = bs[14]
	// 7. service name and method name
	meta := bs[15:req.HeadLength]
	index := bytes.IndexByte(meta, rowSeparator)
	req.ServiceName = string(meta[:index])
	meta = meta[index+1:]
	index = bytes.IndexByte(meta, rowSeparator)
	req.Method = string(meta[:index])
	meta = meta[index+1:]
	// key value in meta
	for len(meta) > 0 {
		// 4 can be increased when more possible meta available
		metaMap := make(map[string]string, 4)
		index = bytes.IndexByte(meta, rowSeparator)
		for index != -1 {
			pair := meta[:index]
			pairIndex := bytes.IndexByte(meta, kvSeparator)
			metaMap[string(pair[:pairIndex])] = string(pair[pairIndex+1:])
			meta = meta[index+1:]
			index = bytes.IndexByte(meta, rowSeparator)
		}
		req.Meta = metaMap
	}

	req.Data = bs[req.HeadLength:]
	return &req
}

func NewRequest(body []byte, serviceName string, methodName string, meta map[string]string, compressor uint8, serializer uint8) *Request {
	r := Request{
		HeadLength:  0,
		BodyLength:  uint32(len(body)),
		MessageId:   atomic.AddUint32(&messageId, 1),
		Version:     0,
		Compressor:  compressor,
		Serializer:  serializer,
		ServiceName: serviceName,
		Method:      methodName,
		Meta:        meta,
		Data:        body,
	}
	r.SetHeadLength()
	return &r
}

type Response struct {
	HeadLength uint32
	BodyLength uint32
	MessageId  uint32
	Version    uint8
	Compressor uint8
	Serializer uint8
	Error      []byte
	Data       []byte
}

func EncodeResp(resp *Response) []byte {
	bs := make([]byte, resp.HeadLength+resp.BodyLength)
	cur := bs
	// HeaderLength
	binary.BigEndian.PutUint32(cur[:4], resp.HeadLength)
	cur = cur[4:]
	// BodyLength
	binary.BigEndian.PutUint32(cur[:4], resp.BodyLength)
	cur = cur[4:]
	// MessageId
	binary.BigEndian.PutUint32(cur[:4], resp.MessageId)
	cur = cur[4:]
	// Version
	cur[0] = resp.Version
	cur = cur[1:]
	// Compressor
	cur[0] = resp.Compressor
	cur = cur[1:]
	// Serializer
	cur[0] = resp.Serializer
	cur = cur[1:]
	// Error
	copy(cur, resp.Error)
	cur = cur[len(resp.Error):]
	// Data
	copy(cur, resp.Data)
	return bs
}

func (resp *Response) SetHeadLength() {
	resp.HeadLength = 15 + uint32(len(resp.Error))
}

func DecodeResp(bs []byte) *Response {
	resp := Response{}
	resp.HeadLength = binary.BigEndian.Uint32(bs[0:4])
	resp.BodyLength = binary.BigEndian.Uint32(bs[4:8])
	resp.MessageId = binary.BigEndian.Uint32(bs[8:12])
	resp.Version = bs[12]
	resp.Compressor = bs[13]
	resp.Serializer = bs[14]
	resp.Error = bs[15:resp.HeadLength]
	resp.Data = bs[resp.HeadLength:]
	return &resp
}

func NewResponse(body []byte, err []byte, compressor uint8, serializer uint8) *Response {
	r := Response{
		HeadLength: 0,
		BodyLength: uint32(len(body)),
		Version:    0,
		Compressor: compressor,
		Serializer: serializer,
		Error:      err,
		Data:       body,
	}
	r.SetHeadLength()
	return &r
}
