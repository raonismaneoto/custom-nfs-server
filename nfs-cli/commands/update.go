package commands

import (
	"os/exec"
	"strings"

	"github.com/raonismaneoto/custom-nfs-server/nfs-cli/helpers"
	"github.com/raonismaneoto/custom-nfs-server/nfs-cli/models"
)

func ExecUpdate(path string, cconfig models.CommandConfiguration) {
	remotePath, err := helpers.GetRemotePath(path + "meta")
	if err != nil {
		panic(err.Error())
	}
	cmd := exec.Command("cp", path, path+".tmp")
	_, err = cmd.Output()
	if err != nil {
		panic(err.Error())
	}
	ExecRm(path, cconfig)
	ExecSave(path+".tmp", remotePath, cconfig)
	splitPath := strings.Split(path, "/")
	pathDir := strings.Join(splitPath[:len(splitPath)-1], "/")
	ExecMount(remotePath, pathDir, cconfig)
	cmd = exec.Command("cp", path+".tmp", path)
	_, err = cmd.Output()
	if err != nil {
		panic(err.Error())
	}
}
