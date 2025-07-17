package provider

import (
	"fmt"
	"reflect"
	"testing"

	broadpeakio "github.com/bashou/bpkio-go-sdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFlattenTranscodingProfiles(t *testing.T) {
	cases := []struct {
		input    []broadpeakio.TranscodingProfile
		expected []map[string]string
	}{
		{
			input: []broadpeakio.TranscodingProfile{
				{
					Id:         1,
					Name:       "Test Profile",
					Content:    `{"bitrate": "5000k"}`,
					InternalId: "profile-abc",
				},
			},
			expected: []map[string]string{
				{
					"id":          "1",
					"name":        "Test Profile",
					"content":     `{"bitrate": "5000k"}`,
					"internal_id": "profile-abc",
				},
			},
		},
	}

	for _, c := range cases {
		values, _, err := FlattenTranscodingProfiles(c.input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var got []map[string]string
		for _, v := range values {
			o := v.(types.Object)
			m := map[string]string{}
			for k, v := range o.Attributes() {
				switch val := v.(type) {
				case types.String:
					m[k] = val.ValueString()
				case types.Int64:
					m[k] = fmt.Sprintf("%d", val.ValueInt64())
				default:
					t.Fatalf("unexpected type for %s: %T", k, v)
				}
			}
			got = append(got, m)
		}

		if !reflect.DeepEqual(got, c.expected) {
			t.Errorf("expected:\n%#v\ngot:\n%#v", c.expected, got)
		}
	}
}
