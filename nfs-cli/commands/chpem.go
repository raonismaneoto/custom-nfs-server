package commands

import (
	"github.com/raonismaneoto/custom-nfs-server/nfs-cli/models"
)

func ExecChpem(path, operation, user string, cconfig models.CommandConfiguration) {
	err := cconfig.Client.Chpem(cconfig.Username+"@"+cconfig.Hostname, user, path, operation)
	if err != nil {
		panic(err.Error())
	}
}
