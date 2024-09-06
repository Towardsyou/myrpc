package message

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
)

const (
	peekLength = 8
)

func ReadMsg(conn net.Conn) ([]byte, error) {
	buffer := make([]byte, peekLength)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	} else if n != peekLength {
		err = errors.New("not enough bytes for header")
		return nil, err
	}
	headerLength := binary.BigEndian.Uint32(buffer[0:4])
	bodyLength := binary.BigEndian.Uint32(buffer[4:8])

	wholeMessage := make([]byte, headerLength + bodyLength)
	copy(wholeMessage[0:peekLength], buffer)
	_, err = io.ReadFull(conn, wholeMessage[peekLength:])
	return wholeMessage, err
}

func WriteMsg(conn net.Conn, message []byte) error {
	_, err := conn.Write(message)
	if err != nil {
		return err
	}
	return nil
}
