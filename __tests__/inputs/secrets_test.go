package client

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/satori/go.uuid"
)

func Test_Secret_Validation(t *testing.T) {
	bundleID := uuid.NewV4()

	nameErroMsg := "name: " + nameValidationErrorMsg + "."

	happyPath := []struct {
		name        string
		secretName  string
		secretValue string
		bundleID    uuid.UUID
		errorMsg    string
	}{
		{
			name:        "Valid Secret is valid",
			secretName:  "valid-name-HERE",
			secretValue: "this is a value",
			bundleID:    bundleID,
		},
		{
			name:        "Empty name is invalid",
			secretName:  "",
			secretValue: "this is a value",
			bundleID:    bundleID,
			errorMsg:    "name: cannot be blank.",
		},
		{
			name:        "Empty value is invalid",
			secretName:  "test",
			secretValue: "",
			bundleID:    bundleID,
			errorMsg:    "value: cannot be blank.",
		},
		{
			name:        "Empty bundlid is invalid",
			secretName:  "test",
			secretValue: "test",
			errorMsg:    "bundleId: cannot be blank.",
		},
		{
			name:        "Special characters in the name are invalid",
			secretName:  `name-cannot-contain-$%^&*()#!@_.,<>?/\|~"'{}[]`,
			secretValue: "this is a value",
			bundleID:    bundleID,
			errorMsg:    nameErroMsg,
		},
	}

	for _, s := range happyPath {
		t.Run(s.name, func(t *testing.T) {
			req := CreateSecretRequest{
				Name:     s.secretName,
				Value:    s.secretValue,
				BundleID: s.bundleID,
			}

			err := req.Validate()
			if s.errorMsg != "" {
				require.EqualError(t, err, s.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
