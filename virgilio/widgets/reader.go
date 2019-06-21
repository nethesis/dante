package widgets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
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

	return layout
}

func ReadDefaultLayout() Layout {
	var widgets []Widget

	// find newest output directory named as day of the month
	libRegEx, e := regexp.Compile("^\\d\\d$")
	if e != nil {
		return Layout{}
	}

	modTime := time.Unix(0, 0)
	var newest string
	e = filepath.Walk(configuration.Config.Ciacco.OutputDirectory, func(path string, info os.FileInfo, err error) error {
		if err == nil && libRegEx.MatchString(info.Name()) {
			if info.Mode().IsDir() {
				if info.ModTime().After(modTime) {
					newest = path
				}
			}
		}

		return nil
	})
	if e != nil {
		return Layout{}
	}

	// read all widgets and generates the default layout
	i := 0
	files, err := ioutil.ReadDir(newest)
	if err != nil {
		return Layout{}
	}
	for _, f := range files {
		var w Widget
		if !f.IsDir() {
			obj := ParseWidget(path.Join(newest, f.Name()))
			w.Type = obj["type"].(string)
			w.Id = obj["minerId"].(string)
			w.I = i
			w.Y = i
			if i%2 == 0 {
				w.X = 0
			} else {
				w.X = 6
			}
			widgets = append(widgets, w)
			i++
		}
	}

	return Layout{widgets}
}

func ReadLayout() Layout {
	_, err := os.Stat(configuration.Config.Virgilio.LayoutFile)
	if os.IsNotExist(err) {
		return ReadDefaultLayout()
	}

	return ParseLayout(configuration.Config.Virgilio.LayoutFile)
}
