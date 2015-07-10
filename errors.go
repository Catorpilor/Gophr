package main

import "errors"

type ValidationError error

var (
	errNoUsername       = ValidationError(errors.New("You must supply a username"))
	errNoEmail          = ValidationError(errors.New("You must supply an email"))
	errNoPassword       = ValidationError(errors.New("You must supply a password"))
	errPasswordTooShort = ValidationError(errors.New("Your password too short"))
	errUsernameExists   = ValidationError(errors.New("That username is taken"))
	errEmailExists      = ValidationError(errors.New("That email is taken"))
	errInfoIncorect     = ValidationError(errors.New("username or password wrong"))
	errPasswordMisMatch = ValidationError(errors.New("Password did not match"))
	errInvalidImageType = ValidationError(errors.New("Please upload only jpeg gif or png images"))
	errNoImage          = ValidationError(errors.New("Please select an image to upload"))
	errInvalidImageURL  = ValidationError(errors.New("Could not download image from url"))
)

func IsValidationError(err error) bool {
	_, ok := err.(ValidationError)
	return ok
}
