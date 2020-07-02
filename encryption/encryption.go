package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/klauspost/compress/zstd"
	"io"
)

func Encrypt(key, data []byte) (encrypted []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, aes.BlockSize+len(data))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], data)

	return cipherText, nil
}

func Decrypt(key, encrypted []byte) (decrypted []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)

	return encrypted, nil
}

func Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := zstd.NewWriter(&buf)
	if err != nil {
		return nil, err
	}

	if _, err := writer.Write(data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Decompress(compressed []byte) ([]byte, error) {
	reader, err := zstd.NewReader(bytes.NewBuffer(compressed))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
