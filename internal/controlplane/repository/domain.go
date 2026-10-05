package repository

import "time"

type Cluster struct {
	Id          string
	PodCidr     string
	AccessToken string
}

type Node struct {
	Id            string
	ClusterId     string
	Ip            string
	PodCIDR       string
	Status        string
	LastHeartbeat time.Time
}

type Pod struct {
	ServiceName string
	Ip          string
	Port        int
	Status      string
}
