package delivery

type JoinRequest struct {
	Token string `json:"token"`
}

type JoinResponse struct {
	NodeID           string `json:"node_id"`
	PodCIDR          string `json:"pod_cidr"`
	RegistryUsername string `json:"registry_username"`
	RegistryPassword string `json:"registry_password"`
}
