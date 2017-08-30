package helpers

import (
	"crypto/md5"
	"fmt"
)

func Md5(data []byte) string {
	h := md5.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}
