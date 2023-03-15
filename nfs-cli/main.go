package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/raonismaneoto/custom-nfs-server/nfs-cli/commands"
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
	cconfig := models.CommandConfiguration{
		Username:           username,
		Hostname:           hostname,
		Client:             *client,
		MaxBytesPerRequest: MaxBytesPerRequest,
	}

	switch command := args[0]; command {
	case "mount":
		if len(args) != 3 {
			panic("Usage: mount <origin> <absolute_destination>")
		}
		log.Println("exec mount command")

		path := args[1]
		absDestPath := args[2]
		commands.ExecMount(path, absDestPath, cconfig)
		log.Println("mount executed successfully")
	case "save":
		log.Println("exec save command")
		if len(args) != 3 {
			panic("Usage: save <absolute_origin> <destination>")
		}

		origin := args[1]
		destination := args[2]
		commands.ExecSave(origin, destination, cconfig)
		log.Println("save executed successfully")
	case "read":
		log.Println("exec read command")
		if len(args) != 2 {
			panic("Usage: read <meta_path>")
		}

		metaPath := args[1]
		commands.ExecRead(metaPath, cconfig)
		log.Println("read executed successfully")
	case "chpem":
		log.Println("exec chpem command")
		if len(args) != 4 {
			panic("Usage: chpem <file_path> <operation> <user>\n Operation: add/rm")
		}

		path := args[1]
		operation := args[2]
		user := args[3]
		commands.ExecChpem(path, operation, user, cconfig)

		log.Println("chpem executed successfully")
	case "rm":
		log.Println("exec rm command")
		if len(args) != 2 {
			panic("Usage: rm <path> \npath can be either a local meta file absolute path or a remote file path")
		}
		path := args[1]
		commands.ExecRm(path, cconfig)

		log.Println("rm executed successfully")
	case "update":
		log.Println("exec update command")
		if len(args) != 2 {
			panic("Usage: update <local_absolute_mounted_path>\npath cannot point to a directory")
		}

		path := args[1]
		commands.ExecUpdate(path, cconfig)

		log.Println("update executed successfully")
	default:
		fmt.Printf("Usage: nfs <command> <args> [options]")
	}
}

func getCmdsCommonData() (*client.Client, string, string) {
	client := client.NewClient(os.Getenv("custom_nfs_server_addr"))
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
