package internal

import "encoding/json"

func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}

func MustUnmarshalAny(b []byte, err error) any {
	if err != nil {
		panic(err)
	}

	var a any
	err = json.Unmarshal(b, &a)
	if err != nil {
		panic(err)
	}

	return &a
}

func ToPtr[T any](value T) *T {
	copy := value
	return &copy
}
