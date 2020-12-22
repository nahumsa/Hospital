package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/appleboy/gofight/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	gin.SetMode(gin.TestMode)
	// Set up context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create mongodb connection
	port := "27017"
	url := "mongodb://localhost:" + port
	client, _ := mongo.Connect(ctx, options.Client().ApplyURI(url))
	_ = client
}

func TestHomepage(t *testing.T) {
	t.Run("HomepageStatus", HomepageStatus)
}

func TestSignUp(t *testing.T) {
	t.Run("SignUpSuccessful", SignUpSuccessful)
	t.Run("SignUpNoKey", SignUpNoKey)
	t.Run("SignUpUsername Already in Use", SignUpUsernameUse)
	t.Run("SignUpNoForm", SignUpNoForm)
}

func TestLogin(t *testing.T) {
	t.Run("LoginSuccessful", LoginSuccessful)
	t.Run("LoginFailed", LoginFailed)
	t.Run("LoginNoUser", LoginNoForm)
}

func TestLogout(t *testing.T) {
	t.Run("LogoutNoUser", LogoutNoUser)
	t.Run("LogoutUser", LogoutUser)
}

func TestPrivate(t *testing.T) {
	t.Run("PrivateNoAuth", PrivateNoAuth)
	t.Run("PrivateTestAuth", PrivateTestAuth)
}

func HomepageStatus(t *testing.T) {

	g := gofight.New()
	e := setupRouter()

	wantCode := http.StatusOK

	// Test without login
	g.GET("/").Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, wantCode, r.Code)
	})
}

func PrivateTestAuth(t *testing.T) {
	// Function adapted from https://github.com/Depado/gin-auth-example/blob/master/main_test.go
	g := gofight.New()
	e := setupRouter()
	wantStatus := http.StatusOK
	wantBody := `{"user":"1"}`

	var cookie string
	g.POST("/login").
		SetForm(gofight.H{"username": "user1", "password": "test"}).
		Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
			cookie = r.HeaderMap.Get("Set-Cookie")
			// Check if there is a cookie
			assert.NotZero(t, cookie)
		})

	g.GET("/private/user").SetHeader(gofight.H{"Cookie": cookie}).Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, wantStatus, r.Code)
		body, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, wantBody, string(body))
	})
}

func PrivateNoAuth(t *testing.T) {
	g := gofight.New()
	e := setupRouter()

	wantBody := `{"error":"unathorized"}`
	wantStatus := http.StatusUnauthorized

	g.GET("/private/user").Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, wantStatus, r.Code)
		body, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, wantBody, string(body))
	})
}

func SignUpSuccessful(t *testing.T) {
	g := gofight.New()
	e := setupRouter()

	wantBody := `{"message": "Signup Complete"}`
	wantStatus := http.StatusOK

	g.POST("/signup").
		SetForm(gofight.H{"username": "user1", "password": "test", "secretkey": "secret"}).
		Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, wantStatus, r.Code)
			body, _ := ioutil.ReadAll(r.Body)
			assert.Equal(t, wantBody, string(body))
		})
}

func SignUpNoKey(t *testing.T) {
	g := gofight.New()
	e := setupRouter()

	wantBody := `{"error":"Invalid Secret Key"}`
	wantStatus := http.StatusBadRequest

	g.POST("/signup").
		SetForm(gofight.H{"username": "user1", "password": "test", "secretkey": ""}).
		Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, wantStatus, r.Code)
			body, _ := ioutil.ReadAll(r.Body)
			assert.Equal(t, wantBody, string(body))
		})
}

func SignUpUsernameUse(t *testing.T) {
	g := gofight.New()
	e := setupRouter()

	wantBody := `{"error":"username already in use"}`
	wantStatus := http.StatusBadRequest

	g.POST("/signup").
		SetForm(gofight.H{"username": "user1", "password": "test", "secretkey": "secret"}).
		Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, wantStatus, r.Code)
			body, _ := ioutil.ReadAll(r.Body)
			assert.Equal(t, wantBody, string(body))
		})
}

func SignUpNoForm(t *testing.T) {
	g := gofight.New()
	e := setupRouter()

	wantBody := `{"error":"username already in use"}`
	wantStatus := http.StatusBadRequest

	g.POST("/signup").
		SetForm(gofight.H{"username": "", "password": "", "secretkey": "secret"}).
		Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, wantStatus, r.Code)
			body, _ := ioutil.ReadAll(r.Body)
			assert.Equal(t, wantBody, string(body))
		})
}

func LoginSuccessful(t *testing.T) {
	g := gofight.New()
	e := setupRouter()

	wantBody := `{"message":"Authentication Successful"}`
	wantStatus := http.StatusOK

	g.POST("/login").
		SetForm(gofight.H{"username": "user1", "password": "test"}).
		Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, wantStatus, r.Code)
			body, _ := ioutil.ReadAll(r.Body)
			assert.Equal(t, wantBody, string(body))
		})
}

func LoginFailed(t *testing.T) {
	g := gofight.New()
	e := setupRouter()

	wantBody := `{"error":"Invalid password or login"}`
	wantStatus := http.StatusUnauthorized

	g.POST("/login").
		SetForm(gofight.H{"username": "asae", "password": "wrr"}).
		Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, wantStatus, r.Code)
			body, _ := ioutil.ReadAll(r.Body)
			assert.Equal(t, wantBody, string(body))
		})
}

func LoginNoForm(t *testing.T) {
	g := gofight.New()
	e := setupRouter()

	wantBody := `{"error":"Parameters can't be empty"}`
	wantStatus := http.StatusBadRequest
	g.POST("/login").
		SetForm(gofight.H{"username": "", "password": ""}).
		Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, wantStatus, r.Code)
			body, _ := ioutil.ReadAll(r.Body)
			assert.Equal(t, wantBody, string(body))
		})
}

func LogoutUser(t *testing.T) {
	// Function adapted from https://github.com/Depado/gin-auth-example/blob/master/main_test.go
	g := gofight.New()
	e := setupRouter()
	wantStatus := http.StatusOK
	wantBody := `{"message":"Successfully logged out"}`

	var cookie string
	g.POST("/login").
		SetForm(gofight.H{"username": "user1", "password": "test"}).
		Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
			cookie = r.HeaderMap.Get("Set-Cookie")
			// Check if there is a cookie
			assert.NotZero(t, cookie)
		})

	g.GET("/logout").SetHeader(gofight.H{"Cookie": cookie}).Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, wantStatus, r.Code)
		body, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, wantBody, string(body))
	})
}

func LogoutNoUser(t *testing.T) {
	g := gofight.New()
	e := setupRouter()

	wantBody := `{"error":"Invalid session token"}`

	// Test without login
	g.GET("/logout").Run(e, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
		assert.Equal(t, http.StatusBadRequest, r.Code)
		body, _ := ioutil.ReadAll(r.Body)
		assert.Equal(t, wantBody, string(body))
	})
}
