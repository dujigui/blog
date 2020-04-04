package visitor

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
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
		log.Fatal(err)
	}
}

type QQState struct {
	Timestamp int64
	Redirect  string
}

func EncodeState(qqState QQState) (string, error) {
	d, err := json.Marshal(qqState)
	if err != nil {
		return "", err
	}
	d = aesEncrypt(d)
	return base64.StdEncoding.EncodeToString(d), nil
}

func DecodeState(s string) (QQState, error) {
	var qqState QQState
	d, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return qqState, err
	}
	d = aesDecrypt(d)
	if err = json.Unmarshal(d, &qqState); err != nil {
		return qqState, err
	}
	return qqState, nil
}

func aesEncrypt(plaintext []byte) []byte {
	stream := cipher.NewCTR(block, []byte(aesIv)[:16])
	stream.XORKeyStream(plaintext, plaintext)
	return plaintext
}

func aesDecrypt(ciptext []byte) []byte {
	stream := cipher.NewCTR(block, []byte(aesIv)[:16])
	stream.XORKeyStream(ciptext, ciptext)
	return ciptext
}
