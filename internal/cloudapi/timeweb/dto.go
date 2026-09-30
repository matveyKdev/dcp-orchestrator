package timeweb

type MetaTotal struct {
	Total int `json:"total"`
}

type OsInfo struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type SoftwareInfo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type DiskInfo struct {
	Id         int    `json:"id"`
	Size       int    `json:"size"`
	Used       int    `json:"used"`
	Type       string `json:"type"`
	IsMounted  bool   `json:"is_mounted"`
	IsSystem   bool   `json:"is_system"`
	SystemName string `json:"system_name"`
	Status     string `json:"status"`
}

type ImageInfo struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	IsCustom bool   `json:"is_custom"`
}

type Ip struct {
	Type   string `json:"type"`
	Ip     string `json:"ip"`
	Ptr    string `json:"ptr"`
	IsMain bool   `json:"is_main"`
}

type Network struct {
	Id             string `json:"id"`
	Type           string `json:"type"`
	NatMode        string `json:"nat_mode"`
	Bandwidth      int    `json:"bandwidth"`
	Ips            []Ip   `json:"ips"`
	IsDdosGuard    bool   `json:"is_ddos_guard"`
	IsImageMounted bool   `json:"is_image_mounted"`
	BockedPorts    []int  `json:"bocked_ports"`
}

type Server struct {
	Id               int          `json:"id"`
	Name             string       `json:"name"`
	Comment          string       `json:"comment"`
	CreatedAt        string       `json:"created_at"`
	Os               OsInfo       `json:"os"`
	Software         SoftwareInfo `json:"software"`
	PresetId         int          `json:"preset_id"`
	Location         string       `json:"location"`
	ConfiguratorId   int          `json:"configurator_id"`
	BootMode         string       `json:"boot_mode"`
	Status           string       `json:"status"`
	StartAt          string       `json:"start_at"`
	IsDdosGuard      bool         `json:"is_ddos_guard"`
	IsMasterSSH      bool         `json:"is_master_ssh"`
	IsDedicatedCPU   bool         `json:"is_dedicated_cpu"`
	GPU              int          `json:"gpu"`
	CPU              int          `json:"cpu"`
	CPUFrequency     string       `json:"cpu_frequency"`
	RAM              int          `json:"ram"`
	Disks            []DiskInfo   `json:"disks"`
	AvatarId         string       `json:"avatar_id"`
	AvatarLink       string       `json:"avatar_link"`
	VNCPass          string       `json:"vnc_pass"`
	RootPass         string       `json:"root_pass"`
	Image            ImageInfo    `json:"image"`
	Networks         []Network    `json:"networks"`
	CloudInit        string       `json:"cloud_init"`
	IsQemuAgent      bool         `json:"is_qemu_agent"`
	AvailabilityZone string       `json:"availability_zone"`
}
type GetNodesResponse struct {
	Meta       MetaTotal `json:"meta"`
	Servers    []Server  `json:"servers"`
	ResponseId string    `json:"response_id"`
}

type ServerConfiguration struct {
	ConfigurationId int `json:"configuration_id"`
	Disk            int `json:"disk"`
	CPU             int `json:"cpu"`
	RAM             int `json:"ram"`
	GPU             int `json:"gpu"`
}
type NetworkConfiguration struct {
	Id              string   `json:"id"`
	FloatingIp      string   `json:"floating_ip"`
	LocalIp         string   `json:"local_ip"`
	Ip              string   `json:"ip"`
	NetworkDriveIds []string `json:"network_drive_ids"`
}
type CreateNodeRequest struct {
	Configuration    ServerConfiguration  `json:"configuration"`
	IsDdosGuard      bool                 `json:"is_ddos_guard"`
	OSId             int                  `json:"os_id"`
	ImageId          string               `json:"image_id"`
	SoftwareId       string               `json:"software_id"`
	PresetId         int                  `json:"preset_id"`
	Bandwidth        int                  `json:"bandwidth"`
	Name             string               `json:"name"`
	AvatarId         string               `json:"avatar_id"`
	Comment          string               `json:"comment"`
	SSHKeysIds       string               `json:"ssh_keys_ids"`
	IsLocalNetwork   bool                 `json:"is_local_network"`
	Network          NetworkConfiguration `json:"network"`
	CloudInit        string               `json:"cloud_init"`
	AvailabilityZone string               `json:"availability_zone"`
	ProjectId        string               `json:"project_id"`
	Hostname         string               `json:"hostname"`
}

type CreateNodeResponse struct {
	Server     Server `json:"server"`
	ResponseId string `json:"response_id"`
}
