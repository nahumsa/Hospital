package routes

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nahumsa/hospital-management/src/db"
)

func init() {
	godotenv.Load()
}

// Logger creates the entry for logs
type Logger struct {
	Username string `json:"username" bson:"username"`
	Room     string `json:"room" bson:"room"`
	Cid      string `json:"cid" bson:"cid"`
	Occupied bool   `json:"occupied" bson:"occupied"`
}

// AddLog generates a Log Handler
func AddLog(c *gin.Context) {
	// Parameters
	userKey := os.Getenv("userKey")

	// get the user
	session := sessions.Default(c)
	userGet := session.Get(userKey)

	// convert interface{} to string
	user := fmt.Sprintf("%v", userGet)

	// Get the form
	cid := c.PostForm("cid")
	room := c.PostForm("room")
	occupied := false

	if cid != "" {
		occupied = true
	}

	// Check login credentials
	log := Logger{Username: user,
		Room:     room,
		Cid:      cid,
		Occupied: occupied,
	}

	// Create mongodb connection
	url := os.Getenv("mongoURL") + os.Getenv("mongoPort") + "/"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, _ := db.Connect(ctx, url)
	collection := client.Client.Database("logsDB").Collection("logs")

	// catch if there already exists an user
	_, err := collection.InsertOne(ctx, log)

	if err != nil {
		// Need to rework error name
		c.JSON(http.StatusBadRequest, gin.H{"database error": err})
		return
	}

	c.Redirect(http.StatusFound, "/private/log")
	// c.JSON(http.StatusOK, gin.H{"message": "Log added"})
}

// GetLog initialize the homepage HTML.
func GetLog(c *gin.Context) {
	Render(c, gin.H{}, "logs.html")
}
