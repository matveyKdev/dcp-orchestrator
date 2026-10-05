package cloudru

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Flavor struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	CPU              int    `json:"cpu"`
	RAM              int    `json:"ram"`
	GPU              int    `json:"gpu"`
	Oversubscription string `json:"oversubscription"`
	FreeTier         bool   `json:"free_tier"`
}

type FloatingIP struct {
	ID        string `json:"id"`
	IPAddress string `json:"ip_address"`
	State     string `json:"state"`
	Name      string `json:"name"`
}

type NetworkInterface struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	IPAddress  string      `json:"ip_address"`
	FloatingIP *FloatingIP `json:"floating_ip,omitempty"`
	Primary    bool        `json:"primary"`
	State      string      `json:"state"`
	Type       string      `json:"type"`
}

type AvailabilityZone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type VirtualMachine struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	ProjectID        string             `json:"project_id"`
	Description      string             `json:"description"`
	Tags             []Tag              `json:"tags"`
	Flavor           Flavor             `json:"flavor"`
	State            string             `json:"state"`
	Locked           bool               `json:"locked"`
	CreatedTime      string             `json:"created_time"`
	ModifiedTime     string             `json:"modified_time"`
	Interfaces       []NetworkInterface `json:"interfaces"`
	AvailabilityZone AvailabilityZone   `json:"availability_zone"`
	IsSerialReady    bool               `json:"is_serial_ready"`
}

type GetNodesResponse struct {
	Items  []VirtualMachine `json:"items"`
	Offset int              `json:"offset"`
	Limit  int              `json:"limit"`
	Total  int              `json:"total"`
}

type CreateInterface struct {
	Type                     string   `json:"type"`
	SubnetID                 string   `json:"subnet_id,omitempty"`
	SubnetName               string   `json:"subnet_name,omitempty"`
	InterfaceSecurityEnabled *bool    `json:"interface_security_enabled,omitempty"`
	SecurityGroups           []string `json:"security_groups,omitempty"`
	SecurityGroupNames       []string `json:"security_group_names,omitempty"`
	NewExternalIP            bool     `json:"new_external_ip,omitempty"`
	AttachExternalIPID       string   `json:"attach_external_ip_id,omitempty"`
	AttachExternalIPName     string   `json:"attach_external_ip_name,omitempty"`
	IPAddress                string   `json:"ip_address,omitempty"`
}

// CreateDisk describes either an existing disk (DiskID/DiskName) or a new one
// (Name, Size and DiskTypeID/DiskTypeName).
type CreateDisk struct {
	DiskID       string `json:"disk_id,omitempty"`
	DiskName     string `json:"disk_name,omitempty"`
	DiskTypeID   string `json:"disk_type_id,omitempty"`
	DiskTypeName string `json:"disk_type_name,omitempty"`
	Name         string `json:"name,omitempty"`
	Size         int    `json:"size,omitempty"`
}

type CreateNodeRequest struct {
	ProjectID            string            `json:"project_id"`
	AvailabilityZoneID   string            `json:"availability_zone_id,omitempty"`
	AvailabilityZoneName string            `json:"availability_zone_name,omitempty"`
	Name                 string            `json:"name"`
	Description          string            `json:"description,omitempty"`
	FlavorID             string            `json:"flavor_id,omitempty"`
	FlavorName           string            `json:"flavor_name,omitempty"`
	ImageID              string            `json:"image_id,omitempty"`
	ImageName            string            `json:"image_name,omitempty"`
	PlacementGroupID     string            `json:"placement_group_id,omitempty"`
	PlacementGroupName   string            `json:"placement_group_name,omitempty"`
	Interfaces           []CreateInterface `json:"interfaces,omitempty"`
	Disks                []CreateDisk      `json:"disks"`
	ImageMetadata        map[string]string `json:"image_metadata,omitempty"`
	TagIDs               []string          `json:"tag_ids,omitempty"`
	TagNames             []string          `json:"tag_names,omitempty"`
	CloudInit            string            `json:"cloud_init,omitempty"`
}

type CreateNodeResponse []VirtualMachine
