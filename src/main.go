package main

import (
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	r := setupRouter()
	r.Run()

}

// User struct in order to keep the first name, last name, user and password.
// user and password are used in order to login, and Firstname and LastName will
// be used in order to keep logs on the recordings
type User struct {
	FirstName string `json:"firstname" bson:"firstname"`
	LastName  string `json:"lastname" bson:"lastname"`
	User      string `json:"user" bson:"user"`
	Password  string `json:"email" bson:"email"`
}

var client *mongo.Client

// setupRouter creates the routing of the API, using Gin Gonic.
func setupRouter() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.LoadHTMLGlob("views/*")
	router.GET("/", getHomepage)
	return router

}

// getHash create a hash for the password in order to keep the hash in the
// database.
func getHash(password []byte) string {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	must(err)
	return string(hash)
}

var secretKey = []byte("secrets")

// GenerateJWT generates the JWT string for authentication
func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Println("Error in JWT token generation")
		return "", err
	}

	return tokenString, nil

}

// getHomepage initialize the homepage HTML.
func getHomepage(c *gin.Context) {
	render(c, gin.H{"title": "login"}, "homepage.html")
}

func postSignUp(c *gin.Context) {

}

// render helper function to render JSON, XML and HTML depending on
// the request.
func render(c *gin.Context, data gin.H, templateName string) {

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		// Respond with XML
		c.XML(http.StatusOK, data["payload"])
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, templateName, data)
	}

}

// must catches errors and prints with log.
func must(err error) {
	if err != nil {
		log.Println(err)
	}
}
