package yandexcloud

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

const defaultBaseURL = "https://compute.api.cloud.yandex.net"

type Client struct {
	baseURL    string
	token      string
	folderID   string
	httpClient *http.Client
}

func NewYandexCloudClient(token, folderID string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		baseURL:    defaultBaseURL,
		token:      token,
		folderID:   folderID,
		httpClient: httpClient,
	}
}

// GetNodes returns one page of Compute Cloud instances. Pass NextPageToken from
// the previous response as pageToken to retrieve the next page.
func (c *Client) GetNodes(ctx context.Context, pageSize int64, pageToken string) (GetNodesResponse, error) {
	endpoint, err := url.Parse(strings.TrimRight(c.baseURL, "/") + "/compute/v1/instances")
	if err != nil {
		return GetNodesResponse{}, fmt.Errorf("build yandex cloud instances URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("folderId", c.folderID)
	if pageSize > 0 {
		query.Set("pageSize", strconv.FormatInt(pageSize, 10))
	}
	if pageToken != "" {
		query.Set("pageToken", pageToken)
	}
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return GetNodesResponse{}, fmt.Errorf("create yandex cloud instances request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return GetNodesResponse{}, fmt.Errorf("list yandex cloud instances: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
		return GetNodesResponse{}, fmt.Errorf("yandex cloud API responded with code %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var result GetNodesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return GetNodesResponse{}, fmt.Errorf("decode yandex cloud instances response: %w", err)
	}
	return result, nil
}

// CreateNode starts asynchronous instance creation and returns the operation
// that can be used to track its progress.
func (c *Client) CreateNode(ctx context.Context, node CreateNodeRequest) (CreateNodeResponse, error) {
	node.FolderID = c.folderID
	body, err := json.Marshal(node)
	if err != nil {
		return CreateNodeResponse{}, fmt.Errorf("encode yandex cloud instance request: %w", err)
	}

	endpoint := strings.TrimRight(c.baseURL, "/") + "/compute/v1/instances"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return CreateNodeResponse{}, fmt.Errorf("create yandex cloud instance request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return CreateNodeResponse{}, fmt.Errorf("create yandex cloud instance: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
		return CreateNodeResponse{}, fmt.Errorf("yandex cloud API responded with code %d: %s", resp.StatusCode, strings.TrimSpace(string(responseBody)))
	}

	var result CreateNodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return CreateNodeResponse{}, fmt.Errorf("decode yandex cloud create instance response: %w", err)
	}
	return result, nil
}
