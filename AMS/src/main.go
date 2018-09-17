// Declare package.
package main

// Import required packages.
import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// The administrative database.
var adminDB = make(map[string]string)

// The title database.
var titleDB = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color.
	// gin.DisableConsoleColor()
	router := gin.Default()

	// Ping test.
	router.GET("/ping", func(context *gin.Context) {
		context.String(http.StatusOK, "pong")
	})

	// Get user value.
	router.GET("/user/:name", func(context *gin.Context) {
		user := context.Params.ByName("name")
		value, ok := adminDB[user]
		if ok {
			context.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			context.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same as:
	// authorized := router.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//      "foo":  "bar",
	//      "manu": "123",
	//}))
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	authorized.POST("admin", func(context *gin.Context) {
		user := context.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if context.Bind(&json) == nil {
			adminDB[user] = json.Value
			context.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	router.GET("/v1/titles", func(context *gin.Context) {
		// Return list of titles.
		for key, value := range titleDB {
			fmt.Println("Key:", key, "Value:", value)
		}
	})

	router.GET("/v1/titles/:title_id", func(context *gin.Context) {
		titleId := context.Params.ByName("title_id")
		value, titleExists := titleDB[titleId]
		if titleExists {
			context.JSON(http.StatusOK, gin.H{"title": titleId, "value": value})
		} else {
			context.JSON(http.StatusOK, gin.H{"title": titleId, "status": "invalid title"})
		}
	})

	router.GET("/v1/titles/:title_id/workprints", func(context *gin.Context) {
		titleId := context.Params.ByName("title_id")
		value, titleExists := titleDB[titleId]
		dwpName := context.Query("name")
		if dwpName == "" {
			fmt.Println("Digital Workprint: not specified")
		} else {
			fmt.Println("Digital Workprint: ", dwpName)
		}
		if titleExists {
			context.JSON(http.StatusOK, gin.H{"title": titleId, "value": value, "workprints": "workprints exist"})
		} else {
			context.JSON(http.StatusOK, gin.H{"title": titleId, "status": "invalid title", "workprints": "no workprints"})
		}
	})

	return router
}

func main() {
	router := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	router.Run(":8080")
}
