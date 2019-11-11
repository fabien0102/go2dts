package client

import (
	"time"

	"github.com/contiamo/labs/pkg/sql"
	uuid "github.com/satori/go.uuid"
)

// TriggerType is a string enum that encodes the source/type of the TriggeredBy struct.TriggerType
// Current this can be `user`, `api` or `schedule`.  Future values may include `webhook` and `eventstream`.
type TriggerType string

func (t TriggerType) String() string {
	return string(t)
}

// ExecutionState describes the overall state/status of an Execution: running, success, failed
type ExecutionState string

func (s ExecutionState) String() string {
	return string(s)
}

const (
	// TriggeredByUser indidcates that a user triggered the job via the API
	TriggeredByUser TriggerType = "user"
	// TriggeredByAPI indidcates that a request token manually triggered the job via the API
	TriggeredByAPI TriggerType = "apikey"
	// TriggeredByCRON indicates that the job execution was created by a CRON schedule
	TriggeredByCRON TriggerType = "schedule"
	// TriggeredByHook is a place holder for a potential webhook flow that is distinguishable from
	// the REST API calls
	TriggeredByHook TriggerType = "webhook"
	// TriggeredByUnknown is a fallback default value, it should rarely or never actually be used
	TriggeredByUnknown TriggerType = "unknown"

	// ExecutionStateUnknown indicates that the execution is in a misconfigured state that we can
	// not determine, this value exists as a fallback.
	ExecutionStateUnknown ExecutionState = "unknown"
	// ExecutionStateFailed indicates that the execution is stopped and returned a non-success
	// status code
	ExecutionStateFailed ExecutionState = "failed"
	// ExecutionStateSuccess indicates the the execution is stopped and returned a success status
	// code
	ExecutionStateSuccess ExecutionState = "success"
	// ExecutionStateRunning indicates that the execution is still running
	ExecutionStateRunning ExecutionState = "running"
)

// JobResponse represents the definition of a bundle job with
// the most recent status and output URL.
// This is distinct from the execution status or logs of
// a job which come from the k8s-job-controller service
type JobResponse struct {
	// ID is provided for consistency and client ease of use, it will always
	// match Name
	ID           string            `json:"id"`
	LastRunAt    *time.Time        `json:"lastRunAt"`
	BundleID     uuid.UUID         `json:"bundleID"`
	URL          string            `json:"url"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	NotebookPath string            `json:"notebookPath"`
	Schedule     string            `json:"schedule"`
	Environment  sql.JSONStringMap `json:"environment"`
	Secrets      []string          `json:"secrets"`
	Success      bool              `json:"success"`
	Internal     bool              `json:"internal"`
}

// JobListResponse contains the paging and the data array for the job List endpoint
type JobListResponse struct {
	Data []*JobResponse `json:"data"`
}

// JobRunRequest is the expected POST body for a manual job invocation, the caller
// is allowed to specify one-time override of the Environment variables on the job
type JobRunRequest struct {
	Environment sql.JSONStringMap `json:"environment"`
}

// JobRunResponse contains the execution id of a specific job execution, this comes
// from the manual Run endpoint and can be used to query for the execution status
// and logs
type JobRunResponse struct {
	ExecutionID string `json:"executionID"`
}

// JobSpecification contains the arguments used and returned by a job execution. These
// values will be populated from k8s-job-controller
type JobSpecification struct {
	Schedule    string            `json:"schedule"`
	Image       string            `json:"image"`
	Command     []string          `json:"command"`
	Environment sql.JSONStringMap `json:"environment"`
	Secrets     []string          `json:"secrets"`
}

// TriggeredBy is used to record who/what is responsible for a specific Job execution
// The type
type TriggeredBy struct {
	Type TriggerType `json:"type"`
	ID   string      `json:"id"`
	Name string      `json:"name"`
}

// ExecutionResponse contains the execution details of a job, as returned and parsed
// from the job controllers, this includes the status as well as the command details
// As will the JobResponse object, the JobID will always match JobName
type ExecutionResponse struct {
	ID            string           `json:"id"`
	JobID         string           `json:"jobID"`
	BundleID      string           `json:"bundleID"`
	JobName       string           `json:"jobName"`
	StartedAt     time.Time        `json:"startedAt"`
	CompletedAt   *time.Time       `json:"completedAt"`
	Success       bool             `json:"success"`
	Status        ExecutionState   `json:"status"`
	Specification JobSpecification `json:"specification"`
	TriggeredBy   TriggeredBy      `db:"-" json:"triggeredBy"`
}

// GetStatus returns the status value as determined from the other fields of the response:
// CompletedAt and Success
func (e ExecutionResponse) GetStatus() (status ExecutionState) {
	status = ExecutionStateRunning

	if e.StartedAt.IsZero() {
		status = ExecutionStateUnknown
	}

	if e.CompletedAt == nil {
		return status
	}

	if e.Success {
		status = ExecutionStateSuccess
	} else {
		status = ExecutionStateFailed
	}

	return status
}

// ExecutionStatus is the subset of fields related to an ExecutionResponse status.  This is used
// in the LogsResponse to let the user know the current status of the logs, this allows displaying
// summary information next to the raw output.
type ExecutionStatus struct {
	StartedAt   time.Time  `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt"`
	Success     bool       `json:"success"`
}

// ExecutionListResponse contains the paging and the data array for the job execution List endpoint
type ExecutionListResponse struct {
	Data []*ExecutionResponse `json:"data"`
}

// LogMessage represents a line of output from a job, it will include the command output and timestamp
type LogMessage struct {
	Timestamp time.Time `json:"timestamp"`
	Msg       string    `json:"msg"`
}

// LogsResponse contains the paging and the data array for the job execution logs
type LogsResponse struct {
	Status ExecutionStatus `json:"status"`
	Data   []LogMessage    `json:"data"`
}
