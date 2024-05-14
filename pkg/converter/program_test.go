package converter

import (
	"os"
	"testing"

	"github.com/rynowak/celery/pkg/documents"
	"github.com/rynowak/celery/pkg/internal"
	"github.com/stretchr/testify/require"
)

func Test_ApplicationsCore_Gateways_2023_01_01(t *testing.T) {
	provider := internal.Must(documents.UnmarshalProvider(internal.Must(os.ReadFile("../testdata/applications.core.yaml"))))
	gateways := provider.Resources["gateways"]

	text, err := GenerateInputConversionProgram(gateways.APIVersions["2023-01-01"], gateways.Datamodel[len(gateways.Datamodel)-1])
	require.NoError(t, err)

	expected := string(internal.Must(os.ReadFile("../testdata/applications.core.gateways.convert-in.2023-01-01.cel")))
	require.Equal(t, expected, *text)
}
