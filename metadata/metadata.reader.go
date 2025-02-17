package metadata

import (
	"os"

	"gopkg.in/yaml.v2"
)

func ReadMetadata(filePath string) (map[string]interface{}, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var metadata map[string]interface{}
	err = yaml.Unmarshal(file, &metadata)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}
