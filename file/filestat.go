package file

import "os"

// ExistsEnum is the file exits enumeration.
type ExistsEnum int

const (
	// Exists means the file exists.
	Exists ExistsEnum = iota
	// NotExists means the file not exists.
	NotExists
	// Unknown means unknown state.
	Unknown
)

// Stat stats the file.
func Stat(name string) (ExistsEnum, error) {
	_, err := os.Stat(name)
	if err == nil {
		return Exists, nil
	}

	if os.IsNotExist(err) {
		return NotExists, nil
	}

	return Unknown, err
}
