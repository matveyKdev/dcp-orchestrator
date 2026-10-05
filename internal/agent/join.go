package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JoinRequest struct {
	Token string `json:"token"`
}

type JoinResponse struct {
	NodeID           string `json:"node_id"`
	PodCIDR          string `json:"pod_cidr"`
	RegistryUsername string `json:"registry_username"`
	RegistryPassword string `json:"registry_password"`
}

// функция которая отправляет запрос на RESP CP и регает ноду в БД CP и response от CP так же регается на самой ноде
func JoinNode(ctx context.Context, agentInterface AgentInterface, host string, token string) (JoinResponse, error) {
	joinRequest := JoinRequest{Token: token}
	reqBody, err := json.Marshal(joinRequest)
	if err != nil {
		return JoinResponse{}, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", host+"/node/join", bytes.NewReader(reqBody))
	if err != nil {
		return JoinResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return JoinResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return JoinResponse{}, errors.New(string(body))
	}

	var result JoinResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return JoinResponse{}, err
	}
	err = agentInterface.InsertNodeData(ctx, result.NodeID, result.PodCIDR)
	if err != nil {
		panic(err)
	}
	return result, nil
}
