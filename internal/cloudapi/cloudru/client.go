package cloudru

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const defaultBaseURL = "https://compute.api.cloud.ru"

type Client struct {
	baseURL    string
	token      string
	projectID  string
	httpClient *http.Client
}

func NewCloudRUClient(token, projectID string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		baseURL:    defaultBaseURL,
		token:      token,
		projectID:  projectID,
		httpClient: httpClient,
	}
}

func (c *Client) GetNodes(ctx context.Context, limit, offset int32) (GetNodesResponse, error) {
	endpoint, err := url.Parse(strings.TrimRight(c.baseURL, "/") + "/api/v1/vms")
	if err != nil {
		return GetNodesResponse{}, fmt.Errorf("build cloud.ru VMs URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("project_id", c.projectID)
	if limit > 0 {
		query.Set("limit", strconv.FormatInt(int64(limit), 10))
	}
	if offset > 0 {
		query.Set("offset", strconv.FormatInt(int64(offset), 10))
	}
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return GetNodesResponse{}, fmt.Errorf("create cloud.ru VMs request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return GetNodesResponse{}, fmt.Errorf("list cloud.ru VMs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
		return GetNodesResponse{}, fmt.Errorf("cloud.ru API responded with code %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var result GetNodesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return GetNodesResponse{}, fmt.Errorf("decode cloud.ru VMs response: %w", err)
	}
	return result, nil
}

// CreateNode creates one VM through the current Cloud.ru v1.1 API. Although
// this method creates one VM, the provider API accepts and returns JSON arrays.
func (c *Client) CreateNode(ctx context.Context, node CreateNodeRequest) (CreateNodeResponse, error) {
	node.ProjectID = c.projectID
	body, err := json.Marshal([]CreateNodeRequest{node})
	if err != nil {
		return nil, fmt.Errorf("encode cloud.ru VM request: %w", err)
	}

	endpoint := strings.TrimRight(c.baseURL, "/") + "/api/v1.1/vms"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create cloud.ru VM request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create cloud.ru VM: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		responseBody, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
		return nil, fmt.Errorf("cloud.ru API responded with code %d: %s", resp.StatusCode, strings.TrimSpace(string(responseBody)))
	}

	var result CreateNodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode cloud.ru create VM response: %w", err)
	}
	return result, nil
}
