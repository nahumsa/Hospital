package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Set up context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	r := setupRouter()
	r.Run()
}

// User struct in order to keep the first name, last name, user and password.
// user and password are used in order to login, and Firstname and LastName will
// be used in order to keep logs on the recordings
type User struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

// getHash create a hash for the password in order to keep the hash in the
// database.
func hashPassword(password []byte) string {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
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
	// Signup handler
	router.POST("/signup", Signup)

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

	collection := client.Database("loginDB").Collection("user")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Signup Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Signup Complete"})

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

	collection := client.Database("loginDB").Collection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// catch if there already exists an user
	var userDB User
	err := collection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&userDB)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
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

func userGet(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// getHomepage initialize the homepage HTML.
func getHomepage(c *gin.Context) {
	render(c, gin.H{"title": "login"}, "homepage.html")
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

// var secretKey = []byte("secrets")

// GenerateJWT generates the JWT string for authentication
// func GenerateJWT() (string, error) {
// 	token := jwt.New(jwt.SigningMethodHS256)
// 	tokenString, err := token.SignedString(secretKey)
// 	if err != nil {
// 		log.Println("Error in JWT token generation")
// 		return "", err
// 	}

// 	return tokenString, nil

// }
