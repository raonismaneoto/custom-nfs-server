package commands

import (
	"log"
	"math"
	"os"

	"github.com/raonismaneoto/custom-nfs-server/helpers"
	"github.com/raonismaneoto/custom-nfs-server/nfs-cli/models"
)

func ExecSave(origin, destination string, cconfig models.CommandConfiguration) {
	f, err := os.Open(origin)
	if err != nil {
		panic("unable to open file")
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		panic("unable to get file stat")
	}
	if stat.IsDir() {
		err := cconfig.Client.Save(cconfig.Username+"@"+cconfig.Hostname, destination, []byte{})
		if err != nil {
			log.Println(err.Error())
			panic(err.Error())
		}
		return
	}
	content := make(chan []byte)
	proceed := make(chan string)
	go cconfig.Client.SaveAsync(cconfig.Username+"@"+cconfig.Hostname, destination, content, proceed)
	chuncks := int32(math.Ceil(float64(stat.Size()) / float64(cconfig.MaxBytesPerRequest)))
	for i := int32(0); i < chuncks; i++ {
		<-proceed
		data, err := helpers.ReadFileChunk(origin, i*cconfig.MaxBytesPerRequest, cconfig.MaxBytesPerRequest)
		if err != nil {
			panic("error while reading file chunk: " + err.Error())
		}
		content <- data
	}
	<-proceed
	close(content)
	<-proceed
}
