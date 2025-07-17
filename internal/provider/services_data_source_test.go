package provider

import (
	"context"
	"testing"

	broadpeakio "github.com/bashou/bpkio-go-sdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func mustList(ctx context.Context, vals []string) types.List {
	l, diags := types.ListValueFrom(ctx, types.StringType, vals)
	if diags.HasError() {
		panic(diags.Errors()[0].Detail())
	}
	return l
}

func TestFlattenService(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name   string
		input  broadpeakio.ServiceOutput
		expect serviceDataSourceModel
	}{
		{
			name: "full",
			input: broadpeakio.ServiceOutput{
				Id:              42,
				Name:            "svc",
				Type:            "ad-insertion",
				Url:             "http://service",
				CreationDate:    "2023-01-01T00:00:00Z",
				UpdateDate:      "2023-01-02T00:00:00Z",
				State:           "enabled",
				EnvironmentTags: []string{"dev", "live"},
			},
			expect: serviceDataSourceModel{
				ID:           types.Int64Value(42),
				Name:         types.StringValue("svc"),
				Type:         types.StringValue("ad-insertion"),
				URL:          types.StringValue("http://service"),
				CreationDate: types.StringValue("2023-01-01T00:00:00Z"),
				UpdateDate:   types.StringValue("2023-01-02T00:00:00Z"),
				State:        types.StringValue("enabled"),
				Tags:         mustList(ctx, []string{"dev", "live"}),
			},
		},
		{
			name: "no tags",
			input: broadpeakio.ServiceOutput{
				Id:              43,
				Name:            "no-tags",
				Type:            "virtual-channel",
				Url:             "http://none",
				CreationDate:    "",
				UpdateDate:      "",
				State:           "disabled",
				EnvironmentTags: nil,
			},
			expect: serviceDataSourceModel{
				ID:           types.Int64Value(43),
				Name:         types.StringValue("no-tags"),
				Type:         types.StringValue("virtual-channel"),
				URL:          types.StringValue("http://none"),
				CreationDate: types.StringValue(""),
				UpdateDate:   types.StringValue(""),
				State:        types.StringValue("disabled"),
				Tags:         mustList(ctx, nil),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := flattenService(tt.input, ctx)
			require.NoError(t, err)
			require.Equal(t, tt.expect, got)
		})
	}
}
