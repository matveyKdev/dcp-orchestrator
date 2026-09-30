package repository

import "time"

type Node struct {
	Name          string
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
