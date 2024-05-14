package converter

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/rynowak/celery/pkg/documents"
	"github.com/rynowak/celery/pkg/internal"
	"github.com/stretchr/testify/require"
)

func Test_Gateway(t *testing.T) {
	provider := internal.Must(documents.UnmarshalProvider(internal.Must(os.ReadFile("../testdata/applications.core.yaml"))))
	input := internal.MustUnmarshalAny(os.ReadFile("../testdata/applications.core.gateways.in.2023-01-01.basic.json"))

	converted := map[string]any{}
	err := Convert(input, &converted, provider, "Applications.Core/gateways", "2023-01-01")
	require.NoError(t, err)

	expected := internal.MustUnmarshalAny(os.ReadFile("../testdata/applications.core.gateways.out.v0.basic.json"))
	require.JSONEq(t, string(internal.Must(json.Marshal(expected))), string(internal.Must(json.Marshal(converted))))
}
