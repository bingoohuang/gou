package enc

import (
	"encoding/json"

	toml "github.com/pelletier/go-toml"
)

func TomlToJSON(data []byte) ([]byte, error) {
	tree, err := toml.LoadBytes(data)
	if err != nil {
		return []byte{}, err
	}
	json, err := json.Marshal(tree.ToMap())
	if err != nil {
		return []byte{}, err
	}
	return json, nil
}
