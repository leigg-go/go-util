package _file

import (
	"fmt"
	"io/ioutil"
	"os"
)

func Each(dir string, fn func(seq int, info os.FileInfo) bool) error {
	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for i, dir := range dirs {
		if !fn(i, dir) {
			break
		}
	}
	return nil
}

// ListDir list files in given dir, flag support "file | folder"
// NOTE: os.FileInfo is interface, does not need to be a pointer
func ListDir(dir string, flag ...string) ([]os.FileInfo, error) {
	var files []os.FileInfo
	err := Each(dir, func(seq int, info os.FileInfo) bool {
		if len(flag) == 0 {
			files = append(files, info)
			return true
		}

		switch flag[0] {
		case "file":
			if !info.IsDir() {
				files = append(files, info)
			}
		case "folder":
			if info.IsDir() {
				files = append(files, info)
			}
		default:
			panic(fmt.Sprintf("_file: unknown flag <%s>", flag[0]))
		}

		return true
	})
	return files, err
}

// IsExist return false if path not exist or err not equal nil
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func IsFile(fPath string) (bool, error) {
	fi, e := os.Stat(fPath)
	if e != nil {
		return false, e
	}
	return !fi.IsDir(), nil
}

func MkdirIfNotExist(fPath string) error {
	if IsExist(fPath) {
		return nil
	}
	return os.Mkdir(fPath, 0755)
}
