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

package configuration

import (
	"os"
	"path"
	"strconv"
	"strings"
)

type Configuration struct {
	Ciacco struct {
		MinersDirectory string
		OutputDirectory string
	}
	Virgilio struct {
		StoreDirectory   string
		LayoutFile       string
		CorsAllowOrigins []string
		MaxDays          int
		MaxEntries       int
		Anonymize        bool
	}
	Beatrice struct {
		BaseDirectory string
	}
}

var Config = Configuration{}

func Init() {
	Config.Ciacco.OutputDirectory = os.Getenv("CIACCO_OUTPUT_DIR")
	Config.Ciacco.MinersDirectory = os.Getenv("CIACCO_MINERS_DIR")

	if os.Getenv("VIRGILIO_ALLOW_ORIGIN") != "" {
		// multiple origins should be separated by space
		Config.Virgilio.CorsAllowOrigins = strings.Split(os.Getenv("VIRGILIO_ALLOW_ORIGIN"), " ")
	}

	if os.Getenv("VIRGILIO_STORE_DIR") != "" {
		Config.Virgilio.StoreDirectory = os.Getenv("VIRGILIO_STORE_DIR")
	} else {
		Config.Virgilio.StoreDirectory = "./"
	}
	Config.Virgilio.LayoutFile = path.Join(Config.Virgilio.StoreDirectory, "layout.json")

	Config.Virgilio.MaxDays = 366
	if os.Getenv("VIRGILIO_MAX_DAYS") != "" {
		max, err := strconv.Atoi(os.Getenv("VIRGILIO_MAX_DAYS"))
		if err == nil {
			Config.Virgilio.MaxDays = max
		}
	}

	Config.Virgilio.MaxEntries = 10
	if os.Getenv("VIRGILIO_MAX_ENTRIES") != "" {
		max, err := strconv.Atoi(os.Getenv("VIRGILIO_MAX_ENTRIES"))
		if err == nil {
			Config.Virgilio.MaxEntries = max
		}
	}

	Config.Virgilio.Anonymize = false
	if os.Getenv("VIRGILIO_ANONYMIZE") != "" {
		b, err := strconv.ParseBool(os.Getenv("VIRGILIO_ANONYMIZE"))
		if err == nil {
			Config.Virgilio.Anonymize = b
		}
	}

	Config.Beatrice.BaseDirectory = os.Getenv("BEATRICE_BASE_DIR")
}
