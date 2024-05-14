package converter

import (
	"fmt"
	"reflect"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/rynowak/celery/pkg/documents"
)

var _ types.Provider = (*TypeProvider)(nil)

type TypeProvider struct {
	types map[string]entry
}

type entry struct {
	cel    *types.Type
	schema *documents.Schema
}

func NewTypeProvider(roots map[string]*documents.Schema) *TypeProvider {
	tp := &TypeProvider{types: map[string]entry{}}
	tp.initialize("", roots)
	return tp
}

func (tp *TypeProvider) initialize(prefix string, roots map[string]*documents.Schema) {
	for name, schema := range roots {
		fqn := name
		if len(prefix) > 0 {
			fqn = prefix + "." + name
		}

		if schema.Type == documents.SchemaTypeObject {

			structType := types.NewObjectType(fqn)
			if schema.Optional {
				structType = types.NewOptionalType(structType)
			}

			tp.types[fqn] = entry{cel: types.NewTypeTypeWithParam(structType), schema: schema}
			tp.initialize(fqn, schema.Properties)
		} else if schema.Type.IsScalar() {
			// Nothing to do
		} else {
			panic("unsupported type: " + schema.Type)
		}
	}
}

// EnumValue implements types.Provider.
func (tp *TypeProvider) EnumValue(enumName string) ref.Val {
	// No enums
	return types.NewErr("unknown enum name '%s'", enumName)
}

// FindIdent implements types.Provider.
func (tp *TypeProvider) FindIdent(identName string) (ref.Val, bool) {
	// No global variables
	return nil, false
}

// FindStructFieldNames implements types.Provider.
func (tp *TypeProvider) FindStructFieldNames(structType string) ([]string, bool) {
	panic("unimplemented")
}

// FindStructFieldType implements types.Provider.
func (tp *TypeProvider) FindStructFieldType(structType string, fieldName string) (*types.FieldType, bool) {
	structTypeEntry, ok := tp.types[structType]
	if !ok {
		return nil, false
	}

	property, ok := structTypeEntry.schema.Properties[fieldName]
	if !ok {
		return nil, false
	}

	return &types.FieldType{
		Type:    tp.findStructFieldType(structType, fieldName, property),
		IsSet:   isSetMapField(fieldName),
		GetFrom: readMapField(fieldName),
	}, true
}

func (tp *TypeProvider) findStructFieldType(structType string, fieldName string, property *documents.Schema) *types.Type {
	if property.Type == documents.SchemaTypeObject {
		st, ok := tp.FindStructType(structType + "." + fieldName)
		if !ok {
			panic("missing type for stuct property: " + structType + "." + fieldName)
		}

		return st
	} else if property.Type == documents.SchemaTypeBoolean && property.Optional {
		return types.NewOptionalType(types.BoolType)
	} else if property.Type == documents.SchemaTypeBoolean {
		return types.BoolType
	} else if property.Type == documents.SchemaTypeNumber && property.Optional {
		return types.NewOptionalType(types.DoubleType)
	} else if property.Type == documents.SchemaTypeNumber {
		return types.DoubleType
	} else if property.Type == documents.SchemaTypeInteger && property.Optional {
		return types.NewOptionalType(types.IntType)
	} else if property.Type == documents.SchemaTypeInteger {
		return types.IntType
	} else if property.Type == documents.SchemaTypeString && property.Optional {
		return types.NewOptionalType(types.StringType)
	} else if property.Type == documents.SchemaTypeString {
		return types.StringType
	}

	panic("unsupported type: " + property.Type)
}

// FindStructType implements types.Provider.
func (tp *TypeProvider) FindStructType(structType string) (*types.Type, bool) {
	structTypeEntry, ok := tp.types[structType]
	if !ok {
		return nil, false
	}

	return structTypeEntry.cel, true
}

// NewValue implements types.Provider.
func (tp *TypeProvider) NewValue(structType string, fields map[string]ref.Val) ref.Val {
	structTypeEntry, ok := tp.types[structType]
	if !ok {
		return types.NewErr("unknown struct type: %s", structType)
	}

	data := map[string]any{}
	for key, value := range fields {
		_, ok := structTypeEntry.schema.Properties[key]
		if !ok {
			return types.NewErr("unknown struct field: %s.%s", structType, key)
		}

		converted, err := value.ConvertToNative(reflect.TypeFor[any]())
		if err != nil {
			return types.NewErr("failed to convert struct field: %s.%s", structType, key)
		}

		data[key] = converted
	}

	return types.NewDynamicMap(&MapAdapter{}, data)
}

func isSetMapField(name string) func(target any) bool {
	return func(target any) bool {
		tv := reflect.ValueOf(target)

		for tv.Kind() == reflect.Pointer || tv.Kind() == reflect.Interface {
			tv = tv.Elem()
		}

		if tv.Kind() != reflect.Map {
			return false
		}

		for _, key := range tv.MapKeys() {
			if key == reflect.ValueOf(name) {
				return true
			}
		}

		return false
	}
}

func readMapField(name string) func(target any) (any, error) {
	return func(target any) (any, error) {
		tv := reflect.ValueOf(target)

		for tv.Kind() == reflect.Pointer || tv.Kind() == reflect.Interface {
			tv = tv.Elem()
		}

		if tv.Kind() != reflect.Map {
			return nil, fmt.Errorf("expected map, got: %T", target)
		}

		mv := tv.MapIndex(reflect.ValueOf(name))
		return mv.Interface(), nil
	}
}

type MapAdapter struct {
}

func (m *MapAdapter) NativeToValue(value any) ref.Val {
	panic("unimplemented")
}
