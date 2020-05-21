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
