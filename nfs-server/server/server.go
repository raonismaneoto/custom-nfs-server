package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/raonismaneoto/custom-nfs-server/nfs-server/models"
	"github.com/raonismaneoto/custom-nfs-server/nfs-server/storage"
	"golang.org/x/exp/slices"
)

const MetaFileSuffix string = "meta"

type Server struct {
	root    string
	storage storage.Storage
}

func New() *Server {
	root := os.Getenv("ROOT_FOLDER")
	if _, err := os.Stat(root); err != nil {
		os.Mkdir(root, 0777)
	}
	sType := os.Getenv("STORAGE_TYPE")
	return &Server{
		root:    root,
		storage: storage.New(sType),
	}
}

func (s *Server) SaveAsync(id, path string, content <-chan []byte, errors chan<- error) {
	log.Println("Save call received.")
	log.Println("saving")
	temp_content := make(chan []byte)
	child_errors := make(chan error)

	// open metadata file
	fm, err := os.OpenFile(s.root+path+MetaFileSuffix, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("unable to open/create %v", path+MetaFileSuffix)
		errors <- err
		close(errors)
		close(temp_content)
	}

	defer fm.Close()

	go s.storage.SaveAsync(id, path, temp_content, child_errors)
	for {
		select {
		case currContent, ok := <-content:
			if !ok {
				close(temp_content)
				close(errors)
				break
			}

			err := s.saveMetaData(id, path, currContent, fm)
			if err != nil {
				log.Println(err.Error())
				errors <- err
				break
			}

			temp_content <- currContent
		case err, _ := <-child_errors:
			if err != nil {
				errors <- err
			}
			close(temp_content)
			close(errors)
			break
		}

	}
}

func (s *Server) Read(id, path string, content chan []byte, errors chan error) {
	//check if id is allowed to access its content
	if _, err := s.readMetaData(id, path); err != nil {
		errors <- err
		close(errors)
		close(content)
		return
	}

	temp_content := make(chan []byte)
	child_errors := make(chan error)

	go s.storage.Read(id, path, temp_content, child_errors)
	for {
		select {
		case currContent, ok := <-temp_content:
			if !ok {
				close(content)
				close(errors)
				return
			}
			content <- currContent
		case err, _ := <-child_errors:
			if err != nil {
				errors <- err
			}
			close(content)
			close(errors)
			return
		}
	}
}

func (s *Server) GetMetaData(id, path string) ([]models.Metadata, error) {
	log.Println("going go call read metadata")
	md, err := s.readMetaData(id, path)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	if !md.Dir {
		return []models.Metadata{*md}, nil
	}

	return s.getDirMetadata(md)
}

func (s *Server) Chpem(ownerId, user, path, op string) error {
	md, err := s.readMetaData(ownerId, path)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if op == "add" {
		if !slices.Contains(md.AllowList, user) {
			md.AllowList = append(md.AllowList, user)
		}
	} else if op == "rm" {
		idx := slices.Index(md.AllowList, user)
		if idx != -1 {
			md.AllowList = slices.Delete(md.AllowList, idx, idx)
		}
	}

	fm, err := os.OpenFile(s.root+path+MetaFileSuffix, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("unable to open/create %v", path+MetaFileSuffix)
	}
	defer fm.Close()

	mMd, err := json.Marshal(md)
	if _, err := fm.Write(mMd); err != nil {
		log.Println("unable to write to %v", path+MetaFileSuffix)
		return err
	}

	return nil
}

func (s *Server) Save(id, path string, content []byte) error {
	return s.storage.Save(id, path, content)
}

func (s *Server) Mkdir(id, path string) error {
	return os.Mkdir(s.root+path, 0644)
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
	log.Println("id: " + id)
	if !slices.Contains(md.AllowList, id) {
		log.Println("The id %v does not have permission to read the file %v", id, path)
		return nil, errors.New("Permission Denied")
	}

	return md, nil
}

func (s *Server) saveMetaData(id, path string, content []byte, fm *os.File) error {
	metadata := models.Metadata{OwnerID: id, Size: float64(len(content)),
		Dir: len(content) == 0, AllowList: []string{}, Path: path + MetaFileSuffix}

	byteValue, err := ioutil.ReadAll(fm)

	if err == nil && len(byteValue) > 0 {
		err := json.Unmarshal(byteValue, &metadata)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		metadata.Size = metadata.Size + float64(len(content))
	}

	if id == metadata.OwnerID && !slices.Contains(metadata.AllowList, id) {
		metadata.AllowList = append(metadata.AllowList, id)
	} else if id != metadata.OwnerID && !slices.Contains(metadata.AllowList, id) {
		return errors.New("the user does not have permission to do such operation")
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

	// save as child
	splitString := strings.Split(s.root+path, "/")
	parentPath := strings.Join(splitString[:len(splitString)-1], "/") + "/" + MetaFileSuffix
	log.Println("parent path: ", parentPath)
	parentMd, err := s.readMetaData(id, "")
	if err != nil {
		// create the dir meta file
		parentMd = &models.Metadata{OwnerID: id, Dir: true, AllowList: []string{}, Path: parentPath}
		log.Println(err.Error())
	}

	if slices.Contains(parentMd.Children, &metadata) {
		return nil
	}
	parentMd.Children = append(parentMd.Children, &metadata)
	mParentMd, err := json.Marshal(parentMd)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	fpm, err := os.OpenFile(parentPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer fpm.Close()
	log.Println("going to wirte to: ", parentPath)
	if _, err := fpm.Write(mParentMd); err != nil {
		log.Println("unable to write to %v", path+MetaFileSuffix)
		return err
	}

	return nil
}

func (s *Server) getDirMetadata(md *models.Metadata) ([]models.Metadata, error) {
	if !md.Dir {
		return nil, errors.New("the file is not a dir")
	}

	var mds []models.Metadata
	mds = append(mds, *md)

	for _, childMd := range md.Children {
		if childMd.Dir {
			dirMds, err := s.getDirMetadata(childMd)
			if err != nil {
				log.Println("error getting mds recursive. Err: ", err.Error())
				return nil, err
			}
			mds = append(mds, dirMds...)
		} else {
			mds = append(mds, *childMd)
		}
	}

	return mds, nil
}
