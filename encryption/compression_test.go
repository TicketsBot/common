package encryption

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCompressDecompress(t *testing.T) {
	bytes := []byte("Hello, World!")

	compressed := Compress(bytes)
	decompressed, err := Decompress(compressed)
	require.NoError(t, err)

	require.Equal(t, bytes, decompressed)
}