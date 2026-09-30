package delivery

import (
	"context"
	"k8s-pet-project/internal/controlplane/application"
)

type Handler struct {
	useCase *application.UseCase
}

func New(useCase application.UseCase) *Handler {
	return &Handler{useCase: &useCase}
}

func (h *Handler) Join(ctx context.Context) (string, error) {
	return "", nil
}

func (h *Handler) HeartBeat(ctx context.Context) (string, error) {
	return "", nil
}
