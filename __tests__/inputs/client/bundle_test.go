package client

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_BundlePatch_RequestIsEmpty(t *testing.T) {
	happyPath := []struct {
		Name    string
		Request string
	}{
		{Name: "if empty json object", Request: "{}"},
		{Name: "if json object with no matching fields", Request: `{"doesnotexist": "ab"}`},
		{Name: "if empty json", Request: ``},
	}

	for _, s := range happyPath {
		t.Run(s.Name, func(t *testing.T) {
			u := &PatchBundleRequest{}
			json.NewDecoder(strings.NewReader(s.Request)).Decode(u)
			require.True(t, u.IsEmpty())
			require.False(t, u.IsNonEmpty())
		})
	}
}

func Test_BundlePatch_RequestIsNotEmpty(t *testing.T) {

	happyPath := []struct {
		Name    string
		Request string
	}{
		{Name: "if name field is not empty", Request: `{"name": "ab"}`},
		{Name: "if gitURL field is not empty", Request: `{"gitURL": "ab"}`},
		{Name: "if branch field is not empty", Request: `{"branch": "ab"}`},
		{Name: "if coverImageURL field is not empty", Request: `{"coverImageURL": "ab"}`},
		{Name: "if tags field is not empty", Request: `{"tags": ["test", "foo"]}`},
	}

	for _, s := range happyPath {
		t.Run(s.Name, func(t *testing.T) {
			u := &PatchBundleRequest{}
			json.NewDecoder(strings.NewReader(s.Request)).Decode(u)
			require.False(t, u.IsEmpty())
			require.True(t, u.IsNonEmpty())
		})
	}
}
