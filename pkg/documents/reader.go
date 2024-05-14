package documents

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

func UnmarshalProvider(b []byte) (*Provider, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(b))
	decoder.KnownFields(true)

	provider := Provider{}
	err := decoder.Decode(&provider)
	if err != nil {
		return nil, err
	}

	return &provider, nil
}
