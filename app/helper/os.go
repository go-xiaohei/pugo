package helper

import (
	"github.com/Unknwon/com"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

// remove all sub dirs and files in directory
func RemoveDir(dir string) error {
	if !com.IsDir(dir) {
		return nil
	}
	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, d := range dirs {
		if d.IsDir() {
			if err = RemoveDir(path.Join(dir, d.Name())); err != nil {
				return err
			}
		}
	}
	return os.RemoveAll(dir)
}

// copy file
func CopyFile(srcFile, destFile string) error {
	os.MkdirAll(filepath.Dir(destFile), os.ModePerm)

	f1, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer f1.Close()

	f2, err := os.OpenFile(destFile, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer f2.Close()

	_, err = io.Copy(f2, f1)
	return err
}
