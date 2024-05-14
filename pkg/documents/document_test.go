package documents

import (
	"os"
	"testing"

	"github.com/rynowak/celery/pkg/internal"
	"github.com/stretchr/testify/require"
)

func Test_Applications_Core(t *testing.T) {
	b, err := os.ReadFile("../testdata/applications.core.yaml")
	require.NoError(t, err)

	provider, err := UnmarshalProvider(b)
	require.NoError(t, err)

	expected := &Provider{
		Namespace: "Applications.Core",
		Resources: map[string]*Resource{
			"gateways": {
				Datamodel: []*Datamodel{
					{
						Schema: &Schema{
							Type: SchemaTypeObject,
							Properties: map[string]*Schema{
								"hostname": {
									Type:     SchemaTypeObject,
									Optional: true,
									Properties: map[string]*Schema{
										"prefix": {
											Type:     SchemaTypeString,
											Optional: true,
										},
										"fullyQualifiedHostname": {
											Type:     SchemaTypeString,
											Optional: true,
										},
									},
								},
								"internal": {
									Type:    SchemaTypeBoolean,
									Default: internal.ToPtr("false"),
								},
								"routes": {
									Type: SchemaTypeArray,
									Element: &Schema{
										Type: SchemaTypeObject,
										Properties: map[string]*Schema{
											"destination": {
												Type:     SchemaTypeString,
												Optional: true,
											},
											"path": {
												Type:     SchemaTypeString,
												Optional: true,
											},
											"replacePrefix": {
												Type:     SchemaTypeString,
												Optional: true,
											},
										},
									},
								},
								"tls": {
									Type:     SchemaTypeObject,
									Optional: true,
									Properties: map[string]*Schema{
										"certificateFrom": {
											Type:     SchemaTypeString,
											Optional: true,
										},
										"minimumProtocolVersion": {
											Type: SchemaTypeString,
											Enum: []string{"1.2", "1.3"},
										},
										"sslPassthrough": {
											Type:    SchemaTypeBoolean,
											Default: internal.ToPtr("false"),
										},
									},
								},
								"url": {
									Type:     SchemaTypeString,
									Optional: true,
								},
							},
						},
					},
				},
				APIVersions: map[string]*APIVersion{
					"2023-01-01": {
						Schema: &Schema{
							Type: SchemaTypeObject,
							Properties: map[string]*Schema{
								"internal": {
									Type: SchemaTypeBoolean,
								},
							},
						},
					},
				},
			},
		},
	}

	require.Equal(t, expected, provider)
}
