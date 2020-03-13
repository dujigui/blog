package files

import (
	"crypto/md5"
	"encoding/hex"
	. "github.com/dujigui/blog/services/cfg"
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/utils"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	dir := Config().GetString("files")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		Logger().Fatal("upload", "创建文件目录失败", Params{"dir": dir, "err": err})
	}
}

func Save(in io.Reader, original string) (string, error) {
	hashed := hash(original)
	dir := Config().GetString("files")
	out, err := os.OpenFile(filepath.Join(dir, hashed), os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return hashed, err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return hashed, err
}

func Remove(fn string) error {
	dir := Config().GetString("files")
	fp := filepath.Join(dir, fn)
	return os.Remove(fp)
}

func hash(original string) (hashed string) {
	ext := filepath.Ext(original)
	original = strings.TrimSuffix(original, ext)

	hasher := md5.New()
	hasher.Write([]byte(original))
	hashed = hex.EncodeToString(hasher.Sum(nil)) + ext
	return
}
