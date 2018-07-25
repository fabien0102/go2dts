package client

import (
	"context"
	"regexp"
	"time"

	"github.com/contiamo/labs/pkg/git/revision"
	"github.com/contiamo/labs/pkg/sql"
	"github.com/contiamo/labs/pkg/sql/null"
	"github.com/contiamo/labs/pkg/types"
	validation "github.com/go-ozzo/ozzo-validation"

	uuid "github.com/satori/go.uuid"
)

var (
	gitURLRegex = regexp.MustCompile(`^((git|ssh|http(s)?)|(git@[\w\.]+))(:(//)?)([\w\.@\:/\-~]+)(\.git)(/)?$`)
)

// BundleClient wraps all functionality needed to interact with the labserver as an enduser
type BundleClient interface {
	Register(ctx context.Context, req *RegisterBundleRequest) (*BundleResponse, error)
	Unregister(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) (*BundleListResponse, error)
	Deploy(ctx context.Context, id uuid.UUID) (*DeployResponse, error)
	Sync(ctx context.Context, id uuid.UUID) (*BundleSyncResponse, error)
	Undeploy(ctx context.Context, id uuid.UUID) error
	StartEditSession(ctx context.Context, id uuid.UUID) (url string, err error)
	StopEditSession(ctx context.Context, id uuid.UUID, deleteVolumes bool) error
	Status(ctx context.Context, id uuid.UUID) (*BundleResponse, error)
	GetTenantID() string
	SetRealm(string)
}

// RegisterBundleRequest is the request body required to create a new bundle instance in
// the Lab Server
type RegisterBundleRequest struct {
	Name     string `json:"name"`
	GitURL   string `json:"gitUrl"`
	Branch   string `json:"branch"`
	TenantID string `json:"tenantId"`
	RealmID  string `json:"realmId"`
}

// IsNonEmpty returns a boolean indicating if any of the PATCHable fields have
// been set. This can be used to short-cicuit empty requests.
func (r *RegisterBundleRequest) IsNonEmpty() bool {
	return (r.Name != "") || (r.GitURL != "") || (r.Branch != "") || (r.TenantID != "")
}

// IsEmpty returns a boolean indicating if any of the PATCHable fields have
// been set. This can be used to short-cicuit empty requests.
func (r *RegisterBundleRequest) IsEmpty() bool {
	return !r.IsNonEmpty()
}

// Validate ensures that the bundle values are valid.  This also
// implements the ozzo-validation.Validatable interface.
func (r *RegisterBundleRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.TenantID, validation.Required),
		validation.Field(&r.Name, validation.Required, validation.Length(3, 255)),
		validation.Field(
			&r.Branch,
			validation.Required,
			validation.Length(1, 255),
			validation.By(revision.ValidationRule),
		),
		validation.Field(
			&r.GitURL,
			validation.Required,
			validation.Length(1, 255),
			validation.Match(gitURLRegex.Copy()).Error("must be a valid git url"),
		),
	)
}

// PatchBundleRequest defines a valid  payload and fields that are editable
// on the Bundle model.
type PatchBundleRequest struct {
	Name          *string             `json:"name"`
	GitURL        *string             `json:"gitUrl"`
	Branch        *string             `json:"branch"`
	CoverImageURL *string             `json:"coverImageURL" `
	Tags          sql.JSONStringArray `json:"tags"`
}

// IsNonEmpty returns a boolean indicating if any of the PATCHable fields have
// been set. This can be used to short-cicuit empty requests.
func (r *PatchBundleRequest) IsNonEmpty() bool {
	return (r.Name != nil) ||
		(r.GitURL != nil) ||
		(r.Branch != nil) ||
		(r.CoverImageURL != nil) ||
		(r.Tags != nil)
}

// IsEmpty returns a boolean indicating if any of the PATCHable fields have
// been set. This can be used to short-cicuit empty requests.
func (r *PatchBundleRequest) IsEmpty() bool {
	return !r.IsNonEmpty()
}

// BundleResponse represents all of the readable fields from the Lab server bundle representation
type BundleResponse struct {
	ID              uuid.UUID           `json:"id"`
	CreatedAt       time.Time           `json:"createdAt"`
	UpdatedAt       time.Time           `json:"updatedAt"`
	SyncedAt        null.Time           `json:"syncedAt"`
	DeployedAt      null.Time           `json:"deployedAt"`
	DeployedSHA     string              `json:"deployedSHA"`
	Name            string              `json:"name"`
	GitURL          string              `json:"gitUrl"`
	GitSHA          string              `json:"gitSHA"`
	Branch          string              `json:"branch"`
	BundleConfig    types.Bundle        `json:"config"`
	PublicKey       string              `json:"publicKey"`
	CoverImageURL   string              `json:"coverImageURL" `
	Tags            sql.JSONStringArray `json:"tags"`
	TenantID        string              `json:"tenantId"`
	RealmID         string              `json:"realmId"`
	ActiveDeployID  null.UUID           `json:"activeDeployID"`
	TopContributors []Contributor       `json:"topContributors"`
	Description     string              `json:"description"`
}

// BundleListResponse contains the paging and the data array for the bundle list endpoint
type BundleListResponse struct {
	Page PageMeta          `json:"page"`
	Data []*BundleResponse `json:"data"`
}

// DeployResponse is returned during a bundle deploy, it describes the Bundle and the Functions deployed
type DeployResponse struct {
	Bundle    BundleResponse      `json:"bundle"`
	Functions []*FunctionResponse `json:"functions"`
}

// DeployLog is an immutable log of events that occur during Bundle deploys.
type DeployLog struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	BundleID  uuid.UUID `json:"bundleId"`
	DeployID  uuid.UUID `json:"deployId"`
	GitURL    string    `json:"gitUrl"`
	SHA       string    `json:"sha"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
}

// DeployLogResponse contains the paging and the data array for the DeployLog list response
type DeployLogResponse struct {
	Page PageMeta    `json:"page"`
	Data []DeployLog `json:"data"`
}

// ContributorListResponse contains the paging and data array of contributors to a bundle
type ContributorListResponse struct {
	Page PageMeta       `json:"page"`
	Data []*Contributor `json:"data"`
}

// BundleEditStartResponse is the url returned from the labs server indicating where
// the edit server is located.
type BundleEditStartResponse struct {
	URL string `json:"url"`
}

// BundleSyncResponse contains the git meta data for the bundle after we finish syncing with the
// remote git repository
type BundleSyncResponse struct {
	SHA     string `json:"sha"`
	Branch  string `json:"branch"`
	GitURL  string `json:"gitURL"`
	Updated bool   `json:"updated"`
}
