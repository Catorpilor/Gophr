package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HandleImageNew(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	RenderTemplate(w, r, "images/new", nil)
}

func HandleImageCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.FormValue("url") != "" {
		//create image from url
		handleImageCreateFromURL(w, r)
		return
	}
	handleImageCreateFromFile(w, r)
}

func handleImageCreateFromURL(w http.ResponseWriter, r *http.Request) {
	user := RequestUser(r)

	image := NewImage(user)
	image.Description = r.FormValue("description")

	err := image.CreateFromURL(r.FormValue("url"))

	if err != nil {
		if IsValidationError(err) {
			RenderTemplate(w, r, "images/new", map[string]interface{}{
				"Error":    err,
				"ImageURL": r.FormValue("url"),
				"Image":    image,
			})
			return
		}
		panic(err)
	}

	http.Redirect(w, r, "/?flash=Image+Uploaded+Successfully&type=success", http.StatusFound)
}

func handleImageCreateFromFile(w http.ResponseWriter, r *http.Request) {
	user := RequestUser(r)

	image := NewImage(user)
	image.Description = r.FormValue("description")

	file, headers, err := r.FormFile("file")

	//no file uploaded
	if file == nil {
		RenderTemplate(w, r, "images/new", map[string]interface{}{
			"Error": errNoImage,
			"Image": image,
		})
		return
	}

	//file uploaded with an error
	if err != nil {
		panic(err)
	}

	defer file.Close()

	err = image.CreateFromFile(file, headers)
	if err != nil {
		RenderTemplate(w, r, "images/new", map[string]interface{}{
			"Error": errNoImage,
			"Image": image,
		})
		return
	}

	http.Redirect(w, r, "/?flash=Image+Uploaded+Successfully&type=success", http.StatusFound)
}

func HandleImageShow(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("imageID: ", params.ByName("imageID"))
	image, err := globalImageStore.Find(params.ByName("imageID"))
	if err != nil {
		panic(err)
	}

	//404
	if image == nil {
		http.NotFound(w, r)
		return
	}

	user, err := globalUserStore.Find(image.UserID)
	if err != nil {
		panic(err)
	}

	if user == nil {
		panic(fmt.Errorf("Couldn't find use %s", image.UserID))
	}

	RenderTemplate(w, r, "images/show", map[string]interface{}{
		"Image": image,
		"User":  user,
	})
}
