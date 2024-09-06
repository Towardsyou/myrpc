package do_nothing_compressor

type DoNothingCompressor struct{}

func (c *DoNothingCompressor) Id() uint8 {
	return 0
}

func (c *DoNothingCompressor) Compress(data []byte) ([]byte, error) {
	return data, nil
}

func (c *DoNothingCompressor) Extract(data []byte) ([]byte, error) {
	return data, nil
}
