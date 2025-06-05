package stringhelper

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5 md5加密字符串
func Md5(content string) string {
	md5obj := md5.New()
	md5obj.Write([]byte(content))
	return hex.EncodeToString(md5obj.Sum(nil))
}
