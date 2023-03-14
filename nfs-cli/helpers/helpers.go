package helpers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/raonismaneoto/custom-nfs-server/nfs-cli/models"
)

func GetRemotePath(metaPath string) (string, error) {
	f, err := os.Open(metaPath)
	if err != nil {
		log.Println("unable to open " + err.Error())
		return "", err
	}
	defer f.Close()

	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	var md *models.Metadata
	json.Unmarshal(byteValue, &md)

	return md.Path[:len(md.Path)-4], nil
}
