package provider

import (
	"testing"

	broadpeakio "github.com/bashou/bpkio-go-sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestFlattenSourceAdServer(t *testing.T) {
	paramObjType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":  types.StringType,
			"name":  types.StringType,
			"value": types.StringType,
		},
	}

	tests := []struct {
		name     string
		input    broadpeakio.AdServer
		expected sourceAdServerDataSourceModel
	}{
		{
			name: "basic",
			input: broadpeakio.AdServer{
				Id:   1,
				Name: "AdServer One",
				Url:  "http://ads.example/ad",
				Type: "ad-server",
				QueryParameters: []broadpeakio.QueryParam{
					{Type: "static", Name: "ad_type", Value: "pre-roll"},
				},
			},
			expected: sourceAdServerDataSourceModel{
				ID:   types.Int64Value(1),
				Name: types.StringValue("AdServer One"),
				Type: types.StringValue("ad-server"),
				URL:  types.StringValue("http://ads.example/ad"),
				QueryParameters: types.ListValueMust(paramObjType, []attr.Value{
					types.ObjectValueMust(paramObjType.AttrTypes, map[string]attr.Value{
						"type":  types.StringValue("static"),
						"name":  types.StringValue("ad_type"),
						"value": types.StringValue("pre-roll"),
					}),
				}),
			},
		},
		{
			name: "empty_parameters",
			input: broadpeakio.AdServer{
				Id:              2,
				Name:            "NoParams",
				Url:             "http://ads.example/none",
				Type:            "ad-server",
				QueryParameters: nil,
			},
			expected: sourceAdServerDataSourceModel{
				ID:              types.Int64Value(2),
				Name:            types.StringValue("NoParams"),
				Type:            types.StringValue("ad-server"),
				URL:             types.StringValue("http://ads.example/none"),
				QueryParameters: types.ListNull(paramObjType),
			},
		},
		{
			name: "multiple_parameters",
			input: broadpeakio.AdServer{
				Id:   3,
				Name: "Multi",
				Url:  "http://ads.example/multi",
				Type: "ad-server",
				QueryParameters: []broadpeakio.QueryParam{
					{Type: "dynamic", Name: "id", Value: "1234"},
					{Type: "static", Name: "placement", Value: "mid-roll"},
				},
			},
			expected: sourceAdServerDataSourceModel{
				ID:   types.Int64Value(3),
				Name: types.StringValue("Multi"),
				Type: types.StringValue("ad-server"),
				URL:  types.StringValue("http://ads.example/multi"),
				QueryParameters: types.ListValueMust(paramObjType, []attr.Value{
					types.ObjectValueMust(paramObjType.AttrTypes, map[string]attr.Value{
						"type":  types.StringValue("dynamic"),
						"name":  types.StringValue("id"),
						"value": types.StringValue("1234"),
					}),
					types.ObjectValueMust(paramObjType.AttrTypes, map[string]attr.Value{
						"type":  types.StringValue("static"),
						"name":  types.StringValue("placement"),
						"value": types.StringValue("mid-roll"),
					}),
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := flattenSourceAdServer(tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.expected, actual)
		})
	}
}
