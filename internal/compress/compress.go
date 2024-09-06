package compress

type Compressor interface {
	Id() uint8
	Compress([]byte) ([]byte, error)
	Extract([]byte) ([]byte, error)
}
