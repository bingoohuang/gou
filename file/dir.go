package file

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// BaseDir returns the common directory for a slice of directories.
func BaseDir(dirs []string) string {
	baseDir := ""

	for _, dir := range dirs {
		d := filepath.Dir(dir)

		if baseDir == "" {
			baseDir = d
		} else {
			for !strings.HasPrefix(d, baseDir) {
				baseDir = filepath.Dir(baseDir)
			}
		}

		if baseDir == "/" {
			break
		}
	}

	if baseDir == "" {
		baseDir = "/"
	}

	return baseDir
}

// IsExist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func IsExist(fp string) bool {
	_, err := os.Stat(fp)
	return err == nil || os.IsExist(err)
}

// DirsUnder list dirs under dirPath
func DirsUnder(dirPath string) ([]string, error) {
	if !IsExist(dirPath) {
		return []string{}, nil
	}

	fs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return []string{}, err
	}

	sz := len(fs)
	if sz == 0 {
		return []string{}, nil
	}

	ret := make([]string, 0, sz)

	for i := 0; i < sz; i++ {
		if fs[i].IsDir() {
			name := fs[i].Name()
			if name != "." && name != ".." {
				ret = append(ret, name)
			}
		}
	}

	return ret, nil
}

// InsureDir insure dir exist
func InsureDir(fp string) error {
	if IsExist(fp) {
		return nil
	}

	return os.MkdirAll(fp, os.ModePerm)
}
