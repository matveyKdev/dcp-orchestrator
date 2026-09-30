package registryapi

import (
	"context"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

type Client struct {
	auth authn.Authenticator
}

func NewClient(username, password string) *Client {
	return &Client{
		auth: &authn.Basic{
			Username: username,
			Password: password,
		},
	}
}

func (c *Client) GetImage(
	ctx context.Context,
	image string,
) error {
	ref, err := name.ParseReference(image)
	if err != nil {
		return err
	}

	img, err := remote.Image(
		ref,
		remote.WithContext(ctx),
		remote.WithAuth(c.auth),
	)
	if err != nil {
		return err
	}

	digest, err := img.Digest()
	if err != nil {
		return err
	}

	_ = digest

	return nil
}
