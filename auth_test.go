package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRequestNewImageUnauthenticated(t *testing.T) {
	req, _ := http.NewRequest("GET", "/images/new", nil)
	recoder := httptest.NewRecorder()

	app := NewApp()
	app.ServeHTTP(recoder, req)

	if recoder.Code != http.StatusFound {
		t.Error("Expected a redirect code, but got ", recoder.Code)
	}
}

func TestRequestNewImageAuthenticated(t *testing.T) {
	oldUserStore := globalUserStore
	defer func() {
		globalUserStore = oldUserStore
	}()

	globalUserStore = &MockUserStore{
		findUser: &User{},
	}

	expiry := time.Now().Add(time.Hour)

	oldSessionStore := globalSessionStore
	defer func() {
		globalSessionStore = oldSessionStore
	}()

	globalSessionStore = &MockSessionStore{
		Session: &Session{
			ID:     "sess_123131321312",
			UserID: "usr_123132132132132131",
			Expiry: expiry,
		},
	}

	//create a cookie
	authCookie := &http.Cookie{
		Name:    sessionCookieName,
		Value:   "sess_123131321312",
		Expires: expiry,
	}

	req, _ := http.NewRequest("GET", "/images/new", nil)
	req.AddCookie(authCookie)

	recoder := httptest.NewRecorder()
	app := NewApp()

	app.ServeHTTP(recoder, req)

	if recoder.Code != http.StatusOK {
		t.Error("Expected a OK code, but got ", recoder.Code)
	}
}

func BenchmarkRequestNewImageUnauthenticated(b *testing.B) {
	req, _ := http.NewRequest("GET", "/images/new", nil)

	recorder := httptest.NewRecorder()

	app := NewApp()

	for i := 0; i < b.N; i++ {
		app.ServeHTTP(recorder, req)
	}
}
