package widgets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

// ParseWidget return a widget generic data from a file
// Return nil if the file can't be parsed
func ParseWidget(filePath string) map[string]interface{} {
	var widgetData map[string]interface{}

	widgetFile, err := os.Open(filePath)
	defer widgetFile.Close()
	if err != nil {
		return nil
	}
	bytes, err := ioutil.ReadAll(widgetFile)
	if err != nil {
		return nil
	}

	json.Unmarshal(bytes, &widgetData)
	if err != nil {
		return nil
	}

	return widgetData
}

// ParseLayout return a layout from a file
// Return nil if the file can't be parsed
func ParseLayout(filePath string) Layout {
	var layout Layout

	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return Layout{}
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return Layout{}
	}

	json.Unmarshal(bytes, &layout)
	if err != nil {
		return Layout{}
	}

	fmt.Print("all ok ->")
	fmt.Println(layout)

	return layout
}
