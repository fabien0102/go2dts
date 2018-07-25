package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/contiamo/labs/pkg/sql"
	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// FunctionTemplate read from template.yml within root of a language template folder
type FunctionTemplate struct {
	Language string `yaml:"language"`
	FProcess string `yaml:"fprocess"`
}

// FunctionDefinition contains the image and command information needed
// execute the function via Labs
type FunctionDefinition struct {
	Image       string            `json:"image"`
	Command     string            `json:"command"`
	Environment sql.JSONStringMap `json:"environment,omitempty" yaml:"environment,omitempty"`
	Secrets     []string          `json:"secrets,omitempty" yaml:"secrets,omitempty"`
	Build       FunctionBuild     `json:"build,omitempty" yaml:"build,omitempty"`
}

// EnvArray returns the environment variables as an array of strings
func (f *FunctionDefinition) EnvArray() []string {
	var envArray []string
	for key, value := range f.Environment {
		envArray = append(envArray, fmt.Sprintf("%s=%s", strings.ToUpper(key), value))
	}
	return envArray
}

// JobDefinition contains the path and configuration information to execute a notebook on
// a schedule (or manually)
type JobDefinition struct {
	NotebookPath string            `json:"notebookPath" yaml:"notebook_path"`
	Environment  sql.JSONStringMap `json:"environment"`
	Secrets      []string          `json:"secrets"`
	Schedule     string            `json:"schedule"`
	Description  string            `json:"description"`
}

// FunctionBuild defines the build arguments for the function.  If this is
// empty, a template will be build from the specified Function.Image that
// uses `watchdog`.  Otherwise, specify a customer dockerfile here to fully
// control the function image.  The `Dockerfile` must be a relative path
// from the project root.
type FunctionBuild struct {
	Dockerfile string   `json:"dockerfile,omitempty"`
	Args       []string `json:"args,omitempty" yaml:"args,omitempty"`
	Labels     []string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Skip       bool     `json:"skip,omitempty" yaml:"skip,omitempty"`
}

// EditConfig describes the docker container used to launch Jupyterlab for editing
// a bundle.
type EditConfig struct {
	Image       string            `json:"image"`
	Environment sql.JSONStringMap `json:"environment,omitempty" yaml:"environment,omitempty"`
	Secrets     []string          `json:"secret,omitempty" yaml:"secrets,omitempty"`
}

// EnvArray returns a string array representing the environment variables.  This
// is the required format for the docker sdk.
func (e *EditConfig) EnvArray() []string {
	var envArray []string
	for key, value := range e.Environment {
		envArray = append(envArray, fmt.Sprintf("%s=%s", strings.ToUpper(key), value))
	}
	return envArray
}

// Bundle contains the Labs configurations
type Bundle struct {
	Version   string                        `json:"apiVersion" yaml:"api_version"`
	Name      string                        `json:"name"`
	Edit      EditConfig                    `json:"edit"`
	Jobs      map[string]JobDefinition       `json:"jobs,omitempty" yaml:"jobs,omitempty"`
	Functions map[string]FunctionDefinition `json:"functions,omitempty" yaml:"functions,omitempty"`
}

// Value implements the Value interfance and provides the the database value in
// a type that the driver can handle, in paritcular as a string.
func (b Bundle) Value() (driver.Value, error) {
	j, err := json.Marshal(b)
	return j, err
}

// Scan implements the Scanner interface that will scan the Postgres JSON payload
// into the Bundle *b
func (b *Bundle) Scan(src interface{}) error {
	var (
		source []byte
	)

	switch v := src.(type) {
	case []byte:
		source, _ = src.([]byte)
	case string:
		// t, ok := src.(string)
		source = []byte(src.(string))
	default:
		return fmt.Errorf("type assertion failed, received %s, %s", v, src)
	}

	err := json.Unmarshal(source, b)
	if err != nil {
		return err
	}

	return nil
}

// BundleInfo contains the bundle as well as path information about that bundle
type BundleInfo struct {
	Name   string
	Path   string
	Bundle Bundle
}

// LoadBundle with parse and load the bundle config at the specified path into
// the BundleInfo
func (info *BundleInfo) LoadBundle(path string, verbose bool) error {

	abs, err := filepath.Abs(path)
	if err != nil {
		return errors.Wrapf(err, "Could not construct absolute path to: %s", path)
	}

	filename := filepath.Base(abs)
	path = filepath.Dir(abs)
	for path != "/" {
		if verbose {
			log.Printf("Looking for config in %s\n", path)
		}

		searchAbs := fmt.Sprintf("%s/%s", path, filename)
		yamlBytes, err := ioutil.ReadFile(searchAbs)
		if err != nil {
			// go up a level and try again
			path = filepath.Dir(path)
			continue
		}

		err = yaml.Unmarshal(yamlBytes, &info.Bundle)
		if err != nil {
			return errors.Wrapf(err, "can not load the yaml in %s", searchAbs)
		}
		info.Name = filename
		info.Path = path

		return ValdiateBundle(yamlBytes)
	}

	return errors.New("Config file not found")
}
