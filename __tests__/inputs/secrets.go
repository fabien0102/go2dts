package client

import (
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/satori/go.uuid"
)

const (
	nameValidation         = `(?i)^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	nameValidationErrorMsg = "The secret name must consist of lower case alphanumeric characters or '-'. The regex used for validation is '" + nameValidation + "'"
)

var (
	nameValidationRegExp = regexp.MustCompile(nameValidation)
)

// CreateSecretRequest defines a valid  payload required for creating new
// secret in the Lab server environment
type CreateSecretRequest struct {
	BundleID uuid.UUID `json:"bundleId"`
	Name     string    `json:"name"`
	Value    string    `json:"value"`
}

// BundleSecretResponse is the lab server response describing a bundle secret
type BundleSecretResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	BundleID  uuid.UUID `json:"bundleId" `
	Name      string    `json:"name" `
}

// SecretListResponse is the lab server list response for the secrets
// associated with a bundle
type SecretListResponse struct {
	Page PageMeta               `json:"page"`
	Data []BundleSecretResponse `json:"data"`
}

// Validate ensures that create request is valid.  This also
// implements the ozzo-validation.Validatable interface.
func (s *CreateSecretRequest) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(
			&s.Name,
			validation.Required,
			validation.Match(nameValidationRegExp).Error(nameValidationErrorMsg),
		),
		validation.Field(&s.Value, validation.Required),
		validation.Field(
			&s.BundleID,
			validation.By(nonEmptyUUIDRule),
			validation.Required,
		),
	)
}
