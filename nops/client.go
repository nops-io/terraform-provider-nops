package nops

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// HostURL - Default nOps URL.
const HostURL string = "https://app.nops.io"

// Client - HTTP client to be used by the provider.
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Auth       AuthStruct
}

// AuthStruct - authentication mechanism with an API Key.
type AuthStruct struct {
	ApiKey string `json:"api_key"`
}

// NewClient - instantiates a client for the provider to use.
func NewClient(host, api_key *string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		// Default nOps URL
		HostURL: HostURL,
	}

	if host != nil {
		c.HostURL = *host
	}

	if api_key == nil {
		return &c, nil
	}

	c.Auth = AuthStruct{
		ApiKey: *api_key,
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	token := c.Auth.ApiKey

	req.Header.Set("X-Nops-Api-Key", token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	statusOK := res.StatusCode >= 200 && res.StatusCode < 300
	if !statusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func (c *Client) GetProjects() ([]Project, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/c/admin/projectaws/", c.HostURL), nil)

	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	projects := []Project{}
	err = json.Unmarshal(body, &projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (c *Client) CreateProject(project NewProject) (*Project, error) {
	rb, err := json.Marshal(project)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/c/admin/projectaws/", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	projects := Project{}
	err = json.Unmarshal(body, &projects)
	if err != nil {
		return nil, err
	}

	return &projects, nil
}

func (c *Client) UpdateProject(id int64, project UpdateProject) (*Project, error) {
	rb, err := json.Marshal(project)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/c/admin/projectaws/%d/", c.HostURL, id), strings.NewReader(string(rb)))

	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	projects := Project{}
	err = json.Unmarshal(body, &projects)
	if err != nil {
		return nil, err
	}

	return &projects, nil
}

func (c *Client) DeleteProject(id int64) error {

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/c/admin/projectaws/%d/", c.HostURL, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) NotifyNops(payload Integration) (*IntegrationResponse, error) {
	rb, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/c/aws/integration/", c.HostURL), strings.NewReader(string(rb)))
	req.Header.Set("X-Aws-Account-Number", payload.AccountNumber)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	status := IntegrationResponse{}
	err = json.Unmarshal(body, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func (c *Client) NotifyComputeCopilotOnboarding(payload ComputeCopilotOnboarding) error {
	rb, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/svc/karpenter_manager/agents/terraform/onboarding-confirmation", c.HostURL), strings.NewReader(string(rb)))

	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetComputeCopilotOnboarding(regionName string, accountId string) (*ComputeCopilotOnboarding, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/svc/karpenter_manager/agents/terraform/onboarding-confirmation", c.HostURL), nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("region_name", regionName)
	q.Add("account_id", accountId)
	req.URL.RawQuery = q.Encode()

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	result := ComputeCopilotOnboarding{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) DeleteComputeCopilotOnboarding(regionName string, accountId string) error {

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/svc/karpenter_manager/agents/terraform/onboarding", c.HostURL), nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("region_name", regionName)
	q.Add("account_id", accountId)
	req.URL.RawQuery = q.Encode()

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) NotifyContainerCostBucketSetup(payload ContainerCostBucketSetup) error {
	rb, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/c/admin/container_cost_bucket/setup/", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetContainerCostBucketSetupStatus() (*[]ContainerCostBucket, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/c/admin/container_cost_bucket/", c.HostURL), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	result := []ContainerCostBucket{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetTargetedContainerCostBucketSetupStatus(id int64) (*ContainerCostBucket, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/c/admin/container_cost_bucket/%d", c.HostURL, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	result := ContainerCostBucket{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) DeleteContainerCostBucket(id int64) error {

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/c/admin/container_cost_bucket/%d/", c.HostURL, id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
