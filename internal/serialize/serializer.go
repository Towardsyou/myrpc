package serialize

type Serializer interface {
	Id() uint8
	Serialize(any) ([]byte, error)
	Deserialize([]byte, any) error
}
