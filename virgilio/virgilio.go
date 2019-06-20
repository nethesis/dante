package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nethesis/dante/virgilio/apis"
	"github.com/nethesis/dante/virgilio/configuration"
)

func main() {
	// initialize configuration using environment variables
	configuration.Init()

	router := gin.Default()

	// cors
	corsConf := cors.DefaultConfig()
	if len(configuration.Config.Virgilio.CorsAllowOrigins) > 0 {
		corsConf.AllowOrigins = configuration.Config.Virgilio.CorsAllowOrigins
	}
	router.Use(cors.New(corsConf))

	router.GET("/widget/:widgetName", apis.ReadWidget)

	// listen on default free port
	router.Run(":8081") // listen and serve
}
