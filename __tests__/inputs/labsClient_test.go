package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"

	test "github.com/contiamo/labs/pkg/testutils"
)

func Test_Bundle_Sync(t *testing.T) {
	tenantID := "1"
	realmID := "2"
	bundleID := uuid.NewV4()

	t.Run("Bundle Sync Success", func(t *testing.T) {
		require := require.New(t)

		successResp := BundleSyncResponse{
			GitURL: "git@github.com:test/example.git",
			Branch: "master",
			SHA:    "abc123",
		}
		requests := []test.Request{
			{
				Method:       http.MethodPost,
				URI:          tenantID + "/api/v1/realms/" + realmID + "/bundles/" + bundleID.String() + "/sync",
				ResponseBody: successResp,
			},
		}
		ts := test.MockHTTPMuxServer(t, requests, 1)
		defer ts.Close()
		c := New(http.DefaultClient, ts.URL, tenantID, realmID)

		resp, err := c.Sync(context.Background(), bundleID)
		require.NoError(err)
		require.Equal(successResp.Branch, resp.Branch)
		require.Equal(successResp.GitURL, resp.GitURL)
		require.Equal(successResp.SHA, resp.SHA)
	})

	t.Run("Bundle Sync Auth Error", func(t *testing.T) {
		require := require.New(t)
		requests := []test.Request{
			{
				Method:             http.MethodPost,
				URI:                tenantID + "/api/v1/realms/" + realmID + "/bundles/" + bundleID.String() + "/sync",
				ResponseStatusCode: http.StatusForbidden,
			},
		}
		ts := test.MockHTTPMuxServer(t, requests, 1)
		defer ts.Close()
		c := New(http.DefaultClient, ts.URL, tenantID, realmID)

		resp, err := c.Sync(context.Background(), bundleID)
		require.Nil(resp)
		require.EqualError(err, "unauthorized request")
	})

	t.Run("Bundle Sync Server Error", func(t *testing.T) {
		require := require.New(t)
		requests := []test.Request{
			{
				Method:             http.MethodPost,
				URI:                tenantID + "/api/v1/realms/" + realmID + "/bundles/" + bundleID.String() + "/sync",
				ResponseStatusCode: http.StatusInternalServerError,
			},
		}
		ts := test.MockHTTPMuxServer(t, requests, 1)
		defer ts.Close()
		c := New(http.DefaultClient, ts.URL, tenantID, realmID)

		resp, err := c.Sync(context.Background(), bundleID)
		require.Nil(resp)
		require.EqualError(err, "unexpected status code 500")
	})
}
