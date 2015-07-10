package main

import (
	"crypto/rand"
	"fmt"
)

//Source string used when generating a random identifier
const IDSource = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

//Save the length in a const so we don't look it up each time
const IDSourceLength = byte(len(IDSource))

func GenerateID(prefix string, length int) string {
	//create an array with the right capacity
	id := make([]byte, length)

	//fill with random numbers
	rand.Read(id)

	//Replace each random number with an alphanumeric value
	for i, b := range id {
		id[i] = IDSource[b%IDSourceLength]
	}

	//return formatted id
	return fmt.Sprintf("%s_%s", prefix, string(id))
}
