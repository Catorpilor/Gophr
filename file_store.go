package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type FileStore struct {
	filename string
}

func (ps *FileStore) NewStore() error {
	contents, err := ioutil.ReadFile(ps.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	err = json.Unmarshal(contents, ps)
	if err != nil {
		return err
	}
	return nil
}

func (ps *FileStore) WriteToFile() error {
	fmt.Println(ps)
	contents, err := json.MarshalIndent(ps, "", "  ")
	fmt.Println(string(contents))
	if err != nil {
		return err
	}
	return ioutil.WriteFile(ps.filename, contents, 0660)
}
