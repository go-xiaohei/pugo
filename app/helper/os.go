package helper

import (
	"github.com/Unknwon/com"
	"io/ioutil"
	"os"
	"path"
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
