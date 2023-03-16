package commands

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/raonismaneoto/custom-nfs-server/nfs-cli/models"
)

func ExecMount(path, absDestPath string, cconfig models.CommandConfiguration) {
	response, err := cconfig.Client.Mount(cconfig.Username+"@"+cconfig.Hostname, path)
	if err != nil {
		panic("error while mounting request: " + err.Error())
	}
	var filesMd []models.Metadata
	err = json.Unmarshal(response, &filesMd)
	if err != nil {
		panic("error unmarshalling mount response. Error: " + err.Error())
	}

	for _, fmd := range filesMd {
		func(currFmd models.Metadata) {
			splitPath := strings.Split(currFmd.Path, "/")
			if currFmd.Dir {
				os.MkdirAll(absDestPath+"/"+strings.Join(splitPath[:len(splitPath)-1], "/"), 0777)
			}
			fm, err := os.OpenFile(absDestPath+"/"+currFmd.Path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
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
}
