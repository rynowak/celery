package resources

import (
	"errors"
	"strings"
)

func ParseFullyQualifiedType(tt string) (string, []string, error) {
	parts := strings.Split(tt, "/")
	if len(parts) == 0 {
		return "", nil, errors.New("resource types should have the format '<namespace>/<type>`")
	}

	return parts[0], parts[1:], nil
}
