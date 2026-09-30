package yandexcloud

import "encoding/json"

type Resources struct {
	Memory       string `json:"memory"`
	Cores        string `json:"cores"`
	CoreFraction string `json:"coreFraction"`
	GPUs         string `json:"gpus"`
}

type OneToOneNAT struct {
	Address   string `json:"address"`
	IPVersion string `json:"ipVersion"`
}

type PrimaryAddress struct {
	Address     string       `json:"address"`
	OneToOneNAT *OneToOneNAT `json:"oneToOneNat,omitempty"`
}

type NetworkInterface struct {
	Index            string          `json:"index"`
	MACAddress       string          `json:"macAddress"`
	SubnetID         string          `json:"subnetId"`
	PrimaryV4Address *PrimaryAddress `json:"primaryV4Address,omitempty"`
	PrimaryV6Address *PrimaryAddress `json:"primaryV6Address,omitempty"`
	SecurityGroupIDs []string        `json:"securityGroupIds"`
}

type AttachedDisk struct {
	Mode       string `json:"mode"`
	DeviceName string `json:"deviceName"`
	AutoDelete bool   `json:"autoDelete"`
	DiskID     string `json:"diskId"`
}

type Instance struct {
	ID                string             `json:"id"`
	FolderID          string             `json:"folderId"`
	CreatedAt         string             `json:"createdAt"`
	Name              string             `json:"name"`
	Description       string             `json:"description"`
	Labels            map[string]string  `json:"labels"`
	ZoneID            string             `json:"zoneId"`
	PlatformID        string             `json:"platformId"`
	Resources         Resources          `json:"resources"`
	Status            string             `json:"status"`
	BootDisk          AttachedDisk       `json:"bootDisk"`
	SecondaryDisks    []AttachedDisk     `json:"secondaryDisks"`
	NetworkInterfaces []NetworkInterface `json:"networkInterfaces"`
	FQDN              string             `json:"fqdn"`
	ServiceAccountID  string             `json:"serviceAccountId"`
}

type GetNodesResponse struct {
	Instances     []Instance `json:"instances"`
	NextPageToken string     `json:"nextPageToken"`
}

type ResourcesSpec struct {
	Memory       string `json:"memory"`
	Cores        string `json:"cores"`
	CoreFraction string `json:"coreFraction,omitempty"`
	GPUs         string `json:"gpus,omitempty"`
}

type DiskSpec struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	TypeID      string `json:"typeId,omitempty"`
	Size        string `json:"size"`
	BlockSize   string `json:"blockSize,omitempty"`
	ImageID     string `json:"imageId,omitempty"`
	SnapshotID  string `json:"snapshotId,omitempty"`
	KMSKeyID    string `json:"kmsKeyId,omitempty"`
}

type AttachedDiskSpec struct {
	Mode       string    `json:"mode,omitempty"`
	DeviceName string    `json:"deviceName,omitempty"`
	AutoDelete bool      `json:"autoDelete"`
	DiskSpec   *DiskSpec `json:"diskSpec,omitempty"`
	DiskID     string    `json:"diskId,omitempty"`
}

type OneToOneNATSpec struct {
	IPVersion string `json:"ipVersion"`
	Address   string `json:"address,omitempty"`
}

type PrimaryAddressSpec struct {
	Address         string           `json:"address,omitempty"`
	OneToOneNATSpec *OneToOneNATSpec `json:"oneToOneNatSpec,omitempty"`
}

type NetworkInterfaceSpec struct {
	Index                string              `json:"index,omitempty"`
	SubnetID             string              `json:"subnetId"`
	PrimaryV4AddressSpec *PrimaryAddressSpec `json:"primaryV4AddressSpec,omitempty"`
	SecurityGroupIDs     []string            `json:"securityGroupIds,omitempty"`
}

type SchedulingPolicy struct {
	Preemptible bool `json:"preemptible"`
}

type CreateNodeRequest struct {
	FolderID              string                 `json:"folderId"`
	Name                  string                 `json:"name"`
	Description           string                 `json:"description,omitempty"`
	Labels                map[string]string      `json:"labels,omitempty"`
	ZoneID                string                 `json:"zoneId"`
	PlatformID            string                 `json:"platformId"`
	ResourcesSpec         ResourcesSpec          `json:"resourcesSpec"`
	Metadata              map[string]string      `json:"metadata,omitempty"`
	BootDiskSpec          AttachedDiskSpec       `json:"bootDiskSpec"`
	SecondaryDiskSpecs    []AttachedDiskSpec     `json:"secondaryDiskSpecs,omitempty"`
	NetworkInterfaceSpecs []NetworkInterfaceSpec `json:"networkInterfaceSpecs"`
	Hostname              string                 `json:"hostname,omitempty"`
	SchedulingPolicy      *SchedulingPolicy      `json:"schedulingPolicy,omitempty"`
	ServiceAccountID      string                 `json:"serviceAccountId,omitempty"`
}

type OperationError struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Details []json.RawMessage `json:"details"`
}

type CreateNodeResponse struct {
	ID          string          `json:"id"`
	Description string          `json:"description"`
	CreatedAt   string          `json:"createdAt"`
	CreatedBy   string          `json:"createdBy"`
	ModifiedAt  string          `json:"modifiedAt"`
	Done        bool            `json:"done"`
	Metadata    json.RawMessage `json:"metadata"`
	Error       *OperationError `json:"error,omitempty"`
	Response    json.RawMessage `json:"response,omitempty"`
}
