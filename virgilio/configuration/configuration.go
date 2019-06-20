package configuration

import (
	"os"
	"strings"
)

type Configuration struct {
	Ciacco struct {
		MinersDirectory string
		OutputDirectory string
	}
	Virgilio struct {
		CorsAllowOrigins []string
		MaxDays          int
	}
}

var Config = Configuration{}

func Init() {
	defaultCors := make([]string, 1)
	defaultCors[0] = "*"

	Config.Ciacco.OutputDirectory = os.Getenv("CIACCO_OUTPUT_DIR")
	Config.Ciacco.MinersDirectory = os.Getenv("CIACCO_MINERS_DIR")

	if os.Getenv("VIRGILIO_ALLOW_ORIGIN") != "" {
		// multiple origins should be separated by space
		Config.Virgilio.CorsAllowOrigins = strings.Split(os.Getenv("VIRGILIO_ALLOW_ORIGIN"), " ")
	}
	Config.Virgilio.MaxDays = 366
}
