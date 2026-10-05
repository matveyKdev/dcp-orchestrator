package timeweb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const defaultBaseURL = "https://api.timeweb.cloud/api/v1"

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewTimewebClient(token string, httpClient *http.Client) *Client {
	return &Client{
		baseURL:    defaultBaseURL,
		token:      token,
		httpClient: httpClient,
	}
}

func (c *Client) GetNodes(ctx context.Context, limit, offset int32) (GetNodesResponse, error) {
	url := fmt.Sprintf(
		"%s/servers?limit=%d&offset=%d",
		c.baseURL,
		limit,
		offset,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return GetNodesResponse{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return GetNodesResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return GetNodesResponse{}, fmt.Errorf("server responded with code %d: %s", resp.StatusCode, body)
	}
	var result GetNodesResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return GetNodesResponse{}, err
	}
	return result, nil
}

func (c *Client) CreateNode(ctx context.Context, node CreateNodeRequest) (CreateNodeResponse, error) {
	var data = strings.NewReader(`{
	  "is_ddos_guard": false,
	  "os_id": 79,
	  "comment": "comment",
	  "name": "name",
	  "preset_id": 2451,
	  "is_local_network": false
	}`)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.timeweb.cloud/api/v1/servers", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return CreateNodeResponse{}, fmt.Errorf("server responded with code %d: %s", resp.StatusCode, body)
	}
	var result CreateNodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return CreateNodeResponse{}, err
	}
	return result, nil
}
