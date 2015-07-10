package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func init() {
	//assign userStore
	store, err := NewFileUserStore("./data/users.json")
	if err != nil {
		panic(fmt.Errorf("Error creating user store: %s", err))
	}
	globalUserStore = store

	//assign sessionStore
	sessionStore, err := NewFileSessionStore("./data/sessions.json")
	if err != nil {
		panic(fmt.Errorf("Error creating session store: %s", err))
	}
	globalSessionStore = sessionStore

	//assign db
	db, err := NewMySQLDB("zhimaa:zhimaa@tcp(127.0.0.1:3306)/gopher?parseTime=true")
	if err != nil {
		panic(err)
	}
	globalMySQLDB = db

	//assign imagestore
	globalImageStore = NewDBImageStore()

}

func main() {
	log.Fatal(http.ListenAndServe(":3000", NewApp()))
	/*
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			RenderTemplate(w, r, "index/home", nil)
		})
		mux.Handle("/asserts/", http.StripPrefix("/asserts/", http.FileServer(http.Dir("asserts/"))))
		log.Fatal(http.ListenAndServe(":3000", mux))
	*/
}

func NewApp() Middleware {
	staticRouter := NewRouter()
	staticRouter.ServeFiles("/asserts/*filepath", http.Dir("asserts/"))
	staticRouter.ServeFiles("/im/*filepath", http.Dir("data/images/"))
	staticRouter.Handle("GET", "/", HandleHome)
	staticRouter.Handle("GET", "/register", HandleUserNew)
	staticRouter.Handle("POST", "/register", HandleUserCreate)
	staticRouter.Handle("GET", "/login", HandleSessionNew)
	staticRouter.Handle("GET", "/image/:imageID", HandleImageShow)
	staticRouter.Handle("POST", "/login", HandleSessionCreate)

	secureRouter := NewRouter()
	secureRouter.Handle("GET", "/sign-out", HandleSessionDestory)
	secureRouter.Handle("GET", "/account", HandleUserEdit)
	secureRouter.Handle("GET", "/user/:userID", HandleUserShow)
	secureRouter.Handle("POST", "/account", HandleUserUpdate)
	secureRouter.Handle("GET", "/images/new", HandleImageNew)
	secureRouter.Handle("POST", "/images/new", HandleImageCreate)

	middleware := Middleware{}
	middleware.Add(staticRouter)
	middleware.Add(http.HandlerFunc(RequireLogin))
	middleware.Add(secureRouter)
	return middleware
}

func NewRouter() *httprouter.Router {
	router := httprouter.New()
	var ct CustomType
	router.NotFound = ct
	return router
}
