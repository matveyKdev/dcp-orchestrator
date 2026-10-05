package yandexcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateNode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/compute/v1/instances" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer token" {
			t.Errorf("Authorization = %q", got)
		}

		var request CreateNodeRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if request.FolderID != "folder-1" || request.Name != "worker-1" {
			t.Errorf("unexpected request: %#v", request)
		}
		if request.BootDiskSpec.DiskSpec == nil || request.BootDiskSpec.DiskSpec.ImageID != "image-1" {
			t.Errorf("unexpected boot disk: %#v", request.BootDiskSpec)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"operation-1","done":false,"metadata":{"instanceId":"vm-1"}}`))
	}))
	defer server.Close()

	client := NewYandexCloudClient("token", "folder-1", server.Client())
	client.baseURL = server.URL
	response, err := client.CreateNode(context.Background(), CreateNodeRequest{
		Name:       "worker-1",
		ZoneID:     "ru-central1-a",
		PlatformID: "standard-v3",
		ResourcesSpec: ResourcesSpec{
			Memory: "4294967296",
			Cores:  "2",
		},
		BootDiskSpec: AttachedDiskSpec{
			AutoDelete: true,
			DiskSpec: &DiskSpec{
				Size:    "21474836480",
				ImageID: "image-1",
			},
		},
		NetworkInterfaceSpecs: []NetworkInterfaceSpec{{SubnetID: "subnet-1"}},
	})
	if err != nil {
		t.Fatalf("CreateNode() error = %v", err)
	}
	if response.ID != "operation-1" || response.Done {
		t.Fatalf("unexpected response: %#v", response)
	}
}

func TestGetNodes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/compute/v1/instances" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if got := r.URL.Query().Get("folderId"); got != "folder-1" {
			t.Errorf("folderId = %q", got)
		}
		if got := r.URL.Query().Get("pageSize"); got != "25" {
			t.Errorf("pageSize = %q", got)
		}
		if got := r.URL.Query().Get("pageToken"); got != "next token" {
			t.Errorf("pageToken = %q", got)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer token" {
			t.Errorf("Authorization = %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"instances":[{"id":"vm-1","name":"worker-1","status":"RUNNING"}],"nextPageToken":"page-2"}`))
	}))
	defer server.Close()

	client := NewYandexCloudClient("token", "folder-1", server.Client())
	client.baseURL = server.URL

	response, err := client.GetNodes(context.Background(), 25, "next token")
	if err != nil {
		t.Fatalf("GetNodes() error = %v", err)
	}
	if len(response.Instances) != 1 || response.Instances[0].ID != "vm-1" {
		t.Fatalf("unexpected instances: %#v", response.Instances)
	}
	if response.NextPageToken != "page-2" {
		t.Fatalf("NextPageToken = %q", response.NextPageToken)
	}
}

func TestGetNodesReturnsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"message":"denied"}`, http.StatusForbidden)
	}))
	defer server.Close()

	client := NewYandexCloudClient("token", "folder-1", server.Client())
	client.baseURL = server.URL

	if _, err := client.GetNodes(context.Background(), 0, ""); err == nil {
		t.Fatal("GetNodes() error = nil")
	}
}
