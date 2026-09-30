package selectel

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

type Client struct {
	computeEndpoint string
	token           string
	httpClient      *http.Client
}

// NewSelectelClient creates a client for the Compute endpoint returned by the
// Selectel/OpenStack service catalog, for example
// https://ru-9.cloud.api.selcloud.ru/compute/v2.1.
func NewSelectelClient(token, computeEndpoint string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		computeEndpoint: strings.TrimRight(computeEndpoint, "/"),
		token:           token,
		httpClient:      httpClient,
	}
}

// GetNodes returns one page of detailed OpenStack servers. Pass the ID of the
// last server from the previous page as marker to retrieve the next page.
func (c *Client) GetNodes(ctx context.Context, limit int, marker string) (GetNodesResponse, error) {
	endpoint, err := url.Parse(c.computeEndpoint + "/servers/detail")
	if err != nil {
		return GetNodesResponse{}, fmt.Errorf("build selectel servers URL: %w", err)
	}

	query := endpoint.Query()
	if limit > 0 {
		query.Set("limit", strconv.Itoa(limit))
	}
	if marker != "" {
		query.Set("marker", marker)
	}
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return GetNodesResponse{}, fmt.Errorf("create selectel servers request: %w", err)
	}
	req.Header.Set("X-Auth-Token", c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return GetNodesResponse{}, fmt.Errorf("list selectel servers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
		return GetNodesResponse{}, fmt.Errorf("selectel API responded with code %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var result GetNodesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return GetNodesResponse{}, fmt.Errorf("decode selectel servers response: %w", err)
	}
	return result, nil
}

func (c *Client) CreateNode(ctx context.Context, node CreateNodeRequest) (CreateNodeResponse, error) {
	body, err := json.Marshal(createNodeEnvelope{Server: node})
	if err != nil {
		return CreateNodeResponse{}, fmt.Errorf("encode selectel server request: %w", err)
	}

	endpoint := c.computeEndpoint + "/servers"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return CreateNodeResponse{}, fmt.Errorf("create selectel server request: %w", err)
	}
	req.Header.Set("X-Auth-Token", c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return CreateNodeResponse{}, fmt.Errorf("create selectel server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		responseBody, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
		return CreateNodeResponse{}, fmt.Errorf("selectel API responded with code %d: %s", resp.StatusCode, strings.TrimSpace(string(responseBody)))
	}

	var result CreateNodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return CreateNodeResponse{}, fmt.Errorf("decode selectel create server response: %w", err)
	}
	return result, nil
}
