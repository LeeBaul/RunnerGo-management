package rao

import "time"

type GetMachineListParam struct {
	ServerType int32  `json:"server_type"`
	Name       string `json:"name"`
}

type GetMachineListResponse struct {
	Region     string    `json:"region"`
	IP         string    `json:"ip"`
	Port       int32     `json:"port"`
	Weight     int       `json:"weight"`
	Name       string    `json:"name"`
	ServerType int32     `json:"server_type"`
	Status     int32     `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
