package db

import (
	"io/ioutil"
	"encoding/json"
	"os"
	"path/filepath"
)

type Archive struct {
	Path string
	Data Data
}

type Data struct {
	Guid []string `json:"guid"`
}

func NewArchive(path string) Archive {
	a := Archive{}
	a.Path = path

	return a
}

func (a *Archive) Read(){
	archiveData, err := ioutil.ReadFile(a.Path)
	if err != nil {
		return
	}

	err = json.Unmarshal(archiveData, &a.Data)
	if err != nil {
		return
	}
}

func (a *Archive) Contains(guid string) bool {
	for _, s := range a.Data.Guid {
		if guid == s {
			return true
		}
	}

	return false
}

func (a *Archive) Add(guid string) error {
	a.Data.Guid = append(a.Data.Guid, guid)

	return nil
}

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

