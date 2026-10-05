package delivery

import (
	"context"
	"encoding/json"
	"k8s-pet-project/internal/controlplane/application"
	"net"
	"net/http"
)

type Handler struct {
	useCase application.UseCase
}

func New(useCase application.UseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) Join(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req JoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "invalid remote address", http.StatusBadRequest)
		return
	}
	result, err := h.useCase.Join(r.Context(), req.Token, host)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	response := JoinResponse{
		NodeID:  result.Id,
		PodCIDR: result.PodCIDR,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HeartBeat(ctx context.Context) (string, error) {
	return "", nil
}
