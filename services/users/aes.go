package users

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	. "github.com/dujigui/blog/services"
	"log"
)

const (
	aesKey = "0thejNSWUhONW5IN3hk"
	aesIv  = "XZsejVUNmRxNi9lK1"
)

var (
	block cipher.Block
	key   []byte
	iv    []byte
)

func init() {
	h := sha256.New()
	h.Write([]byte(aesKey + Pref().Salt))
	key = h.Sum(nil)[:aes.BlockSize]
	h.Reset()

	h.Write([]byte(aesIv + Pref().Salt))
	iv = h.Sum(nil)[:aes.BlockSize]

	if key == nil || len(key) != aes.BlockSize || iv == nil || len(iv) != aes.BlockSize {
		log.Fatal("aes ", "初始化 key/iv 失败 ", fmt.Sprintf("key=%x iv=%x", key, iv))
	}

	var err error
	block, err = aes.NewCipher(key)
	if err != nil {
		log.Fatal("aes ", "初始化 Cipher 失败 ", err)
	}
}

func AESEncrypt(d []byte) []byte {
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(d, d)
	return d
}

func AESDecrypt(d []byte) []byte {
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(d, d)
	return d
}
