package routes

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nahumsa/hospital-management/src/db"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	godotenv.Load()
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
	user := User{Username: username, Password: password}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create mongodb connection
	// port := "27017"
	// url := "mongodb://127.0.0.1:" + port + "/"
	port := os.Getenv("mongoPort")
	url := os.Getenv("mongoURL") + port + "/"
	client, _ := db.Connect(ctx, url)
	collection := client.Client.Database("loginDB").Collection("user")

	// catch if there already exists an user
	var userDB User
	err := collection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&userDB)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password or login"})
		return
	}

	// Compare passwords
	userPassword := []byte(user.Password)
	dbPassword := []byte(userDB.Password)

	passErr := bcrypt.CompareHashAndPassword(dbPassword, userPassword)

	if passErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Save the name in the session
	session.Set(userKey, user.Username)

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

// GetHomepage initialize the homepage HTML.
func GetHomepage(c *gin.Context) {
	Render(c, gin.H{"title": "login"}, "homepage.html")
}

// Render helper function to render JSON, XML and HTML depending on
// the request.
func Render(c *gin.Context, data gin.H, templateName string) {

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

// Signup generates a signup handler
func Signup(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	secretkey := c.PostForm("secretkey")

	// Validate secret key (Need to create env variable)
	if secretkey != "secret" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Secret Key"})
		return
	}

	// Validate post form
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}
	user := User{Username: username, Password: password}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create mongodb connection
	url := os.Getenv("mongoURL") + os.Getenv("mongoPort") + "/"

	client, _ := db.Connect(ctx, url)
	collection := client.Client.Database("loginDB").Collection("user")

	// catch if there already exists an user
	var userDB User
	err := collection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&userDB)
	emptyUser := User{}
	if userDB != emptyUser {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already in use"})
		return
	}

	// Add the user to the database
	user.Password = hashPassword([]byte(user.Password))
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		// Need to rework error name
		c.JSON(http.StatusBadRequest, gin.H{"database error": err})
		return
	}

	c.Redirect(http.StatusFound, "/")
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

// User struct in order to keep the first name, last name, user and password.
// user and password are used in order to login, and Firstname and LastName will
// be used in order to keep logs on the recordings
type User struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

// hashPassword create a hash for the password in order to keep the hash in the
// database.
func hashPassword(password []byte) string {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

var userKey = "user"
