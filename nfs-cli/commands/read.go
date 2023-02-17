package commands

import (
	"log"
	"os"

	"github.com/raonismaneoto/custom-nfs-server/nfs-cli/models"
)

func ExecRead(localPath, remotePath string, cconfig models.CommandConfiguration) {
	_, err := os.Stat(localPath)
	if err == nil {
		err := os.Remove(localPath)
		if err != nil {
			panic(err.Error())
		}
	}
	f, err := os.OpenFile(localPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("unable to open/create file: " + err.Error())
	}

	content := make(chan []byte)
	proceed := make(chan string)

	go cconfig.Client.Read(cconfig.Username+"@"+cconfig.Hostname, remotePath, content, proceed)

	for {
		proceed <- "proceed"

		currContent, ok := <-content
		if !ok {
			close(proceed)
			f.Close()
			return
		}

		if _, err := f.Write(currContent); err != nil {
			log.Println("unable to write to file")
		}
	}
}
