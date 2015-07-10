package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type SessionStore interface {
	Find(string) (*Session, error)
	Save(*Session) error
	Delete(*Session) error
}

type FileSessionStore struct {
	filename string
	//FileStore
	Sessions map[string]Session
}

var globalSessionStore SessionStore

/*
func init() {
	store, err := NewFileSessionStore("./data/sessions.json")
	if err != nil {
		panic(fmt.Errorf("Error creating session store: %s", err))
	}
	globalSessionStore = store
}
*/

func NewFileSessionStore(filename string) (*FileSessionStore, error) {
	store := &FileSessionStore{
		//FileStore{filename},
		//map[string]Session{},
		Sessions: map[string]Session{},
		filename: filename,
	}
	//err := store.NewStore()

	contents, err := ioutil.ReadFile(filename)

	if err != nil {
		//if it's a file not exist,that's ok
		if os.IsNotExist(err) {
			return store, nil
		}
		return nil, err
	}

	err = json.Unmarshal(contents, store)

	if err != nil {
		return nil, err
	}
	fmt.Println("newsessionStore: ", store)

	return store, err
}

func (store *FileSessionStore) Save(session *Session) error {
	store.Sessions[session.ID] = *session

	contents, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(store.filename, contents, 0660)

	//return store.WriteToFile()
}

func (store *FileSessionStore) Find(id string) (*Session, error) {
	sess, ok := store.Sessions[id]
	if ok {
		return &sess, nil
	}
	return nil, nil
}

func (store *FileSessionStore) Delete(session *Session) error {
	delete(store.Sessions, session.ID)

	contents, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(store.filename, contents, 0660)

	//return store.WriteToFile()
}
