// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"reflect"
	"testing"

	broadpeakio "github.com/bashou/bpkio-go-sdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFlattenSources(t *testing.T) {
	testCases := []struct {
		name       string
		filterType *string
		input      []broadpeakio.Source
		expected   []sourcesModel
	}{
		{
			name:       "no filter, multiple types",
			filterType: nil,
			input: []broadpeakio.Source{
				{Id: 1, Name: "Live1", Type: "live", Url: "http://live1"},
				{Id: 2, Name: "Asset1", Type: "asset", Url: "http://asset1"},
			},
			expected: []sourcesModel{
				{ID: types.Int64Value(1), Name: types.StringValue("Live1"), Type: types.StringValue("live"), URL: types.StringValue("http://live1")},
				{ID: types.Int64Value(2), Name: types.StringValue("Asset1"), Type: types.StringValue("asset"), URL: types.StringValue("http://asset1")},
			},
		},
		{
			name:       "filter live only",
			filterType: ptr("live"),
			input: []broadpeakio.Source{
				{Id: 1, Name: "Live1", Type: "live", Url: "http://live1"},
				{Id: 2, Name: "Asset1", Type: "asset", Url: "http://asset1"},
			},
			expected: []sourcesModel{
				{ID: types.Int64Value(1), Name: types.StringValue("Live1"), Type: types.StringValue("live"), URL: types.StringValue("http://live1")},
			},
		},
		{
			name:       "filter matches nothing",
			filterType: ptr("ad-server"),
			input: []broadpeakio.Source{
				{Id: 1, Name: "Live1", Type: "live", Url: "http://live1"},
			},
			expected: []sourcesModel{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := flattenSources(tc.input, tc.filterType)
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("expected:\n%#v\ngot:\n%#v", tc.expected, got)
			}
		})
	}
}

// tiny helper
func ptr[T any](v T) *T {
	return &v
}
