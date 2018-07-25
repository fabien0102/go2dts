package client

import (
	"time"

	"github.com/contiamo/labs/pkg/sql"

	uuid "github.com/satori/go.uuid"
)

// FunctionResponse gives a minimal description of a function in a labs server cluster, this
// is appropriate for use in list responses.
type FunctionResponse struct {
	ID               uuid.UUID           `json:"id,omitempty"`
	CreatedAt        time.Time           `json:"createdAt"`
	UpdatedAt        time.Time           `json:"updatedAt"`
	BundleID         uuid.UUID           `json:"bundleId"`
	Name             string              `json:"name"`
	Image            string              `json:"image"`
	Command          string              `json:"command"`
	Environment      sql.JSONStringMap   `json:"environment"`
	Secrets          sql.JSONStringArray `json:"secrets"`
	Schema           sql.JSONMap         `json:"schema"`
	Deployed         bool                `json:"deployed"`
	DeploymentStatus string              `json:"deploymentStatus"`
	URL              string              `json:"url"`
}

// FunctionInstanceResponse describes a single deployed function in a labs server cluster.  This should
// be used by instance detail  responses.  It may/will include additional information, such as the
// function schema, readme, and other detailed long form data.
type FunctionInstanceResponse struct {
	ID               uuid.UUID           `json:"id,omitempty"`
	CreatedAt        time.Time           `json:"createdAt"`
	UpdatedAt        time.Time           `json:"updatedAt"`
	BundleID         uuid.UUID           `json:"bundleId"`
	Name             string              `json:"name"`
	Image            string              `json:"image"`
	Command          string              `json:"command"`
	Environment      sql.JSONStringMap   `json:"environment"`
	Secrets          sql.JSONStringArray `json:"secrets"`
	Schema           sql.JSONMap         `json:"schema"`
	Deployed         bool                `json:"deployed"`
	DeploymentStatus string              `json:"deploymentStatus"`
	URL              string              `json:"url"`
}

// FunctionListResponse contains the paging and the data array for the bundle list endpoint
type FunctionListResponse struct {
	Page PageMeta            `json:"page"`
	Data []*FunctionResponse `json:"data"`
}
