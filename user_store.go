package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type UserStore interface {
	Find(string) (*User, error)
	FindByEmail(string) (*User, error)
	FindByUsername(string) (*User, error)
	Save(User) error
}

type FileUserStore struct {
	filename string
	//FileStore
	Users       map[string]User
	UserByName  map[string]string
	UserByEmail map[string]string
}

var globalUserStore UserStore

/*
func init() {
	store, err := NewFileUserStore("./data/users.json")
	if err != nil {
		panic(fmt.Errorf("Error creating user store: %s", err))
	}
	globalUserStore = store
}
*/

func NewFileUserStore(filename string) (*FileUserStore, error) {
	store := &FileUserStore{
		//FileStore:   FileStore{filename},
		filename:    filename,
		Users:       map[string]User{},
		UserByEmail: map[string]string{},
		UserByName:  map[string]string{},
	}

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

	//err := store.NewStore()
	fmt.Println("userStore: ", store)
	return store, err
}

func (store FileUserStore) Save(user User) error {
	store.Users[user.ID] = user
	store.UserByEmail[user.Email] = user.ID
	store.UserByName[user.Username] = user.ID

	contents, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(store.filename, contents, 0660)

	//return store.WriteToFile()
}

func (store FileUserStore) Find(id string) (*User, error) {
	user, ok := store.Users[id]
	if ok {
		return &user, nil
	}
	return nil, nil
}

func (store FileUserStore) FindByUsername(username string) (*User, error) {
	if username == "" {
		return nil, nil
	}
	userID, ok := store.UserByName[username]
	if ok {
		user, ok := store.Users[userID]
		if ok {
			return &user, nil
		}
	}
	/*
		for _, user := range store.Users {
			if strings.ToLower(username) == strings.ToLower(user.Username) {
				return &user, nil
			}
		}*/
	return nil, nil
}

func (store FileUserStore) FindByEmail(email string) (*User, error) {
	if email == "" {
		return nil, nil
	}
	userID, ok := store.UserByEmail[email]
	if ok {
		user, ok := store.Users[userID]
		if ok {
			return &user, nil
		}
	}
	/*
		for _, user := range store.Users {
			if strings.ToLower(email) == strings.ToLower(user.Email) {
				return &user, nil
			}
		}
	*/
	return nil, nil
}
