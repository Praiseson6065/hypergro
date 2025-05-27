package main

import (
	_ "Praiseson6065/Hypergro-assign/config"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	err := Server()
	log.Printf("Server is starting on 8000")

	if err != nil {
		log.Fatal(err)
	}
}
