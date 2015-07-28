package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

const (
	AuthTokenUrl = "http://api.meican.loc/v2.1/accounts/authtoken"
)

func HandleHome(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	//display home page
	fmt.Println("Cookies: ", r.Cookies())
	remember, err := r.Cookie("remember")
	if err != nil || remember == nil {
		fmt.Println("no valid cookie found, err: ", err)
	}

	req, err := http.NewRequest("GET", AuthTokenUrl, nil)
	req.AddCookie(remember)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("failed at calling api, error: %s", err)
	}
	defer resp.Body.Close()

	code := resp.StatusCode
	if code != http.StatusOK {
		fmt.Errorf("failed at calling api, got status code: %d", code)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("failed at reading responseBody, error: %s", err)
	}
	var data map[string]interface{}
	if err = json.Unmarshal(body, &data); err != nil {
		fmt.Errorf("json.Umarshal failed, error: %s", err)
	}
	fmt.Println(data)
	uid, ok := data["userId"]
	if !ok {
		fmt.Errorf("couldn't get userId from api")
	}

	userId, err := strconv.ParseInt(uid.(string), 10, 64)
	if err != nil {
		fmt.Errorf("invalid userId: %s", uid.(string))
	}

	fmt.Println("Userid is :", userId)
	images, err := globalImageStore.FindAll(0)
	if err != nil {
		panic(err)
	}
	RenderTemplate(w, r, "index/home", map[string]interface{}{
		"Images": images,
	})
}
