package main

import "net/http"

type CustomType int

func (c CustomType) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
