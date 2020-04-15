package users

import (
	"crypto/aes"
	"crypto/cipher"
	. "github.com/dujigui/blog/services"
	"log"
)

const (
	aesKey = "0thejNSWUhONW5IN3hk"
	aesIv  = "XZsejVUNmRxNi9lK1"
)

var block cipher.Block

func init() {
	var err error
	block, err = aes.NewCipher([]byte(aesKey + Pref().Salt)[:16])
	if err != nil {
		log.Fatal("aes ", "初始化加密失败 ", err)
	}
}

func AESEncrypt(d []byte) []byte {
	stream := cipher.NewCTR(block, []byte(aesIv)[:16])
	stream.XORKeyStream(d, d)
	return d
}

func AESDecrypt(d []byte) []byte {
	stream := cipher.NewCTR(block, []byte(aesIv)[:16])
	stream.XORKeyStream(d, d)
	return d
}