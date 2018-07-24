package client

import (
	"time"

	"github.com/contiamo/labs/pkg/constants"
)

// EditorStageStatus is an enum of names for to the editor status stage
type EditorStageStatus string

// EditorState is an enum of name for the possible editor status states
type EditorState string

const (
	// EditorStageTodo indicates the the current editor status stage is still todo
	EditorStageTodo EditorStageStatus = "todo"
	// EditorStageRunning indicates the the current editor status stage is in progress
	EditorStageRunning = "running"
	// EditorStageDone indicates the the current editor status stage is complete and successful
	EditorStageDone = "done"
	// EditorStageFailed indicates the the current editor status stage is complete and failed
	EditorStageFailed = "failed"

	// EditorStateUnknown is a fall back EditorStatus.Status
	EditorStateUnknown EditorState = "unknown"
	// EditorStateFailed indicates that the editor failed to start, generally this is a resource
	// constraint or some other issue in the cluster
	EditorStateFailed = "failed"
	// EditorStateStarting indicates that the pod or Jupyter server is starting
	EditorStateStarting = "starting"
	// EditorStateDoesNotExist indicates that the neither a Deployment or a Service exists for the
	// editor
	EditorStateDoesNotExist = "does not exist"
	// EditorStateRunning indicates a healthy running server
	EditorStateRunning = "running"
	// EditorStateStopping indicates that the Service exists but the Deployment has been marked for
	// deletion
	EditorStateStopping = "stopping"
	// EditorStateStopped indicates that the Service exists but the Deployment has been deleted
	EditorStateStopped = "stopped"
	// EditorStateResuming indicates that the Service exists but the Deployment is starting
	// (currently unused)
	EditorStateResuming = "resuming"
	// EditorStateDestroying indicates that the Service exists but the Deployment has been deleted
	// and that the Volume is also in the process of being deleted (currently unused)
	EditorStateDestroying = "destroying"
)

// EditorStatus represents the current pod state and Jupyter server status of a user's editor session
// This information is pulled from a combination of the Kubernetes API and the JupyterLab API
type EditorStatus struct {
	// Name is the Name value for the BundleMeta with BundleID
	Name string `json:"name"`
	// RealmID is the id of the project the bundle belongs to
	RealmID string `json:"realmID"`
	// BundleID is the id string of the bundle this editor corresponds to
	BundleID string `json:"bundleID"`
	// CreatedAt indicates the time that the editor was started
	CreatedAt *time.Time `json:"createdAt"`
	// LastActivity is a timestamp, it will update with the websocket heartbeat of the editor,
	// so it will be "fresh" if there is an active browser tab out there
	LastActiveAt *time.Time `json:"lastActiveAt"`
	// Status indicates the current state of the editor session
	Status EditorState `json:"status"`
	// StatusMessage is an optional string that can be returned to provide additional details about
	// the current status
	StatusMessage string `json:"statusMessage"`
	// Stages indicates the current flow/process required to complete the current status, this is
	// often empty
	Stages []EditorStage `json:"stages"`
	// EditorURL is the url used to access the jupyterlabs server
	EditorURL string `json:"editorURL"`
	// Labels are the labels from the editor object in k8s, these are for internal use only
	Labels map[string]string `json:"-"`
	// Annotations are the Annotations from the editor object in k8s, these are for internal use only
	Annotations map[string]string `json:"-"`
}

// GetBundleID gets and returns the bundle id for the session from the labels/annotations,
// if the annotation is malformed or missing, the value will be uuid.Nil
func (e *EditorStatus) GetBundleID() string {
	return e.getValue(constants.LabelBundleID)
}

// GetUserID gets and returns the user id for the session from the labels/annotations
func (e *EditorStatus) GetUserID() string {
	return e.getValue(constants.LabelUserID)
}

// GetTenantID gets and returns the tenant id for the session from the labels/annotations
func (e *EditorStatus) GetTenantID() string {
	return e.getValue(constants.LabelTenantID)
}

// GetRealmID extracts the editor realmid from the labels/annotations
func (e *EditorStatus) GetRealmID() string {
	return e.getValue(constants.LabelRealmID)
}

// GetServerAuthToken extracts the editor jupyterlabs token from the labels/annotations
func (e *EditorStatus) GetServerAuthToken() string {
	// should only be found in the annotations
	return e.Annotations[constants.AnnotationTokenKey]
}

// GetBaseURL extracts the editor jupyterlabs baseurl from the labels/annotations
func (e *EditorStatus) GetBaseURL() string {
	return e.Annotations[constants.AnnotationJupyterLabBaseURL]
}

// getValue helps to extract values from first the labels and if not found, then the annotations
func (e *EditorStatus) getValue(key string) string {
	value := e.Labels[key]
	if value == "" {
		value = e.Annotations[key]
	}
	return value
}

// EditorListResponse is the list response of EditorSessions for a user
type EditorListResponse struct {
	Page PageMeta        `json:"page"`
	Data []*EditorStatus `json:"data"`
}

// EditorStage indicates the the steps or stages that the current editor status will transition through.
// For example, if the the Status is "starting" the stages will consist of
// - "pod starting",
// - "pulling from git remote",
// -  "starting server",
// etc.
type EditorStage struct {
	Message string            `json:"message"`
	Status  EditorStageStatus `json:"status"` // this is an enum of 'todo', 'running', 'done', 'failed'
}
