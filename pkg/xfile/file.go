/**
 * @Author: yangon
 * @Description
 * @Date: 2020/12/22 19:31
 **/
package xfile

import (
	"os"
	"runtime"
	"strings"
)

func MkdirIfNecessary(createDir string) error {
	var (
		path string
		err  error
	)
	if os.IsPathSeparator('\\') {
		path = "\\"
	} else {
		path = "/"
	}

	s := strings.Split(createDir, path)
	startIndex := 0
	dir := ""
	if s[0] == "" {
		startIndex = 1
	} else {
		dir, _ = os.Getwd()
	}
	for i := startIndex; i < len(s); i++ {
		d := dir + path + strings.Join(s[startIndex:i+1], path)
		if _, e := os.Stat(d); os.IsNotExist(e) {
			err = os.Mkdir(d, os.ModePerm)
			if err != nil {
				break
			}
		}
	}
	return err
}

func CheckAndGetParentDir(path string) string {
	// check path is the directory
	isDir, err := IsDirectory(path)
	if err != nil || isDir {
		return path
	}
	return getParentDirectory(path)
}

func getParentDirectory(dirctory string) string {
	if runtime.GOOS == "windows" {
		dirctory = strings.Replace(dirctory, "\\", "/", -1)
	}
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

// IsDirectory ...
func IsDirectory(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	switch mode := f.Mode(); {
	case mode.IsDir():
		return true, nil
	case mode.IsRegular():
		return false, nil
	}
	return false, nil
}
