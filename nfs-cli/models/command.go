package models

import (
	client "github.com/raonismaneoto/custom-nfs-server/nfs-client"
)

type CommandConfiguration struct {
	Username           string
	Hostname           string
	MaxBytesPerRequest int32
	Client             client.Client
}
