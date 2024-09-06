package json

import "encoding/json"

type JsonSerializer struct{}

func (s *JsonSerializer) Id() uint8 {
	return 0
}

func (s *JsonSerializer) Serialize(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (s *JsonSerializer) Deserialize(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
