package main

import (
	//"fmt"
	"main/api"
	//"os"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

func main(){
    router := gin.Default()
	config := cors.Config{
		Origins:        "*",
		Methods:        "GET, PUT, POST, DELETE",
		RequestHeaders: "Origin, Authorization, Content-Type",
		ExposedHeaders: "",
		MaxAge: 1 * time.Minute,
		Credentials: false,
		ValidateHeaders: false,
	}

	// Apply the middleware to the router (works on groups too)
	router.Use(cors.Middleware(config))
	router.Static("images","./uploaded/images")
	api.Setup(router)
	router.Run(":8081")
}