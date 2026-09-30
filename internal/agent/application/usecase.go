package application

import (
	"context"
	"net/http"
)

type UseCase interface {
	Join(ctx context.Context, host string, token string) (string, error)
}

type useCase struct {
	http *http.Client
}

func NewUseCase(http *http.Client) UseCase {
	return &useCase{http}
}

func (uc *useCase) Join(ctx context.Context, host string, token string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://"+host, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := uc.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return resp.Status, nil
}
