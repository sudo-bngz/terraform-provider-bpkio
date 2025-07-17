package provider

import (
	"context"
	"testing"

	broadpeakio "github.com/bashou/bpkio-go-sdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestFlattenAdInsertionOutput_basic(t *testing.T) {
	ctx := context.Background()
	in := broadpeakio.AdInsertionOutput{
		Id:                  1,
		Name:                "AI Service",
		Type:                "ad-insertion",
		Url:                 "http://service",
		State:               "enabled",
		Tags:                []string{"ad", "insertion"},
		CreationDate:        "2024-01-01T00:00:00Z",
		UpdateDate:          "2024-02-01T00:00:00Z",
		EnableAdTranscoding: true,
		ServerSideAdTracking: broadpeakio.ServerSideAdTracking{
			Enable:                          true,
			CheckAdMediaSegmentAvailability: false,
		},
		Source: broadpeakio.Source{
			Id:   2,
			Name: "SRC",
			Type: "live",
			Url:  "http://src",
			Origin: broadpeakio.Origin{
				CustomHeaders: []broadpeakio.CustomHeader{
					{Name: "X-Test", Value: "42"},
				},
			},
		},
		TranscodingProfile: broadpeakio.TranscodingProfile{
			Id:         3,
			Name:       "Default",
			InternalId: "int-3",
			Content:    "{\"foo\":\"bar\"}",
		},
		AdvancedOptions: broadpeakio.AdvancedOptions{
			AuthorizationHeader: broadpeakio.AuthorizationHeader{
				Name:  "X-Auth",
				Value: "token",
			},
		},
		LiveAdPreRoll: broadpeakio.LiveAdPreRollOutput{
			AdServer: broadpeakio.AdServer{
				Id:   4,
				Name: "adserver",
				Type: "ad-server",
				Url:  "http://adserver",
				QueryParameters: []broadpeakio.QueryParam{
					{Name: "foo", Value: "bar", Type: "static"},
				},
			},
			MaxDuration: 10,
			Offset:      2,
		},
		LiveAdReplacement: broadpeakio.LiveAdReplacementOutput{
			AdServer: broadpeakio.AdServer{
				Id:   5,
				Name: "adserver2",
				Type: "ad-server",
				Url:  "http://adserver2",
				QueryParameters: []broadpeakio.QueryParam{
					{Name: "baz", Value: "qux", Type: "dynamic"},
				},
			},
			GapFiller: broadpeakio.GapFiller{
				Id:   6,
				Name: "gapfiller",
				Type: "slate",
				Url:  "http://gapfiller",
			},
			SpotAware: broadpeakio.SpotAware{
				Mode: "spot_to_live",
			},
		},
	}

	out, err := flattenAdInsertionOutput(in, ctx)
	require.NoError(t, err)
	require.Equal(t, types.Int64Value(1), out.ID)
	require.Equal(t, types.StringValue("AI Service"), out.Name)
	require.Equal(t, types.StringValue("ad-insertion"), out.Type)
	require.Equal(t, types.StringValue("enabled"), out.State)
	require.Equal(t, types.BoolValue(true), out.EnableAdTranscoding)
	require.NotNil(t, out.ServerSideAdTracking)
	require.Equal(t, types.BoolValue(true), out.ServerSideAdTracking.Enable)
	require.NotNil(t, out.Source)
	require.Equal(t, types.StringValue("SRC"), out.Source.Name)
	require.NotNil(t, out.LiveAdPreRoll)
	require.Equal(t, types.StringValue("adserver"), out.LiveAdPreRoll.AdServer.Name)
	require.NotNil(t, out.LiveAdReplacement)
	require.Equal(t, types.StringValue("adserver2"), out.LiveAdReplacement.AdServer.Name)
	require.Equal(t, types.StringValue("gapfiller"), out.LiveAdReplacement.GapFiller.Name)
	require.Equal(t, types.StringValue("spot_to_live"), out.LiveAdReplacement.SpotAware.Mode)
	require.Len(t, out.Tags.Elements(), 2)
}
