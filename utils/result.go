package utils

import "github.com/kataras/iris/v12"

func Result(ok bool, msg string, result interface{}, kv ...interface{}) iris.Map {
	m := iris.Map{"ok": ok, "msg": msg, "result": result}
	for i:=0;i<len(kv);i+=2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}
