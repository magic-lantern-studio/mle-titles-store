// Declare package.
package main

// Import required packages.
import (
    "net/http"
    "github.com/gin-gonic/gin"
)

var DB = make(map[string]string)

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
        value, ok := DB[user]
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
            DB[user] = json.Value
            context.JSON(http.StatusOK, gin.H{"status": "ok"})
        }
    })

    return router
}

func main() {
    router := setupRouter()
    // Listen and Server in 0.0.0.0:8080
    router.Run(":8080")
}
