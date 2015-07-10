package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HandleUserNew(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	next := r.URL.Query().Get("next")
	RenderTemplate(w, r, "users/new", map[string]interface{}{
		"Next": next,
	})
}

func HandleUserCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//process creating a user
	user, err := NewUser(
		r.FormValue("username"),
		r.FormValue("email"),
		r.FormValue("password"),
	)
	next := r.FormValue("next")
	if err != nil {
		if IsValidationError(err) {
			RenderTemplate(w, r, "users/new", map[string]interface{}{
				"Error": err.Error(),
				"User":  user,
				"Next":  next,
			})
			return
		}
		panic(err)
	}
	err = globalUserStore.Save(user)
	if err != nil {
		panic(err)
		return
	}

	//create a new session
	session := NewSession(w)
	session.UserID = user.ID
	err = globalSessionStore.Save(session)
	if err != nil {
		panic(err)
	}
	if next == "" {
		next = "/"
	}
	http.Redirect(w, r, next+"?flash=User+created&type=success", http.StatusFound)
}

func HandleUserEdit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := RequestUser(r)
	RenderTemplate(w, r, "users/edit", map[string]interface{}{
		"User": user,
	})
}
func HandleUserUpdate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	curUser := RequestUser(r)
	email := r.FormValue("email")
	curPwd := r.FormValue("currentPassword")
	newPwd := r.FormValue("newPassword")

	user, err := UpdateUser(curUser, email, curPwd, newPwd)
	if err != nil {
		if IsValidationError(err) {
			RenderTemplate(w, r, "users/edit", map[string]interface{}{
				"Error": err.Error(),
				"User":  user,
			})
			return
		}
		panic(err)
	}
	err = globalUserStore.Save(*curUser)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/account?flash=User+updated&type=success", http.StatusFound)
}

func HandleUserShow(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("user_id: ", params.ByName("userID"))
	user, err := globalUserStore.Find(params.ByName("userID"))
	if err != nil {
		panic(err)
	}

	//404
	if user == nil {
		http.NotFound(w, r)
		return
	}

	images, err := globalImageStore.FindAllByUser(user, 0)
	if err != nil {
		panic(err)
	}

	RenderTemplate(w, r, "users/show", map[string]interface{}{
		"Images": images,
		"User":   user,
	})
}
