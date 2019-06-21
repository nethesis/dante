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

package apis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/nethesis/dante/virgilio/configuration"
	"github.com/nethesis/dante/virgilio/utils"
	"github.com/nethesis/dante/virgilio/widgets"
)

type miner struct {
	FullName string `json:"fullName"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}

// Validate query parameters
func validate(widgetName string, startDateString string, endDateString string) (string, time.Time, int) {
	deltaDays := 0
	startDate := time.Unix(0, 0)
	if widgetName == "" {
		return "widgetName is mandatory", startDate, deltaDays
	}

	if startDateString == "" {
		return "startDate is mandatory", startDate, deltaDays
	}
	if endDateString == "" {
		return "endDate is mandatory", startDate, deltaDays
	}

	startDateTokens := strings.Split(startDateString, "-")
	endDateTokens := strings.Split(endDateString, "-")
	if len(startDateTokens) < 3 {
		return "Malformed startDate", startDate, deltaDays
	}
	if len(endDateTokens) < 3 {
		return "Malformed endDate", startDate, deltaDays
	}

	startDateYear, startDateYearErr := strconv.Atoi(startDateTokens[0])
	startDateMonth, startDateMonthErr := strconv.Atoi(startDateTokens[1])
	startDateDay, startDateDayErr := strconv.Atoi(startDateTokens[2])
	if startDateYearErr != nil || startDateMonthErr != nil || startDateDayErr != nil {
		return "startDate must be numeric", startDate, deltaDays
	}
	if startDateMonth < 1 || startDateMonth > 12 {
		return "startDate month is invalid", startDate, deltaDays
	}
	if startDateDay < 1 || startDateDay > 31 {
		return "startDate day of month is invalid", startDate, deltaDays
	}

	startDate = time.Date(startDateYear, time.Month(startDateMonth), startDateDay, 0, 0, 0, 0, time.UTC)

	endDateYear, endDateYearErr := strconv.Atoi(endDateTokens[0])
	endDateMonth, endDateMonthErr := strconv.Atoi(endDateTokens[1])
	endDateDay, endDateDayErr := strconv.Atoi(endDateTokens[2])
	if endDateYearErr != nil || endDateMonthErr != nil || endDateDayErr != nil {
		return "endDate must be numeric", startDate, deltaDays
	}
	if endDateMonth < 1 || endDateMonth > 12 {
		return "endDate day of month is invalid", startDate, deltaDays
	}
	if endDateDay < 1 || endDateDay > 31 {
		return "endDate day of month is invalid", startDate, deltaDays
	}

	endDate := time.Date(endDateYear, time.Month(endDateMonth), endDateDay, 0, 0, 0, 0, time.UTC)
	if endDate.Before(startDate) {
		return "endDate can't be before startDate", startDate, deltaDays
	}

	delta := endDate.Sub(startDate)

	// difference in days between startDate and endDate
	deltaDays = int(delta.Hours() / 24)
	if deltaDays > configuration.Config.Virgilio.MaxDays {
		return "Maximum number of days exceeded", startDate, deltaDays
	}

	return "", startDate, deltaDays
}

// ReadWidget parses query and validate it
// Finally, return the widget result
func ReadWidget(c *gin.Context) {
	widgetName := c.Param("widgetName")
	startDateString := c.Query("startDate")
	endDateString := c.Query("endDate")

	message, startDate, deltaDays := validate(widgetName, startDateString, endDateString)
	if message != "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": message})
		return
	}

	filePaths := widgets.GetFileLists(widgetName, startDate, deltaDays)

	var widgetData map[string]interface{}

	var labelData widgets.Label

	var counterData widgets.Counter
	var valueOutputCounter float64

	var chartData widgets.Chart
	var numSeries int
	var numCategories int
	var seriesOutputChart []widgets.Series

	var tableData widgets.Table
	var numRows int
	var numColumns int
	var rowsOutputTable [][]float64
	breakLoop := false
	firstIteration := true

	// widget files are read starting from the most recent
	for index := len(filePaths) - 1; index >= 0; index-- {
		if breakLoop {
			break
		}
		filePath := filePaths[index]
		widgetFile, err := os.Open(filePath)
		defer widgetFile.Close()
		if err != nil {
			// missing files are skipped
			continue
		}
		bytes, err := ioutil.ReadAll(widgetFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error reading file", "error": err.Error()})
			return
		}

		json.Unmarshal(bytes, &widgetData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Unmarshaling error", "error": err.Error()})
			return
		}

		switch widgetType := widgetData["type"]; widgetType {
		case "label":
			mapstructure.Decode(widgetData, &labelData)
			// label widget are never aggregated, exiting loop
			breakLoop = true
		case "counter":
			mapstructure.Decode(widgetData, &counterData)
			valueOutputCounter += counterData.Value

			// if snapshot = true, don't aggregate
			if counterData.Snapshot {
				breakLoop = true
			}
		case "chart":
			mapstructure.Decode(widgetData, &chartData)
			series := chartData.Series

			if firstIteration {
				firstIteration = false
				numSeries = len(series)
				numCategories = len(chartData.Categories)
				seriesOutputChart = make([]widgets.Series, numSeries)

				for i := range seriesOutputChart {
					seriesOutputChart[i].Data = make([]float64, numCategories)
				}
			}

			for i := 0; i < numSeries; i++ {
				for j := 0; j < numCategories; j++ {
					seriesOutputChart[i].Data[j] += series[i].Data[j]
				}
			}

			// if snapshot = true, don't aggregate
			if chartData.Snapshot {
				breakLoop = true
			}
		case "table":
			mapstructure.Decode(widgetData, &tableData)
			rows := tableData.Rows

			if firstIteration {
				firstIteration = false
				numRows = len(rows)
				numColumns = len(rows[0])
				rowsOutputTable = make([][]float64, numRows)

				for i := range rowsOutputTable {
					rowsOutputTable[i] = make([]float64, numColumns)
				}
			}

			for i := 0; i < numRows; i++ {
				for j := 0; j < numColumns; j++ {
					rowsOutputTable[i][j] += rows[i][j]
				}
			}

			// if snapshot = true, don't aggregate
			if tableData.Snapshot {
				breakLoop = true
			}
		case nil:
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot retrieve widget type for " + filePath})
			return
		default:
			c.JSON(http.StatusNotImplemented, gin.H{"message": "Widget type not implemented: " + widgetType.(string)})
			return
		}
	}
	var widget interface{}

	switch widgetType := widgetData["type"]; widgetType {
	case "label":
		widget = labelData
	case "counter":
		counterData.Value = valueOutputCounter
		widget = counterData
	case "chart":
		chartData.Series = seriesOutputChart

		// if it's a pie chart, change structure of output json
		if chartData.ChartType == "pie" {
			pieChart := utils.MapChartToPieChart(chartData)
			widget = pieChart
		} else {
			widget = chartData
		}
	case "table":
		tableData.Rows = rowsOutputTable
		widget = tableData
	}

	c.JSON(http.StatusOK, gin.H{
		"widget": widget,
	})
}

// ListMiners list all Ciacco miners
func ListMiners(c *gin.Context) {
	miners, errString := GetMiners()
	if errString != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errString})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"miners": miners,
	})
}

func GetMiners() ([]miner, string) {
	var miners []miner

	files, err := ioutil.ReadDir(configuration.Config.Ciacco.MinersDirectory)
	if err != nil {
		return nil, "Miners directory not found " + configuration.Config.Ciacco.MinersDirectory
	}

	for _, f := range files {
		var m miner
		if !f.IsDir() {
			parts := strings.Split(f.Name(), "-")
			m.FullName = f.Name()
			m.Name = parts[0]
			m.Type = parts[1]
			miners = append(miners, m)
		}
	}
	return miners, ""
}

// DeleteLayout delete saved layout
func DeleteLayout(c *gin.Context) {
	var err = os.Remove(configuration.Config.Virgilio.LayoutFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	} else {
		c.Status(http.StatusOK)
	}
}

// GetLayout returns the saved layout from VIRGILIO_STORE_DIR,
// if no existing layout can be found, return the default one
func GetLayout(c *gin.Context) {
	layout := widgets.ReadLayout()

	c.JSON(http.StatusOK, gin.H{
		"layout":  layout.Widgets,
		"default": layout.Default,
	})
}

// SetLayoyt saves the layout inside VIRGILIO_STORE_DIR
func SetLayout(c *gin.Context) {
	var layout widgets.Layout

	if err := c.ShouldBindJSON(&layout); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
	file, _ := json.Marshal(layout)

	_ = ioutil.WriteFile(configuration.Config.Virgilio.LayoutFile, file, 0644)
}

// GetLang returns the i18n strings used by Beatrice, including those used by the widgets
func GetLang(c *gin.Context) {
	langCode := c.Param("langCode")
	basePath := configuration.Config.Beatrice.BaseDirectory + "/i18n/"
	i18nFile := basePath + langCode + ".json"
	i18nMap, err := utils.ReadJson(i18nFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	miners, errString := GetMiners()
	if errString != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": errString})
		return
	}

	for _, miner := range miners {
		minerName := miner.FullName
		minerI18nFile := basePath + minerName + "-" + langCode + ".json"
		minerI18nMap, err := utils.ReadJson(minerI18nFile)
		if err != nil {
			// non-fatal error
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		i18nMap[minerName] = minerI18nMap
	}

	c.JSON(http.StatusOK, gin.H{
		"lang": i18nMap,
	})
}
