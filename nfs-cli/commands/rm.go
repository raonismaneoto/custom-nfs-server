package commands

import (
	"os"
	"strings"

	"github.com/raonismaneoto/custom-nfs-server/nfs-cli/helpers"
	"github.com/raonismaneoto/custom-nfs-server/nfs-cli/models"
)

func ExecRm(path string, cconfig models.CommandConfiguration) {
	rmPath := path

	if strings.Contains(path, "meta") {
		err := os.RemoveAll(rmPath)
		if err != nil {
			panic(err.Error())
		}
		rmPath, err = helpers.GetRemotePath(path)
	}

	err := cconfig.Client.Rm(cconfig.Username+"@"+cconfig.Hostname, rmPath)
	if err != nil {
		panic(err.Error())
	}
}
