package services

import (
	"encoding/json"
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/utils"
	"log"
	"os"
	"path/filepath"
)

const (
	prefName = "pref.json"
)

var p = &Preferences{}

type Preferences struct {
	BlogName      string `json:"blog_name" form:"blog_name"`
	AdminPageName string `json:"admin_page_name" form:"admin_page_name"`
	AboutPostID   int    `json:"about_id" form:"about_id"`
	Init          int64  `json:"init"`
	Email         string `json:"email"`
	QQAppID       string `json:"qq_app_id"`
	QQAppKey      string `json:"qq_app_key"`
	QQRedirect    string `json:"qq_redirect"`
	Salt          string `json:"salt"`
}

func init() {
	dir := "data/favicon/"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatal("pref ", "创建 favicon 目录失败 ", Params{"dir": dir}.Err(err))
	}

	dir = "data/pref/"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatal("pref ", "创建配置目录失败 ", Params{"dir": dir}.Err(err))
	}

	f, err := os.Open(filepath.Join(dir, prefName))
	if err != nil {
		Logger().Warning("pref", "配置文件不存在", Params{"dir": dir, "preName": prefName, "err": err})
		return
	}

	if err = json.NewDecoder(f).Decode(p); err != nil {
		Logger().Warning("pref", "读取配置文件失败", Params{"dir": dir, "preName": prefName, "err": err})
	} else {
		Logger().Trace("pref", "读取配置文件成功", Params{"dir": dir, "preName": prefName})
	}
}

// todo 初始化时保存配置
func (p *Preferences) Save() error {
	dir := "data/pref/"
	f, err := os.OpenFile(filepath.Join(dir, prefName), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		Logger().Fatal("pref", "配置文件无法打开", Params{"dir": dir, "preName": prefName, "err": err})
		return err
	}
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")

	// 清空配置文件
	if err := f.Truncate(0); err != nil {
		Logger().Warning("pref", "写入配置文件失败", Params{"dir": dir, "preName": prefName, "err": err})
		return err
	}

	err = encoder.Encode(p)
	if err != nil {
		Logger().Warning("pref", "写入配置文件失败", Params{"dir": dir, "preName": prefName, "err": err})
		return err
	}
	Logger().Trace("pref", "写入配置文件成功", Params{"dir": dir, "preName": prefName})
	return nil
}

func Pref() *Preferences {
	return p
}
