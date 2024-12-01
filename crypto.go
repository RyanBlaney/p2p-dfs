package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

func hashKey(key string) string {
	hash := md5.Sum([]byte(key))

	return hex.EncodeToString(hash[:])
}

func newEncryptionKey() []byte {
	keyBuf := make([]byte, 32)
	io.ReadFull(rand.Reader, keyBuf)
	return keyBuf
}

func copyStream(stream cipher.Stream, blockSize int, src io.Reader, dest io.Writer) (int, error) {
	var (
		buf = make([]byte, 32*1024)
		nw  = blockSize
	)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf, buf[:n])
			nn, err := dest.Write(buf[:n])
			if err != nil {
				return 0, err
			}
			nw += nn
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return 0, err
		}
	}
	return nw, nil
}

func copyEncrypt(key []byte, src io.Reader, dest io.Writer) (int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return 0, err
	}

	// prepend the IV to the file
	if _, err := dest.Write(iv); err != nil {
		return 0, err
	}

	stream := cipher.NewCTR(block, iv)
	return copyStream(stream, block.BlockSize(), src, dest)
}

func copyDecrypt(key []byte, src io.Reader, dest io.Writer) (int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	// Read IV from the given io.Reader
	iv := make([]byte, block.BlockSize())
	if _, err := src.Read(iv); err != nil {
		return 0, err
	}

	stream := cipher.NewCTR(block, iv)
	return copyStream(stream, block.BlockSize(), src, dest)
}
