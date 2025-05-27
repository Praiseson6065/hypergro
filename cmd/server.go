package main

import (
	"Praiseson6065/Hypergro-assign/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Server() error {
	env := viper.GetString("ENVIRONMENT")
	port := viper.GetString(env + ".server.port")
	r := gin.New()
	r.Use(middleware.CORS())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Hello": "World",
		})
	})
	AuthRouter(r)
	ApiRouter(r) 
	err := r.Run(port)
	if err != nil {
		return err
	}
	return nil
}
