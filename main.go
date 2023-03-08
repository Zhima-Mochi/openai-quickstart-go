package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/Zhima-Mochi/openai-quickstart-go/api/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

var generateHandler *handlers.GenerateHandler

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("OPENAI_API_KEY")

	generateHandler = handlers.NewGenerateHandler(openai.NewClient(apiKey))
}

func main() {
	router := gin.Default()
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(c.Writer, nil)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	})

	apiRouter := router.Group("api")
	apiRouter.POST("/generate", generateHandler.TextCompletionHandler)

	router.Run()
}
