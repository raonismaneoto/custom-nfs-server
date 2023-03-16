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
		remotePath, err := helpers.GetRemotePath(path)
		if err != nil {
			panic(err.Error())
		}
		err = os.RemoveAll(rmPath)
		if err != nil {
			panic(err.Error())
		}
		os.RemoveAll(rmPath[:len(rmPath)-4])
		rmPath = remotePath
	}

	err := cconfig.Client.Rm(cconfig.Username+"@"+cconfig.Hostname, rmPath)
	if err != nil {
		panic(err.Error())
	}
}
