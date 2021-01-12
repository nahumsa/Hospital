package main

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/nahumsa/hospital-management/src/routes"
)

func main() {
	// Setup router
	r := setupRouter()
	r.Run()
}

// setupRouter creates the routing of the API, using Gin Gonic.
func setupRouter() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)
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
	}

	return router
}

var userKey = "user"

func userGet(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)
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
