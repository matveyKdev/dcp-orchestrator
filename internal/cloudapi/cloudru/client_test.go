package cloudru

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateNode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1.1/vms" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer token" {
			t.Errorf("Authorization = %q", got)
		}

		var request []CreateNodeRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if len(request) != 1 || request[0].ProjectID != "project-1" || request[0].Name != "worker-1" {
			t.Fatalf("unexpected request: %#v", request)
		}
		if len(request[0].Disks) != 1 || request[0].Disks[0].Size != 20 {
			t.Errorf("unexpected disks: %#v", request[0].Disks)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`[{"id":"vm-1","name":"worker-1","project_id":"project-1","state":"creating"}]`))
	}))
	defer server.Close()

	client := NewCloudRUClient("token", "project-1", server.Client())
	client.baseURL = server.URL
	response, err := client.CreateNode(context.Background(), CreateNodeRequest{
		Name:               "worker-1",
		AvailabilityZoneID: "zone-1",
		FlavorID:           "flavor-1",
		ImageID:            "image-1",
		Disks:              []CreateDisk{{Name: "worker-1-boot", Size: 20, DiskTypeID: "disk-type-1"}},
		Interfaces:         []CreateInterface{{Type: "regular", SubnetID: "subnet-1"}},
	})
	if err != nil {
		t.Fatalf("CreateNode() error = %v", err)
	}
	if len(response) != 1 || response[0].ID != "vm-1" {
		t.Fatalf("unexpected response: %#v", response)
	}
}

func TestGetNodes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/vms" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if got := r.URL.Query().Get("project_id"); got != "project-1" {
			t.Errorf("project_id = %q", got)
		}
		if got := r.URL.Query().Get("limit"); got != "20" {
			t.Errorf("limit = %q", got)
		}
		if got := r.URL.Query().Get("offset"); got != "40" {
			t.Errorf("offset = %q", got)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer token" {
			t.Errorf("Authorization = %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"items":[{"id":"vm-1","name":"worker-1","project_id":"project-1","state":"running","interfaces":[{"id":"nic-1","ip_address":"10.0.0.2","primary":true}]}],"offset":40,"limit":20,"total":61}`))
	}))
	defer server.Close()

	client := NewCloudRUClient("token", "project-1", server.Client())
	client.baseURL = server.URL

	response, err := client.GetNodes(context.Background(), 20, 40)
	if err != nil {
		t.Fatalf("GetNodes() error = %v", err)
	}
	if len(response.Items) != 1 || response.Items[0].ID != "vm-1" {
		t.Fatalf("unexpected VMs: %#v", response.Items)
	}
	if response.Total != 61 {
		t.Fatalf("Total = %d", response.Total)
	}
}
