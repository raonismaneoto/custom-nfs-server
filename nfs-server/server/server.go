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
	fm, err := os.OpenFile(s.root+path+MetaFileSuffix, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
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

func (s *Server) Read(id, path string, content chan<- []byte, errors chan<- error) {
	//check if id is allowed to access its content
	log.Println("going to read: ", path)
	if _, err := s.readMetaData(id, path); err != nil {
		errors <- err
		close(errors)
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
				return
			}
			content <- currContent
		case err, _ := <-child_errors:
			if err != nil {
				errors <- err
			}
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

	resp, err := s.getMetadataRecursively(md)
	log.Println(resp)
	return resp, err
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

	fm, err := os.OpenFile(s.root+path+MetaFileSuffix, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
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
	fm, err := os.OpenFile(s.root+path+MetaFileSuffix, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println("unable to open/create %v", path+MetaFileSuffix, err.Error())
		return err
	}

	err = s.saveMetaData(id, path, content, fm)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return s.storage.Save(id, path, content)
}

func (s *Server) Rm(id, path string) error {
	_, err := s.readMetaData(id, path)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	f, err := os.Open(s.root + path)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	var metapath string
	stat, err := f.Stat()
	if stat.IsDir() {
		if path[len(path)-1:] == "/" {
			metapath = path + MetaFileSuffix
		} else {
			metapath = path + "/" + MetaFileSuffix
		}
	} else {
		metapath = path + MetaFileSuffix
	}

	err = os.RemoveAll(s.root + metapath)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return s.storage.Rm(id, path)
}

func (s *Server) Mkdir(id, path string) error {
	err := os.Mkdir(s.root+path, 0777)
	if err != nil {
		log.Println("unable to create ", s.root+path)
		return err
	}

	fm, err := os.OpenFile(s.root+path+"/"+MetaFileSuffix, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println("unable to open/create %v", s.root+path+"/"+MetaFileSuffix)
		return err
	}

	err = s.saveMetaData(id, path, []byte{}, fm)
	if err != nil {
		log.Println("unable to create meta file")
		return err
	}

	return nil
}

func (s *Server) readMetaData(id, path string) (*models.Metadata, error) {
	var filePath string
	if strings.Contains(path, "meta") {
		filePath = path[:len(path)-4]
	} else {
		filePath = path
	}

	log.Println("opening file")
	f, err := os.Open(s.root + filePath)
	if err != nil {
		log.Println("unable to read %v", filePath)
		return nil, err
	}

	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		log.Println("unable to get file stat")
		return nil, err
	}

	var metaPath string
	if stat.IsDir() {
		metaPath = s.root + filePath + "/" + MetaFileSuffix
	} else {
		metaPath = s.root + filePath + MetaFileSuffix
	}

	fm, err := os.Open(metaPath)
	if err != nil {
		log.Println("unable to read %v", metaPath)
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
	var metaPath string
	if len(content) == 0 {
		metaPath = path + "/" + MetaFileSuffix
	} else {
		metaPath = path + MetaFileSuffix
	}

	metadata := models.Metadata{OwnerID: id, Size: float64(len(content)),
		Dir: len(content) == 0, AllowList: []string{}, Path: metaPath}

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
	splitString := strings.Split(path, "/")
	parentPath := strings.Join(splitString[:len(splitString)-1], "/") + "/"
	log.Println("parent path: ", parentPath)
	parentMd, err := s.readMetaData(id, parentPath)
	if err != nil {
		// create the dir meta file
		parentMd = &models.Metadata{OwnerID: id, Dir: true, AllowList: []string{}, Path: parentPath + MetaFileSuffix}
		log.Println(err.Error())
	}

	if slices.Contains(parentMd.Children, metadata.Path) {
		return nil
	}
	parentMd.Children = append(parentMd.Children, metadata.Path)
	mParentMd, err := json.Marshal(parentMd)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	fpm, err := os.OpenFile(s.root+parentPath+MetaFileSuffix, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer fpm.Close()
	log.Println("going to wirte to: ", s.root+parentPath+MetaFileSuffix)
	if _, err := fpm.Write(mParentMd); err != nil {
		log.Println("unable to write to %v", s.root+parentPath+MetaFileSuffix)
		return err
	}

	return nil
}

func (s *Server) getMetadataRecursively(md *models.Metadata) ([]models.Metadata, error) {
	var mds []models.Metadata
	mds = append(mds, *md)

	for _, childMdPath := range md.Children {
		childMd, err := s.readMetaData(md.OwnerID, childMdPath)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		if childMd.Dir {
			dirMds, err := s.getMetadataRecursively(childMd)
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
