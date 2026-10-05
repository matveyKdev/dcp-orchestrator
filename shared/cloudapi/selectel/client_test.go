package selectel

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateNode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/compute/v2.1/servers" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("X-Auth-Token"); got != "token" {
			t.Errorf("X-Auth-Token = %q", got)
		}

		var request createNodeEnvelope
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if request.Server.Name != "worker-1" || request.Server.FlavorRef != "flavor-1" {
			t.Errorf("unexpected request: %#v", request)
		}
		if len(request.Server.Networks) != 1 || request.Server.Networks[0].UUID != "network-1" {
			t.Errorf("unexpected networks: %#v", request.Server.Networks)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"server":{"id":"vm-1","adminPass":"generated","links":[]}}`))
	}))
	defer server.Close()

	client := NewSelectelClient("token", server.URL+"/compute/v2.1", server.Client())
	response, err := client.CreateNode(context.Background(), CreateNodeRequest{
		Name:      "worker-1",
		FlavorRef: "flavor-1",
		ImageRef:  "image-1",
		Networks:  []CreateNetwork{{UUID: "network-1"}},
		KeyName:   "default-key",
	})
	if err != nil {
		t.Fatalf("CreateNode() error = %v", err)
	}
	if response.Server.ID != "vm-1" || response.Server.AdminPass != "generated" {
		t.Fatalf("unexpected response: %#v", response)
	}
}

func TestGetNodes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/compute/v2.1/servers/detail" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if got := r.URL.Query().Get("limit"); got != "50" {
			t.Errorf("limit = %q", got)
		}
		if got := r.URL.Query().Get("marker"); got != "vm-previous" {
			t.Errorf("marker = %q", got)
		}
		if got := r.Header.Get("X-Auth-Token"); got != "token" {
			t.Errorf("X-Auth-Token = %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"servers":[{"id":"vm-1","name":"worker-1","status":"ACTIVE","addresses":{"private":[{"addr":"10.0.0.2","version":4}]}}]}`))
	}))
	defer server.Close()

	client := NewSelectelClient("token", server.URL+"/compute/v2.1", server.Client())
	response, err := client.GetNodes(context.Background(), 50, "vm-previous")
	if err != nil {
		t.Fatalf("GetNodes() error = %v", err)
	}
	if len(response.Servers) != 1 || response.Servers[0].ID != "vm-1" {
		t.Fatalf("unexpected servers: %#v", response.Servers)
	}
	if got := response.Servers[0].Addresses["private"][0].Address; got != "10.0.0.2" {
		t.Fatalf("private address = %q", got)
	}
}
