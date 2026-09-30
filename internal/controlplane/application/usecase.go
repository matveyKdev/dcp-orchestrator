package application

import (
	"context"
	"k8s-pet-project/internal/controlplane/repository"
)

type UseCase interface {
	Join(ctx context.Context, host string, token string) (string, error)
}

type useCase struct {
	repo *repository.DatabaseInterface
}

func NewUseCase(repo *repository.DatabaseInterface) UseCase {
	return &useCase{repo: repo}
}

// реализация логики по коннекту ноды в кластер
func (uc *useCase) Join(ctx context.Context, host string, token string) (string, error) {

	return "", nil
}
