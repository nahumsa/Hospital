package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := setupRouter()
	r.Run()
}

// setupRouter creates the routing of the Books API
func setupRouter() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.LoadHTMLGlob("views/*")
	router.GET("/", getHomepage)
	return router

}

func getHomepage(c *gin.Context) {
	c.HTML(http.StatusOK, "homepage.html", nil)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
