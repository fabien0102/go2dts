package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/contiamo/idp/http/client"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/util/homedir"
)

const (
	labsAPIRoot = "/api"
)

// HTTPClient interface to represents the minimal required http client interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewIDPClient creates a new authenticated HttpClient
func NewIDPClient(idpAddr, tenant, user, password string) (idpClient *client.Client, err error) {
	idpClient, err = client.New(idpAddr)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create idp client")
	}
	idpClient.SetInsecureSkipVerify()
	idpClient.SetSSHKey(homedir.HomeDir() + "/.ssh/id_labs")
	idpClient.SetTenant(tenant)
	idpClient.SetEmail(user)
	idpClient.SetPassword(password)

	if err = idpClient.Login(); err != nil {
		return nil, errors.Wrap(err, "idp client failed to login")
	}
	// call verify to ensure that we have a request token loaded
	err = idpClient.Verify()
	if err != nil {
		return nil, errors.Wrap(err, "unable to verify authentication")
	}
	err = idpClient.SetupSSHAuthentication()
	if err != nil {
		return nil, errors.Wrap(err, "can not configure SSH authentication")
	}

	return idpClient, err
}

type bundleClientImpl struct {
	idpClient HTTPClient
	labsAddr  string
	TenantID  string
	RealmID   string
}

// New returns a Labs client instance, that can be used to interact with the Labs Bundle APIs.
func New(authedClient HTTPClient, labsAddr, tenantID, realmID string) BundleClient {
	return &bundleClientImpl{
		idpClient: authedClient,
		labsAddr:  labsAddr,
		TenantID:  tenantID,
		RealmID:   realmID,
	}
}

func (cli *bundleClientImpl) Do(r *http.Request) (*http.Response, error) {
	logrus.Debugf("%s %s", r.Method, r.URL.String())
	return cli.idpClient.Do(r)
}

func (cli *bundleClientImpl) GetTenantID() string {
	return cli.TenantID
}

func (cli *bundleClientImpl) SetRealm(id string) {
	cli.RealmID = id
}

func (cli *bundleClientImpl) Register(ctx context.Context, req *RegisterBundleRequest) (*BundleResponse, error) {
	req.TenantID = cli.TenantID
	reqBytes, err := json.Marshal(&req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to serialize the bundle request")
	}

	requestURI := cli.getURL()
	request, err := http.NewRequest(http.MethodPost, requestURI, bytes.NewReader(reqBytes))
	if err != nil {
		return nil, errors.Wrap(err, "could not create http request")
	}
	request = request.WithContext(ctx)
	request.Header.Set("Content-Type", "application/json")
	res, err := cli.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "bundle post failed")
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	serverBundle := &BundleResponse{}

	// s, _ := httputil.DumpResponse(res, true)
	// fmt.Printf("DEBUG: labs api client response\n\n:%s", s)

	switch res.StatusCode {
	case http.StatusOK, http.StatusCreated:
		err := json.NewDecoder(res.Body).Decode(serverBundle)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read response body")
		}
	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, errors.New("unauthorized request")
	default:
		bytesOut, err := ioutil.ReadAll(res.Body)
		if err == nil && len(bytesOut) > 0 {
			return nil, errors.Errorf("unexpected status code %v, message: %v", res.StatusCode, string(bytesOut))
		}
		return nil, errors.Errorf("unexpected status code %v", res.StatusCode)
	}

	return serverBundle, nil
}

func (cli *bundleClientImpl) Unregister(ctx context.Context, id uuid.UUID) error {
	requestURI := cli.getURL(id.String())
	request, err := http.NewRequest(http.MethodDelete, requestURI, nil)
	if err != nil {
		return errors.Wrap(err, "could not create http request")
	}
	request = request.WithContext(ctx)
	resp, err := cli.Do(request)
	if err != nil {
		return errors.Wrap(err, "could not send http request")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	switch resp.StatusCode {
	case http.StatusOK, http.StatusNoContent:
		return nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return errors.New("unauthorized request")
	default:
		bytesOut, err := ioutil.ReadAll(resp.Body)
		if err == nil && len(bytesOut) > 0 {
			return errors.Errorf("unexpected status code %v, message: %v", resp.StatusCode, string(bytesOut))
		}
		return errors.Errorf("unexpected status code %v", resp.StatusCode)
	}
}

func (cli *bundleClientImpl) List(ctx context.Context) (*BundleListResponse, error) {
	requestURI := cli.getURL()
	request, err := http.NewRequest(http.MethodGet, requestURI, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create http request")
	}
	request = request.WithContext(ctx)
	resp, err := cli.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "could not send http request")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	var bundles *BundleListResponse
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		err := json.NewDecoder(resp.Body).Decode(&bundles)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read response body")
		}
	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, errors.New("unauthorized request")
	default:
		bytesOut, err := ioutil.ReadAll(resp.Body)
		if err == nil && len(bytesOut) > 0 {
			return nil, errors.Errorf("unexpected status code %v, message: %v", resp.StatusCode, string(bytesOut))
		}
		return nil, errors.Errorf("unexpected status code %v", resp.StatusCode)
	}
	return bundles, nil
}

func (cli *bundleClientImpl) Sync(ctx context.Context, id uuid.UUID) (*BundleSyncResponse, error) {
	requestURI := cli.getURL(id.String(), "sync")
	request, err := http.NewRequest(http.MethodPost, requestURI, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create http request")
	}
	request = request.WithContext(ctx)
	resp, err := cli.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "could not send http request")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	response := &BundleSyncResponse{}
	switch resp.StatusCode {
	case http.StatusOK:
		err := json.NewDecoder(resp.Body).Decode(response)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read response body")
		}
	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, errors.New("unauthorized request")
	default:
		bytesOut, err := ioutil.ReadAll(resp.Body)
		if err == nil && len(bytesOut) > 0 {
			return nil, errors.Errorf("unexpected status code %v, message: %v", resp.StatusCode, string(bytesOut))
		}
		return nil, errors.Errorf("unexpected status code %v", resp.StatusCode)
	}

	return response, nil
}

func (cli *bundleClientImpl) Deploy(ctx context.Context, id uuid.UUID) (*DeployResponse, error) {
	requestURI := cli.getURL(id.String(), "deploy")
	request, err := http.NewRequest(http.MethodPost, requestURI, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create http request")
	}
	request = request.WithContext(ctx)
	resp, err := cli.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "could not send http request")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	deployResponse := &DeployResponse{}
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		err := json.NewDecoder(resp.Body).Decode(deployResponse)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read response body")
		}
	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, errors.New("unauthorized request")
	default:
		bytesOut, err := ioutil.ReadAll(resp.Body)
		if err == nil && len(bytesOut) > 0 {
			return nil, errors.Errorf("unexpected status code %v, message: %v", resp.StatusCode, string(bytesOut))
		}
		return nil, errors.Errorf("unexpected status code %v", resp.StatusCode)
	}

	return deployResponse, nil
}

func (cli *bundleClientImpl) Undeploy(ctx context.Context, id uuid.UUID) error {
	requestURI := cli.getURL(id.String(), "undeploy")
	request, err := http.NewRequest(http.MethodPost, requestURI, nil)
	if err != nil {
		return errors.Wrap(err, "could not create http request")
	}
	request = request.WithContext(ctx)
	resp, err := cli.Do(request)
	if err != nil {
		return errors.Wrap(err, "could not send http request")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	switch resp.StatusCode {
	case http.StatusOK, http.StatusNoContent:
		return nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return errors.New("unauthorized request")
	default:
		bytesOut, err := ioutil.ReadAll(resp.Body)
		if err == nil && len(bytesOut) > 0 {
			return errors.Errorf("unexpected status code %v, message: %v", resp.StatusCode, string(bytesOut))
		}
		return errors.Errorf("unexpected status code %v", resp.StatusCode)
	}
}

func (cli *bundleClientImpl) StartEditSession(ctx context.Context, id uuid.UUID) (url string, err error) {
	requestURI := cli.getURL(id.String(), "edit")
	request, err := http.NewRequest(http.MethodGet, requestURI, nil)
	if err != nil {
		return "", errors.Wrap(err, "could not create http request")
	}
	request = request.WithContext(ctx)
	resp, err := cli.Do(request)
	if err != nil {
		return "", errors.Wrap(err, "could not send http request")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	urlStruct := &struct {
		URL string `json:"url"`
	}{}
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated:
		if err = json.NewDecoder(resp.Body).Decode(urlStruct); err != nil {
			return "", err
		}
	case http.StatusUnauthorized, http.StatusForbidden:
		return "", errors.New("unauthorized request")
	default:
		bytesOut, err := ioutil.ReadAll(resp.Body)
		if err == nil && len(bytesOut) > 0 {
			return "", errors.Errorf("unexpected status code %v, message: %v", resp.StatusCode, string(bytesOut))
		}
		return "", errors.Errorf("unexpected status code %v", resp.StatusCode)
	}
	return urlStruct.URL, nil
}

func (cli *bundleClientImpl) StopEditSession(ctx context.Context, id uuid.UUID, deleteVolumes bool) error {
	requestURI := cli.getURL(id.String(), "edit")
	request, err := http.NewRequest(http.MethodDelete, requestURI, nil)
	if err != nil {
		return errors.Wrap(err, "could not create http request")
	}
	request = request.WithContext(ctx)
	resp, err := cli.Do(request)
	if err != nil {
		return errors.Wrap(err, "could not send http request")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		return nil
	case http.StatusUnauthorized, http.StatusForbidden:
		return errors.New("unauthorized request")
	case http.StatusBadRequest, http.StatusNotFound:
		// stop is already done
		return nil
	default:
		bytesOut, err := ioutil.ReadAll(resp.Body)
		if err == nil && len(bytesOut) > 0 {
			return errors.Errorf("unexpected status code %v, message: %v", resp.StatusCode, string(bytesOut))
		}
		return errors.Errorf("unexpected status code %v", resp.StatusCode)
	}
}

func (cli *bundleClientImpl) Status(ctx context.Context, id uuid.UUID) (*BundleResponse, error) {
	requestURI := cli.getURL(id.String())
	request, err := http.NewRequest(http.MethodGet, requestURI, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create http request")
	}
	request = request.WithContext(ctx)
	res, err := cli.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "bundle post failed")
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	serverBundle := &BundleResponse{}
	switch res.StatusCode {
	case http.StatusOK:
		err := json.NewDecoder(res.Body).Decode(serverBundle)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read response body")
		}
	case http.StatusUnauthorized, http.StatusForbidden:
		return nil, errors.New("unauthorized request")
	default:
		bytesOut, err := ioutil.ReadAll(res.Body)
		if err == nil && len(bytesOut) > 0 {
			return nil, errors.Errorf("unexpected status code %v, message: %v", res.StatusCode, string(bytesOut))
		}
		return nil, errors.Errorf("unexpected status code %v", res.StatusCode)
	}
	return serverBundle, nil
}

func (cli *bundleClientImpl) getURL(optionals ...string) string {
	args := []string{
		cli.labsAddr,
		cli.TenantID,
		"api",
		"v1",
		"realms",
		cli.RealmID,
		"bundles",
	}
	args = append(args, optionals...)
	res := strings.Join(args, "/")

	return res
}
