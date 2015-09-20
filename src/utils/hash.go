// Package utils provides common functions
package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

// md5 encode
func Md5String(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

// sha256 encode
func Sha256String(str string) string {
	s := sha256.New()
	s.Write([]byte(str))
	return hex.EncodeToString(s.Sum(nil))
}
