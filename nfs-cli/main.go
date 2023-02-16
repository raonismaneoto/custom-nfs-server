package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strings"

	"github.com/raonismaneoto/custom-nfs-server/helpers"
	"github.com/raonismaneoto/custom-nfs-server/nfs-cli/models"
	client "github.com/raonismaneoto/custom-nfs-server/nfs-client"
)

const MaxBytesPerRequest int32 = 10000

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		panic("no command has been provided")
	}

	client, username, hostname := getCmdsCommonData()

	switch command := args[0]; command {
	case "mount":
		if len(args) != 3 {
			panic("Usage: mount <origin> <absolute_destination>")
		}
		log.Println("exec mount command")

		path := args[1]
		response, err := client.Mount(username+"@"+hostname, path)
		if err != nil {
			panic("error while mounting request: " + err.Error())
		}
		absDestPath := args[2]

		var filesMd []models.Metadata
		err = json.Unmarshal(response, &filesMd)
		if err != nil {
			panic("error unmarshalling mount response. Error: " + err.Error())
		}

		for _, fmd := range filesMd {
			func(currFmd models.Metadata) {
				fm, err := os.OpenFile(absDestPath+"/"+currFmd.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
				if err != nil {
					panic("unable to open/create " + absDestPath + ". Error: " + err.Error())
				}
				defer fm.Close()
				parsedFmd, err := json.Marshal(currFmd)
				if err != nil {
					panic(err.Error())
				}
				if _, err := fm.Write(parsedFmd); err != nil {
					panic("unable to write to: " + absDestPath + ". Error: " + err.Error())
				}
			}(fmd)
		}
	case "save":
		log.Println("exec save command")
		if len(args) != 3 {
			panic("Usage: save <absolute_origin> <destination>")
		}
		content := make(chan []byte)
		proceed := make(chan string)
		origin := args[1]
		destination := args[2]
		go client.Save(username+"@"+hostname, destination, content, proceed)
		f, err := os.Open(origin)
		if err != nil {
			panic("unable to open file")
		}
		defer f.Close()
		stat, err := f.Stat()
		if err != nil {
			panic("unable to get file stat")
		}
		chuncks := int32(math.Ceil(float64(stat.Size()) / float64(MaxBytesPerRequest)))
		for i := int32(0); i < chuncks; i++ {
			<-proceed
			data, err := helpers.ReadFileChunk(origin, i*MaxBytesPerRequest, MaxBytesPerRequest)
			if err != nil {
				panic("error while reading file chunk: " + err.Error())
			}
			content <- data
		}
		<-proceed
		close(content)
		<-proceed
	case "read":
		log.Println("exec read command")
	case "chpem":
		log.Println("exec chpem command")
	default:
		fmt.Printf("Usage: \nexecutable <command> <args> [options]")
	}
}

func getCmdsCommonData() (*client.Client, string, string) {
	client := client.NewClient(os.Getenv("server_addr"))
	cmd := exec.Command("whoami")
	out, err := cmd.Output()
	if err != nil {
		panic("error while getting username: " + err.Error())
	}
	username := string(out)
	username = strings.Replace(username, "\n", "", -1)
	cmd = exec.Command("hostname")
	out, err = cmd.Output()
	if err != nil {
		panic("error while getting hostname: " + err.Error())
	}
	hostname := string(out)
	hostname = strings.Replace(hostname, "\n", "", -1)
	return client, username, hostname
}
