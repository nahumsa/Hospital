package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func performRequest(r http.Handler, method, path string, w *httptest.ResponseRecorder) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	r.ServeHTTP(w, req)
	return w
}

func TestHomepageStatus(t *testing.T) {
	// Setup router
	router := setupRouter()
	// Parameters of the test
	wantCode := http.StatusOK
	method := "GET"
	url := "/"

	// Test
	// Perform a GET request with that handler.
	w := httptest.NewRecorder()
	w = performRequest(router, method, url, w)

	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, wantCode, w.Code)
}

func performPost(r http.Handler, method, path string, payload io.Reader) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.ServeHTTP(w, req)
	return w
}

func TestLogin(t *testing.T) {
	t.Run("LoginSuccessful", LoginSuccessful)
	t.Run("LoginFailed", LoginFailed)
	t.Run("LoginNoUser", LoginNoForm)
}

func TestLogout(t *testing.T) {
	t.Run("LogoutNoUser", LogoutNoUser)
	// t.Run("LogoutUser", LogoutUser)
}

func LoginSuccessful(t *testing.T) {
	router := setupRouter()

	wantBody := `{"message":"Authentication Successful"}`
	wantCode := http.StatusOK
	method := "POST"
	url := "/login"
	payload := strings.NewReader("username=user1&password=test")

	w := performPost(router, method, url, payload)
	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, wantCode, w.Code)

	// Check body
	body, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, wantBody, string(body))
}

func LoginFailed(t *testing.T) {
	router := setupRouter()

	wantBody := `{"error":"Invalid password or login"}`
	wantCode := http.StatusUnauthorized
	method := "POST"
	url := "/login"
	payload := strings.NewReader("username=user2&password=test")

	w := performPost(router, method, url, payload)

	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, wantCode, w.Code)

	// Check body
	body, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, wantBody, string(body))
}

func LoginNoForm(t *testing.T) {
	router := setupRouter()

	wantBody := `{"error":"Parameters can't be empty"}`
	wantCode := http.StatusBadRequest
	method := "POST"
	url := "/login"
	payload := strings.NewReader("username=&password=")

	w := performPost(router, method, url, payload)

	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, wantCode, w.Code)

	// Check body
	body, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, wantBody, string(body))
}

// func LogoutUser(t *testing.T) {
// 	router := setupRouter()

// 	wantBody := `{"message":"Successfully logged out"}`
// 	wantCode := http.StatusOK
// 	method := "GET"
// 	url := "/logout"

// 	// Login first
// 	payload := strings.NewReader("username=user1&password=test")
// 	w := performPost(router, "POST", "/login", payload)

// 	cookie := w.HeaderMap.Get("Set-Cookie")
// 	w = httptest.NewRecorder()
// 	http.SetCookie(w, &http.Cookie{Name: "mysession", Value: cookie})
// 	w = performRequest(router, method, url, w)

// 	// Assert we encoded correctly,
// 	// the request gives a 200
// 	assert.Equal(t, wantCode, w.Code)

// 	// Check body
// 	body, _ := ioutil.ReadAll(w.Body)
// 	assert.Equal(t, wantBody, string(body))
// }

func LogoutNoUser(t *testing.T) {
	// Setup router
	router := setupRouter()
	// Parameters of the test
	wantCode := http.StatusBadRequest
	method := "GET"
	url := "/logout"

	// Test
	// Perform a GET request with that handler.
	w := httptest.NewRecorder()
	w = performRequest(router, method, url, w)

	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, wantCode, w.Code)
}
