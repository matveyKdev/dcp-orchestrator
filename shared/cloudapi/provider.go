package cloudapi

import "context"

type Provider interface {
	CreateNode(ctx context.Context, spec NodeSpec) (*Node, error)
	DeleteNode(ctx context.Context, id string) error
	GetNode(ctx context.Context, id string) (*Node, error)
}

type NodeSpec struct {
	Name string

	CPU       int
	MemoryMiB int
	DiskGiB   int

	UserData string
}

type Node struct {
	ID        string
	Name      string
	PrivateIP string
	PublicIP  string
	Status    string
}
