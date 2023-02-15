package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/raonismaneoto/custom-nfs-server/nfs-server/models"
	"golang.org/x/exp/slices"
)

const MetaFileSuffix string = "meta"

type Server struct {
	root string
}

func New(root string) *Server {
	if _, err := os.Stat(root); err != nil {
		os.Mkdir(root, 0777)
	}
	return &Server{
		root: root,
	}
}

func (s *Server) Save(id, path string, content []byte) error {
	log.Println("Save call received.")
	log.Println("saving")
	f, err := os.OpenFile(s.root+path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("unable to open/create %v", path)
		return err
	}

	defer f.Close()

	fm, err := os.OpenFile(s.root+path+MetaFileSuffix, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("unable to open/create %v", path+MetaFileSuffix)
		return err
	}

	defer fm.Close()

	metadata := models.Metadata{OwnerID: id, Size: float64(len(content)),
		Dir: len(content) == 0, AllowList: []string{}}

	byteValue, err := ioutil.ReadAll(fm)
	if err != nil {
		log.Println(err.Error())
		// return err
	}

	if err == nil && len(byteValue) > 0 {
		json.Unmarshal(byteValue, &metadata)
		metadata.Size = metadata.Size + float64(len(content))
	}

	if id == metadata.OwnerID && !slices.Contains(metadata.AllowList, id) {
		metadata.AllowList = append(metadata.AllowList, id)
	}

	mMd, err := json.Marshal(metadata)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	if _, err := fm.Write(mMd); err != nil {
		log.Println("unable to write to %v", path+MetaFileSuffix)
		return err
	}

	// TODO: check if it appends the content
	if _, err := f.Write(content); err != nil {
		log.Println("unable to write to %v", path)
	}

	return err
}

func (s *Server) Read(id, path string, offset, limit int32) ([]byte, error) {
	//check if id is allowed to access its content
	if _, err := s.readMetaData(id, path); err != nil {
		return nil, err
	}

	f, err := os.Open(s.root + path)
	if err != nil {
		log.Println("unable to open file %v", path)
		log.Println(err)
		return nil, err
	}

	stat, err := f.Stat()

	if err != nil {
		log.Println("unable to get file stat")
		return nil, err
	}

	if int32(stat.Size()) <= offset {
		log.Println("the file size is smaller than or equal to the offset")
		return nil, errors.New("the file size is smaller than or equal to the offset")
	}

	if int32(stat.Size()) < limit {
		limit = int32(stat.Size())
	}

	content := make([]byte, limit)

	if _, err := f.ReadAt(content, int64(offset)); err != nil {
		log.Println("unable to read file chunck: ", err)
		return nil, err
	}

	return content, nil
}

func (s *Server) GetMetaData(id, path string) ([]models.Metadata, error) {
	log.Println("going go call read metadata")
	md, err := s.readMetaData(id, path)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	log.Println("checking dir")
	if !md.Dir {
		log.Println("not a dir")
		return []models.Metadata{*md}, nil
	}

	log.Println("it is a dir, going to execute ls")
	cmd := exec.Command("ls", "/")
	out, _ := cmd.Output()
	fileNames := string(out)
	var mds []models.Metadata
	for _, fn := range strings.Split(fileNames, "\n") {
		if fn != "" {
			currMd, err := s.readMetaData(id, path+"/"+fn)
			if err != nil {
				return nil, err
			}
			mds = append(mds, *currMd)
		}
	}

	return mds, nil
}

func (s *Server) readMetaData(id, path string) (*models.Metadata, error) {
	log.Println("opening meta file")
	fm, err := os.Open(s.root + path + MetaFileSuffix)
	if err != nil {
		log.Println("unable to read %v", path+MetaFileSuffix)
		return nil, err
	}

	defer fm.Close()

	byteValue, err := ioutil.ReadAll(fm)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var md *models.Metadata
	json.Unmarshal(byteValue, &md)

	log.Println("checking for permission")
	if !slices.Contains(md.AllowList, id) {
		log.Println("The id %v does not have permission to read the file %v", id, path)
		return nil, errors.New("Permission Denied")
	}

	return md, nil
}
