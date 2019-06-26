/*
 * Copyright (C) 2019 Nethesis S.r.l.
 * http://www.nethesis.it - info@nethesis.it
 *
 * This file is part of Dante project.
 *
 * Dante is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * Dante is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dante.  If not, see COPYING.
 */
package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nethesis/dante/virgilio/apis"
	"github.com/nethesis/dante/virgilio/configuration"
	"github.com/nethesis/dante/virgilio/utils"
)

func main() {
	// initialize configuration using environment variables
	configuration.Init()

	router := gin.Default()

	// cors
	corsConf := cors.DefaultConfig()

	if len(configuration.Config.Virgilio.CorsAllowOrigins) > 0 {
		if utils.ContainsString(configuration.Config.Virgilio.CorsAllowOrigins, "*") {
			corsConf.AllowAllOrigins = true
		} else {
			corsConf.AllowOrigins = configuration.Config.Virgilio.CorsAllowOrigins
		}
	} else {
		corsConf.AllowAllOrigins = true
	}
	router.Use(cors.New(corsConf))

	router.GET("/widget/:widgetName", apis.ReadWidget)
	router.GET("/miners", apis.ListMiners)
	router.GET("/layout", apis.GetLayout)
	router.POST("/layout", apis.SetLayout)
	router.DELETE("/layout", apis.DeleteLayout)
	router.GET("/lang/:langCode", apis.GetLang)

	// listen on default free port
	router.Run(os.Getenv("PORT"))
}
