package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HandleSessionNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	next := r.URL.Query().Get("next")
	fmt.Sprintf("1st loging with next: ", r.URL.Query())
	RenderTemplate(w, r, "sessions/new", map[string]interface{}{
		"Next": next,
	})
}

func HandleSessionCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	next := r.FormValue("next")
	user, err := FindUser(username, password)
	if err != nil {
		if IsValidationError(err) {
			RenderTemplate(w, r, "sessions/new", map[string]interface{}{
				"Error": err,
				"User":  user,
				"Next":  next,
			})
			return
		}
		panic(err)
	}

	sess := FindOrCreateSession(w, r)
	sess.UserID = user.ID
	err = globalSessionStore.Save(sess)
	if err != nil {
		panic(err)
	}

	if next == "" {
		next = "/"
	}
	http.Redirect(w, r, next+"?flash=Signed+in&type=info", http.StatusFound)
}

func HandleSessionDestory(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := RequestSession(r)
	if session != nil {
		err := globalSessionStore.Delete(session)
		if err != nil {
			panic(err)
		}
	}
	RenderTemplate(w, r, "sessions/destory", nil)
}
