package db

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Archive represents the location and the data of a given archive
type Archive struct {
	Path string
	Data Data
}

// Data is just a list of guids that have already been processed for a given feed
type Data struct {
	Guid []string `json:"guid"`
}

// NewArchive will create a new Archive struct with the given path
func NewArchive(path string) Archive {
	a := Archive{}
	a.Path = path

	return a
}

// Read will refresh the data included in the archive from the set path
func (a *Archive) Read() {
	archiveData, err := ioutil.ReadFile(a.Path)
	if err != nil {
		return
	}

	err = json.Unmarshal(archiveData, &a.Data)
	if err != nil {
		return
	}
}

// Contains will return true, if the guid is already part of the archive
func (a *Archive) Contains(guid string) bool {
	for _, s := range a.Data.Guid {
		if guid == s {
			return true
		}
	}

	return false
}

// Add will add the guid to the archive. Keep in mind that this is only in memory until Persist() is called
func (a *Archive) Add(guid string) error {
	a.Data.Guid = append(a.Data.Guid, guid)

	return nil
}

// Persist will write the current data to the disk at the given path
func (a *Archive) Persist() error {
	path := filepath.Dir(a.Path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
	}

	data, err := json.Marshal(a.Data)
	if err != nil {
		return err
	}
	ioutil.WriteFile(a.Path, data, 0600)

	return nil
}
