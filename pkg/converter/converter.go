package converter

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/rynowak/celery/pkg/documents"
	"github.com/rynowak/celery/pkg/resources"
)

func Convert(input any, datamodel any, provider *documents.Provider, tt string, apiVersion string) error {
	namespace, resourceTypes, err := resources.ParseFullyQualifiedType(tt)
	if err != nil {
		return fmt.Errorf("error parsing type: %w", err)
	}

	if provider.Namespace != namespace {
		return fmt.Errorf("expected namespace '%s', got: %s", namespace, provider.Namespace)
	}

	resourceType, ok := provider.Resources[strings.Join(resourceTypes, "/")]
	if !ok {
		return fmt.Errorf("could not find type: %q", tt)
	}

	api, ok := resourceType.APIVersions[apiVersion]
	if !ok {
		return fmt.Errorf("could not find API version: %q", apiVersion)
	}

	if len(resourceType.Datamodel) == 0 {
		return fmt.Errorf("no datamodel found for type: %q", tt)
	}

	dm := resourceType.Datamodel[len(resourceType.Datamodel)-1]

	env, err := cel.NewEnv(
		cel.CustomTypeProvider(NewTypeProvider(map[string]*documents.Schema{
			"API":       api.Schema,
			"DataModel": dm.Schema,
		})),
		cel.Variable("input", cel.ObjectType("API")),
	)
	if err != nil {
		return fmt.Errorf("CEL infrastructure error: %w", err)
	}

	text, err := GenerateInputConversionProgram(api, dm)
	if err != nil {
		return err
	}

	ast, issues := env.Compile(*text)
	if issues != nil && issues.Err() != nil {
		return fmt.Errorf("CEL infrastructure error: %w", issues.Err())
	}

	prg, err := env.Program(ast)
	if err != nil {
		return fmt.Errorf("CEL infrastructure error: %w", err)
	}

	value, _, err := prg.Eval(map[string]any{
		"input": input,
	})
	if err != nil {
		return err
	}

	converted, err := value.ConvertToNative(reflect.TypeFor[any]())
	if err != nil {
		return err
	}

	b, err := json.Marshal(converted)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &datamodel)
	if err != nil {
		return err
	}

	return nil
}

func walkpath[T any](path string, root *documents.Schema, inital T, step func(part string, previous T) T) T {
	current := inital

	parts := strings.Split(path, ".")
	for _, part := range parts[1:] {
		current = step(part, current)
	}

	return current
}
