package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/klauspost/compress/zstd"
	"io"
	"runtime"
)

func Encrypt(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	cipherText := gcm.Seal(nil, nonce, data, nil)
	return append(nonce, cipherText...), nil
}

func Decrypt(key, encrypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	nonce := encrypted[:12]
	cipherText := encrypted[12:]

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	decrypted, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return decrypted, err
}

var compressor, _ = zstd.NewWriter(nil, zstd.WithEncoderConcurrency(runtime.NumCPU()))
func Compress(data []byte) []byte {
	return compressor.EncodeAll(data, make([]byte, 0, len(data)))
}

var decompressor, _ = zstd.NewReader(nil, zstd.WithDecoderConcurrency(runtime.NumCPU()))
func Decompress(data []byte) ([]byte, error) {
	return decompressor.DecodeAll(data, nil)
}
