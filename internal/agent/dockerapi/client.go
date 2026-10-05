package dockerapi

import (
	"github.com/moby/moby/client"
)

type Client struct {
	client *client.Client
}

func New() (*Client, error) {
	cli, err := client.New(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: cli,
	}, nil
}
