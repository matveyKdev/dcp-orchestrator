package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"k8s-pet-project/internal/controlplane/platform/access_token"
	"k8s-pet-project/internal/controlplane/repository"
	"k8s-pet-project/shared/ipam"
	"net/netip"
	"time"
)

type UseCase interface {
	Join(ctx context.Context, token string, host string) (*repository.Node, error)
}

type useCase struct {
	repo repository.ClusterInterface
}

func NewUseCase(repo repository.ClusterInterface) UseCase {
	return &useCase{repo: repo}
}

// реализация логики по коннекту ноды в кластер
func (uc *useCase) Join(ctx context.Context, token string, host string) (*repository.Node, error) {
	clusterInfo, err := uc.repo.GetClusterInfo(ctx)
	if err != nil {
		return nil, err
	}
	isTokenValid, err := access_token.ValidateToken(token, clusterInfo.AccessToken)
	if err != nil {
		return nil, err
	}
	if !isTokenValid {
		return nil, fmt.Errorf("token is invalid")
	}

	nodes, err := uc.repo.GetNodes(ctx)
	if err != nil {
		return nil, err
	}

	clusterCidr, err := netip.ParsePrefix(clusterInfo.PodCidr)
	if err != nil {
		return nil, err
	}
	usedCIDRs := make([]netip.Prefix, 0, len(nodes))
	for _, node := range nodes {
		if node.PodCIDR == "" {
			continue
		}
		prefix, err := netip.ParsePrefix(node.PodCIDR)
		if err != nil {
			return nil, err
		}
		usedCIDRs = append(usedCIDRs, prefix)
	}

	nodeCIDR, err := ipam.AllocateSubnet(clusterCidr, 24, usedCIDRs)
	if err != nil {
		return nil, err
	}

	node := repository.Node{
		Id:            uuid.NewString(),
		ClusterId:     clusterInfo.Id,
		Ip:            host,
		PodCIDR:       nodeCIDR.String(),
		Status:        "joining",
		LastHeartbeat: time.Now(),
	}

	createNode, err := uc.repo.InsertNode(ctx, node)
	if err != nil {
		return nil, err
	}

	return createNode, nil
}
