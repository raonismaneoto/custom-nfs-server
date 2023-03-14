package commands

import (
	"log"
	"os"

	"github.com/raonismaneoto/custom-nfs-server/nfs-cli/helpers"
	"github.com/raonismaneoto/custom-nfs-server/nfs-cli/models"
)

func ExecRead(metaPath string, cconfig models.CommandConfiguration) {
	filePath := metaPath[:len(metaPath)-4]
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("unable to open/create file: " + err.Error())
	}
	defer f.Close()

	content := make(chan []byte)
	errors := make(chan error)
	remotePath, err := helpers.GetRemotePath(metaPath)
	if err != nil {
		panic(err.Error())
	}

	go cconfig.Client.Read(cconfig.Username+"@"+cconfig.Hostname, remotePath, content, errors)

	for {
		select {
		case currContent, ok := <-content:
			if !ok {
				f.Close()
				return
			}
			if _, err := f.Write(currContent); err != nil {
				log.Println("unable to write to file")
			}
		case err, ok := <-errors:
			if !ok {
				f.Close()
				return
			}
			panic(err.Error())
		}
	}
}
