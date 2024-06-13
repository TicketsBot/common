package encryption

import (
	"github.com/klauspost/compress/zstd"
	"runtime"
)

var compressor, _ = zstd.NewWriter(nil,
	zstd.WithEncoderConcurrency(runtime.NumCPU()),
	zstd.WithWindowSize(2<<20),
)

func Compress(data []byte) []byte {
	return compressor.EncodeAll(data, make([]byte, 0, len(data)))
}

var decompressor, _ = zstd.NewReader(nil, zstd.WithDecoderConcurrency(runtime.NumCPU()))

func Decompress(data []byte) ([]byte, error) {
	return decompressor.DecodeAll(data, nil)
}
