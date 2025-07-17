package provider

import (
	"reflect"
	"testing"

	broadpeakio "github.com/bashou/bpkio-go-sdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFlattenSourceSlate(t *testing.T) {
	input := broadpeakio.Source{
		Id:          42,
		Name:        "Slate Sample",
		Type:        "slate",
		Url:         "http://cdn.example/slate.mp4",
		Description: "Fallback video",
	}

	expected := sourceSlateDataSourceModel{
		ID:          types.Int64Value(42),
		Name:        types.StringValue("Slate Sample"),
		Type:        types.StringValue("slate"),
		URL:         types.StringValue("http://cdn.example/slate.mp4"),
		Description: types.StringValue("Fallback video"),
	}

	got := flattenSourceSlate(input)

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("flattenSourceSlate() mismatch:\nexpected: %#v\ngot:      %#v", expected, got)
	}
}
