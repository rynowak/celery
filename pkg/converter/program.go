package converter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rynowak/celery/pkg/documents"
	"golang.org/x/exp/maps"
)

func GenerateInputConversionProgram(source *documents.APIVersion, destination *documents.Datamodel) (*string, error) {
	text, errs := convertValue("$", source.Schema, destination.Schema, destination.Schema)
	if len(errs) > 0 {
		lines := []string{}
		for i := range errs {
			lines = append(lines, fmt.Sprintf(" - %s", errs[i].Error()))
		}
		return nil, fmt.Errorf("failed to create conversion:\n%s", strings.Join(lines, "\n"))
	}

	return text, nil
}

func convertValue(path string, source *documents.Schema, destination *documents.Schema, root *documents.Schema) (*string, []error) {
	// Destination can never be nil.
	if source == nil && destination.Default == nil {
		return nil, nil
	} else if source == nil {
		return destination.Default, nil
	}

	if source.Type == documents.SchemaTypeObject && destination.Type == documents.SchemaTypeObject {
		return convertObjects(path, source, destination, root)
	} else if source.Type == documents.SchemaTypeArray && destination.Type == documents.SchemaTypeArray {
		return convertArrays(path, source, destination)
	} else if source.Type == documents.SchemaTypeMap && destination.Type == documents.SchemaTypeMap {
		return convertMaps(path, source, destination)
	} else if source.Type.IsScalar() && destination.Type.IsScalar() {
		return convertScalars(path, source, destination)
	}

	return nil, []error{fmt.Errorf("unsupported conversion at '%s' %+v -> %+v", path, source.Type, destination.Type)}
}

func convertObjects(path string, source *documents.Schema, destination *documents.Schema, root *documents.Schema) (*string, []error) {
	keys := maps.Keys(destination.Properties)
	sort.Strings(keys)

	errs := []error{}
	properties := map[string]string{}
	for i := range keys {
		key := keys[i]
		expr, err := convertValue(path+"."+key, source.Properties[key], destination.Properties[key], root)
		if err != nil {
			errs = append(errs, err...)
		} else if expr != nil {
			properties[key] = *expr
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}

	lines := []string{}

	for i := range keys {
		key := keys[i]
		value, hasValue := properties[key]

		if hasValue {
			lines = append(lines, fmt.Sprintf("%s: %s", key, value))
		}
	}

	if len(lines) == 0 && destination.Optional {
		return nil, nil
	}

	typeName := walkpath(path, root, "DataModel", func(part string, previous string) string {
		return previous + "." + part
	})

	expr := fmt.Sprintf("%s{\n%s\n}", typeName, strings.Join(lines, "\n"))
	return &expr, nil
}

func convertArrays(path string, source *documents.Schema, destination *documents.Schema) (*string, []error) {
	return nil, []error{fmt.Errorf("unsupported conversion at '%s' %+v -> %+v", path, source.Type, destination.Type)}
}

func convertMaps(path string, source *documents.Schema, destination *documents.Schema) (*string, []error) {
	return nil, []error{fmt.Errorf("unsupported conversion at '%s' %+v -> %+v", path, source.Type, destination.Type)}
}

func convertScalars(path string, source *documents.Schema, destination *documents.Schema) (*string, []error) {
	expr := strings.Replace(path, "$", "input", 1)
	if source.Optional && destination.Default != nil {
		expr = expr + " || " + *destination.Default
	}

	return &expr, nil
}
