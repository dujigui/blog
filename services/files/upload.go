package files

import (
	"crypto/md5"
	"encoding/hex"
	. "github.com/dujigui/blog/utils"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	dir := "data/files/"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatal("upload ", "创建文件目录失败 ", Params{"dir": dir}.Err(err))
	}
}

func Save(in io.Reader, original string) (string, error) {
	hashed := hash(original)
	dir := "data/files/"
	out, err := os.OpenFile(filepath.Join(dir, hashed), os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return hashed, err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return hashed, err
}

func Remove(fn string) error {
	dir := "data/files/"
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
