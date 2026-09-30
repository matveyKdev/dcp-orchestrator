package repository

type Node struct {
	Ip      string
	PodCIDR string
}

type Pod struct {
	ServiceName string
	Ip          string
	Port        int
	Status      string
}
