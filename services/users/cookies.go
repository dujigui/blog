package users

import (
	"encoding/base64"
	"encoding/json"
)

type c struct {
	ID    int
	Admin bool
}

func CreateToken(id int, admin bool) string {
	d, err := json.Marshal(c{
		ID:    id,
		Admin: admin,
	})
	if err != nil {
		return ""
	}
	d = AESEncrypt(d)
	return base64.StdEncoding.EncodeToString(d)
}

func ParseToken(token string) (ok bool, id int, admin bool) {
	var c c
	d, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false, 0, false
	}
	d = AESDecrypt(d)
	if err = json.Unmarshal(d, &c); err != nil {
		return false, 0, false
	}
	return true, c.ID, c.Admin
}
