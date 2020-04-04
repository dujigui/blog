package visitor

import (
	"encoding/base64"
	"encoding/json"
	. "github.com/dujigui/blog/services/users"
)



type QQState struct {
	Timestamp int64
	Redirect  string
}

func EncodeState(qqState QQState) (string, error) {
	d, err := json.Marshal(qqState)
	if err != nil {
		return "", err
	}
	d = AESEncrypt(d)
	return base64.StdEncoding.EncodeToString(d), nil
}

func DecodeState(s string) (QQState, error) {
	var qqState QQState
	d, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return qqState, err
	}
	d = AESDecrypt(d)
	if err = json.Unmarshal(d, &qqState); err != nil {
		return qqState, err
	}
	return qqState, nil
}

