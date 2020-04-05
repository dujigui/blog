package users

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

const (
	CookieExpire = 24 * time.Hour
)

type c struct {
	ID     int
	Admin  bool
	Expire int64
}

func CreateToken(id int, admin bool, duration time.Duration) string {
	d, err := json.Marshal(c{
		ID:     id,
		Admin:  admin,
		Expire: time.Now().Add(duration).Unix(),
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
	if c.Expire == 0 || c.Expire < time.Now().Unix() {
		return false, 0, false
	}
	return true, c.ID, c.Admin
}
