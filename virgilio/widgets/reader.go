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

package widgets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/nethesis/dante/virgilio/configuration"
	"github.com/nethesis/dante/virgilio/utils"
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
	topWidgets := []string{"df-chart-pie", "df-counter", "hardware-label", "hostname-label", "release-label"}
	typesOrderInLayout := []string{"label", "counter", "chart", "list", "table"}

	// sort widgets
	sort.Slice(files, func(i, j int) bool {
		fileNameI := strings.Split(files[i].Name(), ".")[0]
		fileNameJ := strings.Split(files[j].Name(), ".")[0]

		// some widgets stay at the top of the layout
		isTopWidgetI := utils.ContainsString(topWidgets, fileNameI)
		isTopWidgetJ := utils.ContainsString(topWidgets, fileNameJ)

		if isTopWidgetI && !isTopWidgetJ {
			return true
		}
		if isTopWidgetJ && !isTopWidgetI {
			return false
		}

		// split file name by dash or dot
		tokensI := strings.FieldsFunc(files[i].Name(), utils.IsDashOrDot)
		tokensJ := strings.FieldsFunc(files[j].Name(), utils.IsDashOrDot)
		// get widget types
		typeI := tokensI[1]
		typeJ := tokensJ[1]
		// sort types
		rankI := utils.IndexOf(typeI, typesOrderInLayout)
		rankJ := utils.IndexOf(typeJ, typesOrderInLayout)
		return rankI < rankJ
	})

	for _, f := range files {
		var w Widget

		// exclude directories and file names starting with "." (e.g. Vim swap files)
		if !f.IsDir() && string(f.Name()[0]) != "." {
			obj := ParseWidget(path.Join(newest, f.Name()))
			w.Type = obj["type"].(string)
			w.Id = obj["minerId"].(string)
			w.I = i
			w.Y = float64(i / 2)
			if i%2 == 0 {
				w.X = 0
			} else {
				w.X = 6
			}
			widgets = append(widgets, w)
			i++
		}
	}
	return Layout{widgets, true}
}

func ReadLayout() Layout {
	_, err := os.Stat(configuration.Config.Virgilio.LayoutFile)
	if os.IsNotExist(err) {
		return ReadDefaultLayout()
	}

	layout := ParseLayout(configuration.Config.Virgilio.LayoutFile)
	layout.Default = false

	return layout
}
