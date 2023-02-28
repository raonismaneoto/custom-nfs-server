package dht

import (
	Client "github.com/raonismaneoto/CustomDHT/core/client"
)

type DhtStorage struct {
	client   Client.Client
	nodeAddr string
}

func New() DhtStorage {
	s := DhtStorage{}
	s.client = *Client.New()
	return s
}

func (s DhtStorage) Save(id, path string, content <-chan []byte, errors chan<- error) {
	temp_content := make(chan []byte)
	temp_errors := make(chan error)

	go s.client.SaveAsync(s.nodeAddr, path, temp_content, temp_errors)

	for {
		select {
		case currContent, ok := <-content:
			if !ok {
				close(temp_content)
				close(errors)
				return
			}
			temp_content <- currContent
		case err, _ := <-temp_errors:
			close(temp_content)
			close(errors)
			if err != nil {
				errors <- err
			}
			return
		}

	}
}

func (s DhtStorage) Read(id, path string, content chan<- []byte, errors chan<- error) {
	temp_content := make(chan []byte)
	temp_errors := make(chan error)

	go s.client.QueryAsync(s.nodeAddr, path, temp_content, temp_errors)

	for {
		select {
		case currContent, ok := <-temp_content:
			if !ok {
				close(content)
				close(errors)
				return
			}
			content <- currContent
		case err, _ := <-temp_errors:
			close(content)
			close(errors)
			if err != nil {
				errors <- err
			}
			return
		}

	}
}
