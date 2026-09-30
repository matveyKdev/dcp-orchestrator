package selectel

import "encoding/json"

type Link struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
}

type Address struct {
	Address    string `json:"addr"`
	Version    int    `json:"version"`
	Type       string `json:"OS-EXT-IPS:type"`
	MACAddress string `json:"OS-EXT-IPS-MAC:mac_addr"`
}

type SecurityGroup struct {
	Name string `json:"name"`
}

type Server struct {
	ID               string               `json:"id"`
	Name             string               `json:"name"`
	Status           string               `json:"status"`
	TenantID         string               `json:"tenant_id"`
	UserID           string               `json:"user_id"`
	Created          string               `json:"created"`
	Updated          string               `json:"updated"`
	AccessIPv4       string               `json:"accessIPv4"`
	AccessIPv6       string               `json:"accessIPv6"`
	AvailabilityZone string               `json:"OS-EXT-AZ:availability_zone"`
	PowerState       int                  `json:"OS-EXT-STS:power_state"`
	VMState          string               `json:"OS-EXT-STS:vm_state"`
	TaskState        *string              `json:"OS-EXT-STS:task_state"`
	Addresses        map[string][]Address `json:"addresses"`
	Flavor           json.RawMessage      `json:"flavor"`
	Image            json.RawMessage      `json:"image"`
	SecurityGroups   []SecurityGroup      `json:"security_groups"`
	Tags             []string             `json:"tags"`
	Links            []Link               `json:"links"`
}

type GetNodesResponse struct {
	Servers     []Server `json:"servers"`
	ServerLinks []Link   `json:"servers_links"`
}

type CreateNetwork struct {
	UUID    string `json:"uuid,omitempty"`
	Port    string `json:"port,omitempty"`
	FixedIP string `json:"fixed_ip,omitempty"`
	Tag     string `json:"tag,omitempty"`
}

type CreateSecurityGroup struct {
	Name string `json:"name"`
}

type BlockDeviceMapping struct {
	UUID                string `json:"uuid"`
	SourceType          string `json:"source_type"`
	DestinationType     string `json:"destination_type"`
	BootIndex           string `json:"boot_index,omitempty"`
	VolumeSize          int    `json:"volume_size,omitempty"`
	DeleteOnTermination bool   `json:"delete_on_termination,omitempty"`
}

type CreateNodeRequest struct {
	Name                 string                `json:"name"`
	FlavorRef            string                `json:"flavorRef"`
	ImageRef             string                `json:"imageRef,omitempty"`
	Networks             []CreateNetwork       `json:"networks"`
	SecurityGroups       []CreateSecurityGroup `json:"security_groups,omitempty"`
	KeyName              string                `json:"key_name,omitempty"`
	UserData             string                `json:"user_data,omitempty"`
	Metadata             map[string]string     `json:"metadata,omitempty"`
	ConfigDrive          bool                  `json:"config_drive,omitempty"`
	AvailabilityZone     string                `json:"OS-EXT-AZ:availability_zone,omitempty"`
	BlockDeviceMappingV2 []BlockDeviceMapping  `json:"block_device_mapping_v2,omitempty"`
}

type createNodeEnvelope struct {
	Server CreateNodeRequest `json:"server"`
}

type CreatedServer struct {
	ID        string `json:"id"`
	AdminPass string `json:"adminPass"`
	Links     []Link `json:"links"`
}

type CreateNodeResponse struct {
	Server CreatedServer `json:"server"`
}
