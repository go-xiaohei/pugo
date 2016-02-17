package helper

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// Md5 generate md5 hash from string
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// Md5File generate file hash
func Md5File(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	var buff *bufio.Reader
	if fi, _ := f.Stat(); fi != nil && fi.Size() > 1024*100 {
		buff = bufio.NewReaderSize(f, 102400)
	} else {
		buff = bufio.NewReader(f)
	}
	for {
		line, err := buff.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		h.Write([]byte(line))
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
