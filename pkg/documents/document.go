package documents

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Provider struct {
	Namespace string               `yaml:"namespace"`
	Resources map[string]*Resource `yaml:"resources"`
}

type Resource struct {
	Datamodel   []*Datamodel           `yaml:"datamodel"`
	APIVersions map[string]*APIVersion `yaml:"apiVersions"`
}

type Datamodel struct {
	Schema *Schema `yaml:"schema"`
}

type APIVersion struct {
	Capabilities []string `yaml:"capabilities"`
	Schema       *Schema  `yaml:"schema"`
}

type Schema struct {
	Type       SchemaType         `yaml:"type"`
	Optional   bool               `yaml:"optional,omitempty"`
	Default    *string            `yaml:"default,omitempty"`
	Enum       []string           `yaml:"enum,omitempty"`
	Element    *Schema            `yaml:"element,omitempty"`
	Properties map[string]*Schema `yaml:"properties,omitempty"`
}

type SchemaType string

const (
	SchemaTypeArray   SchemaType = "array"
	SchemaTypeMap     SchemaType = "map"
	SchemaTypeObject  SchemaType = "object"
	SchemaTypeString  SchemaType = "string"
	SchemaTypeNumber  SchemaType = "number"
	SchemaTypeBoolean SchemaType = "boolean"
	SchemaTypeInteger SchemaType = "integer"
)

var knownSchemaTypes []string = []string{
	string(SchemaTypeArray),
	string(SchemaTypeMap),
	string(SchemaTypeObject),
	string(SchemaTypeString),
	string(SchemaTypeNumber),
	string(SchemaTypeBoolean),
	string(SchemaTypeInteger),
}

func (s SchemaType) IsScalar() bool {
	return s == SchemaTypeBoolean || s == SchemaTypeInteger || s == SchemaTypeNumber || s == SchemaTypeString
}

func (s *SchemaType) UnmarshalYAML(value *yaml.Node) error {
	ss := ""
	err := value.Decode(&ss)
	if err != nil {
		return err
	}

	found := false
	for i := range knownSchemaTypes {
		if knownSchemaTypes[i] == ss {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("unknown schema type: %q. Valid values: %s", ss, strings.Join(knownSchemaTypes, ", "))
	}

	*s = SchemaType(ss)
	return nil
}
