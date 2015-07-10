package main

import (
	"crypto/md5"
	"fmt"

	"code.google.com/p/crypto.go/bcrypt"
)

type User struct {
	ID             string
	Email          string
	HashedPassword string
	Username       string
}

const (
	hashCost       = 10
	passwordLength = 6
	userIDLength   = 16
)

func NewUser(username, email, password string) (User, error) {
	user := User{
		Email:    email,
		Username: username,
	}
	if username == "" {
		return user, errNoUsername
	}
	if email == "" {
		return user, errNoEmail
	}
	if password == "" {
		return user, errNoPassword
	}
	if len(password) < passwordLength {
		return user, errPasswordTooShort
	}
	//check if user exist
	existingUser, err := globalUserStore.FindByUsername(username)
	if err != nil {
		return user, err
	}
	if existingUser != nil {
		return user, errUsernameExists
	}
	existingUser, err = globalUserStore.FindByEmail(email)
	if err != nil {
		return user, err
	}
	if existingUser != nil {
		return user, errEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	user.HashedPassword = string(hashedPassword)
	user.ID = GenerateID("usr", userIDLength)
	return user, err
}

func FindUser(username, password string) (*User, error) {
	ret := &User{
		Username: username,
	}

	extuser, err := globalUserStore.FindByUsername(username)
	if err != nil {
		return ret, nil
	}
	if extuser == nil {
		return ret, errInfoIncorect
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(extuser.HashedPassword),
		[]byte(password),
	) != nil {
		return ret, errInfoIncorect
	}

	return extuser, nil
}

func UpdateUser(user *User, email, curPwd, newPwd string) (User, error) {
	ret := *user
	if email == "" {
		return ret, errNoEmail
	}
	ret.Email = email

	//check if email exists
	extUser, err := globalUserStore.FindByEmail(email)
	if err != nil {
		return ret, err
	}
	if extUser != nil && extUser.ID != user.ID {
		return ret, errEmailExists
	}

	//update email address
	user.Email = email

	if curPwd == "" {
		return ret, errNoPassword
	}

	if bcrypt.CompareHashAndPassword(
		[]byte(user.HashedPassword),
		[]byte(curPwd),
	) != nil {
		return ret, errPasswordMisMatch
	}

	if newPwd == "" {
		return ret, errNoPassword
	}

	if len(newPwd) < passwordLength {
		return ret, errPasswordTooShort
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(newPwd), hashCost)
	user.HashedPassword = string(hashedPwd)
	return ret, err
}

func (user *User) AvatarURL() string {
	return fmt.Sprintf(
		"//en.gravatar.com/avatar/%x?size=30",
		md5.Sum([]byte(user.Email)),
	)
}

func (user *User) ImagesRoute() string {
	return "/user/" + user.ID
}
