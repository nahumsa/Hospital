package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/sessions"
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
var userKey = "user"

// setupRouter creates the routing of the API, using Gin Gonic.
func setupRouter() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.LoadHTMLGlob("views/*")
	router.GET("/", getHomepage)

	// Create login handler
	store := sessions.NewCookieStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	router.POST("/login", Login)
	router.GET("/logout", Logout)

	// Create private handlers
	private := router.Group("/private")
	private.Use(AuthRequired)
	{
		private.GET("/user", userGet)
	}

	return router
}

// AuthRequired is a middleware to check the session
func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unathorized"})
		return
	}
	// Continue down the chain
	c.Next()
}

// Login generates a login validation Handler
func Login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Validate post form
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	// Check login credentials
	// TODO: Use database
	userID := "1"
	if username != "user1" || password != "test" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password or login"})
		return
	}

	// Save the name in the session
	session.Set(userKey, userID)

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save Session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Authentication Successful"})
}

// Logout generates a handler that logouts current session
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)

	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}

	session.Delete(userKey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save Session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func userGet(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)
	c.JSON(http.StatusOK, gin.H{"user": user})
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
