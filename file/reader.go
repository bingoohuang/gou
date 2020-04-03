package file

import (
	"io/ioutil"
	"strings"
)

// ToBytes reads the file to a byte slice.
func ToBytes(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

// ToString reads the file content to string.
func ToString(filePath string) (string, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// ToTrimString reads the file to a  trimmed string.
func ToTrimString(filePath string) (string, error) {
	str, err := ToString(filePath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(str), nil
}
