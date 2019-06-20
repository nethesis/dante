package widgets

import (
	"fmt"
	"path"
	"strconv"
	"time"

	"github.com/nethesis/dante/virgilio/configuration"
)

func GetFileLists(widgetName string, startDate time.Time, deltaDays int) []string {
	// list of files to aggregate
	filePaths := make([]string, deltaDays+1)
	fileName := ""

	for i := 0; i <= deltaDays; i++ {
		date := startDate.AddDate(0, 0, i)
		year := strconv.Itoa(date.Year())
		month := fmt.Sprintf("%02d", int(date.Month()))
		day := fmt.Sprintf("%02d", date.Day())
		fileName = year + "/" + month + "/" + day + "/" + widgetName + ".json"
		fullPath := path.Join(configuration.Config.Ciacco.OutputDirectory, fileName)
		filePaths[i] = fullPath
	}
	return filePaths
}
