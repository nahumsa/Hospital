package main

import (
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nahumsa/hospital-management/src/routes"
)

func init() {
	godotenv.Load()
}

func main() {
	// Setup router
	r := setupRouter()
	r.Run(":8000")
}

// setupRouter creates the routing of the API, using Gin Gonic.
func setupRouter() *gin.Engine {
	router := gin.New()
	router.LoadHTMLGlob("views/*")
	router.GET("/", routes.GetHomepage)

	// Create login handler
	store := sessions.NewCookieStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	router.POST("/login", routes.Login)
	router.GET("/logout", routes.Logout)
	// Signup handler
	router.POST("/signup", routes.Signup)

	// Create private handlers
	private := router.Group("/private")
	private.Use(routes.AuthRequired)
	{
		private.GET("/user", userGet)
		private.GET("/log", routes.GetLog)
		private.POST("/addlog", routes.AddLog)
	}

	return router
}

func userGet(c *gin.Context) {
	// Parameters
	userKey := os.Getenv("userKey")

	// get the user
	session := sessions.Default(c)
	user := session.Get(userKey)
	log.Println(reflect.TypeOf(user))
	c.JSON(http.StatusOK, gin.H{"user": user})
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
