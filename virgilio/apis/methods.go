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
	"math"
	"net/http"
	"os"
	"sort"
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
	Name string `json:"name"`
	Type string `json:"type"`
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

// AggregateCounter decodes a Counter from a map and aggregates its value // todo del
// func AggregateCounter(widgetData map[string]interface{}, counterData widgets.Counter, valueOutputCounter float64) (widgets.Counter, float64) {
// 	mapstructure.Decode(widgetData, &counterData)
// 	valueOutputCounter += counterData.Value
// 	return counterData, valueOutputCounter
// }

func initChart(widgetData map[string]interface{}, chartData widgets.Chart, seriesOutputChart []widgets.Series) (widgets.Chart, []widgets.Series) {
	mapstructure.Decode(widgetData, &chartData)
	series := chartData.Series
	numSeries := len(series)
	numCategories := len(chartData.Categories)
	seriesOutputChart = make([]widgets.Series, numSeries)

	for i := range seriesOutputChart {
		seriesOutputChart[i].Data = make([]float64, numCategories)
	}

	for i := 0; i < numSeries; i++ {
		for j := 0; j < numCategories; j++ {
			seriesOutputChart[i].Data[j] += series[i].Data[j]
		}
	}
	return chartData, seriesOutputChart
}

// AggregateChart decodes a Chart from a map and aggregates its series
func AggregateChart(widgetData map[string]interface{}, chartData widgets.Chart, seriesOutputChart []widgets.Series, initialize bool) (widgets.Chart, []widgets.Series) {
	mapstructure.Decode(widgetData, &chartData)
	series := chartData.Series
	numSeries := len(series)
	numCategories := len(chartData.Categories)

	if initialize {
		seriesOutputChart = make([]widgets.Series, numSeries)

		for i := range seriesOutputChart {
			seriesOutputChart[i].Data = make([]float64, numCategories)
		}
	}

	for i := 0; i < numSeries; i++ {
		for j := 0; j < numCategories; j++ {
			seriesOutputChart[i].Data[j] += series[i].Data[j]
		}
		seriesOutputChart[i].Name = series[i].Name
	}
	return chartData, seriesOutputChart
}

func computeSnapshotTrendValue(mostRecentValueTrend float64, leastRecentValueTrend float64, counterData widgets.Counter) (float64, string) {
	// retrieve least recent value to compute trend
	// leastRecentValueTrend, errorString := GetLeastRecentValueSnapshotTrend(filePaths)
	// if errorString != "" {
	// 	return 0, errorString
	// }

	if counterData.TrendType == "number" {
		trend := mostRecentValueTrend - leastRecentValueTrend
		return trend, ""
	} else if counterData.TrendType == "percentage" {
		trend := (mostRecentValueTrend - leastRecentValueTrend) / leastRecentValueTrend * 100
		fmt.Println("mostRecentValueTrend", mostRecentValueTrend, "leastRecentValueTrend", leastRecentValueTrend, "trend%", trend) // todo del
		trendRounded := math.Round(trend*100) / 100
		return trendRounded, ""
	} else {
		return 0, "TrendType not implemented: " + counterData.TrendType
	}
}

// func GetLeastRecentValueSnapshotTrend(filePaths []string) (float64, string) {
// 	firstFileRead := false
// 	var counterData widgets.Counter

// 	// read least recent widget file
// 	for index := 0; index < len(filePaths) && !firstFileRead; index++ {
// 		filePath := filePaths[index]
// 		widgetData, openError, err := utils.ReadJsonIgnoreOpenError(filePath)
// 		if err != nil {
// 			if openError {
// 				// skip to next most recent widget file
// 				continue
// 			} else {
// 				return 0, err.Error()
// 			}
// 		}
// 		firstFileRead = true

// 		switch widgetType := widgetData["type"]; widgetType {
// 		case "counter":
// 			mapstructure.Decode(widgetData, &counterData)
// 		case nil:
// 			errorString := "Cannot retrieve widget type for " + filePath
// 			return 0, errorString
// 		default:
// 			errorString := "Widget type not implemented: " + widgetType.(string)
// 			return 0, errorString
// 		}
// 	}
// 	return counterData.Value, ""
// }

func initTable(widgetData map[string]interface{}, tableData widgets.Table, rowsOutputTable [][]float64) (widgets.Table, [][]float64) {
	mapstructure.Decode(widgetData, &tableData)
	rows := tableData.Rows
	numRows := len(rows)
	numColumns := len(rows[0])
	rowsOutputTable = make([][]float64, numRows)

	for i := range rowsOutputTable {
		rowsOutputTable[i] = make([]float64, numColumns)
	}

	for i := 0; i < numRows; i++ {
		for j := 0; j < numColumns; j++ {
			rowsOutputTable[i][j] += rows[i][j]
		}
	}
	return tableData, rowsOutputTable
}

// AggregateTable decodes a Table from a map and aggregates its data cells
func AggregateTable(widgetData map[string]interface{}, tableData widgets.Table, rowsOutputTable [][]float64, initialize bool) (widgets.Table, [][]float64) {
	mapstructure.Decode(widgetData, &tableData)
	rows := tableData.Rows
	numRows := len(rows)
	numColumns := len(rows[0])

	if initialize {
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
	return tableData, rowsOutputTable
}

// AggregateList decodes a List from a map and aggregates its values in a map
func AggregateList(widgetData map[string]interface{}, listData widgets.List, listMap map[string]float64, initialize bool) (widgets.List, map[string]float64) {
	mapstructure.Decode(widgetData, &listData)
	numElems := len(listData.Data)

	if initialize {
		listMap = make(map[string]float64)
	}

	for i := 0; i < numElems; i++ {
		currentCount, exists := listMap[listData.Data[i].Name]
		if !exists {
			listMap[listData.Data[i].Name] = listData.Data[i].Count
		} else {
			// sum
			listMap[listData.Data[i].Name] = currentCount + listData.Data[i].Count
		}
	}
	return listData, listMap
}

// ReadWidget parses query and validate it
// Finally, return the widget result
func ReadWidget(c *gin.Context) {
	var widgetData map[string]interface{}
	var lastValidWidgetData map[string]interface{}
	var openError bool
	var err error
	var labelData widgets.Label
	var counterData widgets.Counter
	var valueOutputCounter float64
	var chartData widgets.Chart
	var seriesOutputChart []widgets.Series
	var tableData widgets.Table
	var rowsOutputTable [][]float64
	firstFileRead := false
	var index int
	aggregate := false
	var mostRecentValueTrend float64
	var leastRecentValueTrend float64
	var errorString string
	var trend float64
	var widgetType string
	var trendSeries []float64
	var trendCategories []string
	var listMap map[string]float64
	var listData widgets.List

	widgetName := c.Param("widgetName")
	startDateString := c.Query("startDate")
	endDateString := c.Query("endDate")

	message, startDate, deltaDays := validate(widgetName, startDateString, endDateString)
	if message != "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": message})
		return
	}
	filePaths := widgets.GetFileLists(widgetName, startDate, deltaDays)

	// read most recent widget file
	for index = len(filePaths) - 1; index >= 0 && !firstFileRead; index-- {
		filePath := filePaths[index]
		fmt.Println("reading most recent", filePath) // todo del
		widgetData, openError, err = utils.ReadJsonIgnoreOpenError(filePath)
		if err != nil {
			if openError {
				widgetData = lastValidWidgetData
				// skip to next most recent widget file
				continue
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
		}
		lastValidWidgetData = widgetData
		firstFileRead = true

		if widgetData["aggregationType"] == "sum" {
			aggregate = true
		}

		widgetType = widgetData["type"].(string)
		switch widgetType {
		case "label":
			mapstructure.Decode(widgetData, &labelData)
		case "counter":
			mapstructure.Decode(widgetData, &counterData)
			valueOutputCounter = counterData.Value
			trendSeries = append(trendSeries, counterData.Value)
			dateString := utils.GetDateStringFromWidgetFilePath(filePath)
			trendCategories = append(trendCategories, dateString)

			if counterData.AggregationType == "snapshot" {
				// save most recent value to compute trend
				mostRecentValueTrend = counterData.Value
			}
			// even if it's snapshot, aggregation is needed to display trend
			aggregate = true
		case "chart":
			chartData, seriesOutputChart = AggregateChart(widgetData, chartData, seriesOutputChart, true)
		case "table":
			tableData, rowsOutputTable = AggregateTable(widgetData, tableData, rowsOutputTable, true)
		case "list":
			listData, listMap = AggregateList(widgetData, listData, listMap, true)
		default:
			c.JSON(http.StatusNotImplemented, gin.H{"message": "Widget type not implemented: " + widgetType})
			return
		}
	}

	// todo refactoring...
	// if aggregate {
	// 	switch widgetType {
	// 	case "counter":
	// 		aggregateCounter(mostRecentValueTrend)
	// 	case "chart":
	// 		aggregateChart()
	// 	case "table":
	// 		aggregateTable()
	// 	case "list":
	// 		// todo
	// 	}
	// }

	if aggregate {
		// aggregate widget data
		for ; index >= 0; index-- {
			filePath := filePaths[index]
			fmt.Println("aggregating", filePath) // todo del
			widgetData, openError, err = utils.ReadJsonIgnoreOpenError(filePath)
			if err != nil {
				if openError {
					widgetData = lastValidWidgetData
					// skip to next most recent widget file
					continue
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
					return
				}
			}
			lastValidWidgetData = widgetData

			widgetType = widgetData["type"].(string)
			switch widgetType {
			// case "label":
			// can't aggregate label widget
			case "counter":
				mapstructure.Decode(widgetData, &counterData)

				if counterData.AggregationType == "sum" {
					valueOutputCounter += counterData.Value
				}
				// trend management
				leastRecentValueTrend = counterData.Value
				trendSeries = append(trendSeries, counterData.Value)
				dateString := utils.GetDateStringFromWidgetFilePath(filePath)
				trendCategories = append(trendCategories, dateString)
			case "chart":
				chartData, seriesOutputChart = AggregateChart(widgetData, chartData, seriesOutputChart, false)
			case "table":
				tableData, rowsOutputTable = AggregateTable(widgetData, tableData, rowsOutputTable, false)
			case "list":
				listData, listMap = AggregateList(widgetData, listData, listMap, false)
			default:
				c.JSON(http.StatusNotImplemented, gin.H{"message": "Widget type not implemented: " + widgetType})
				return
			}
		}
	}
	var widget interface{}

	widgetType = widgetData["type"].(string)
	switch widgetType {
	case "label":
		widget = labelData
	case "counter":
		fmt.Println("valueOutputCounter", valueOutputCounter) // todo del
		counterData.Value = valueOutputCounter

		if counterData.AggregationType == "snapshot" {
			trend, errorString = computeSnapshotTrendValue(mostRecentValueTrend, leastRecentValueTrend, counterData)
			if errorString != "" {
				c.JSON(http.StatusInternalServerError, gin.H{"message": errorString})
				return
			}

		} else if counterData.AggregationType == "sum" {
			// second aggregation
			finalValueOutputCounter := valueOutputCounter
			mostRecentValueTrend = valueOutputCounter
			valueOutputCounter = 0
			startDate = startDate.AddDate(0, 0, -deltaDays-1)
			filePaths = widgets.GetFileLists(widgetName, startDate, deltaDays)

			for index = len(filePaths) - 1; index >= 0; index-- {
				filePath := filePaths[index]
				widgetData, openError, err = utils.ReadJsonIgnoreOpenError(filePath)
				if err != nil {
					if openError {
						widgetData = lastValidWidgetData
						// skip to next most recent widget file
						continue
					} else {
						c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
						return
					}
				}

				widgetType = widgetData["type"].(string)
				switch widgetType {
				case "counter":
					mapstructure.Decode(widgetData, &counterData)
					valueOutputCounter += counterData.Value
				default:
					c.JSON(http.StatusNotImplemented, gin.H{"message": "Widget type not expected, expecting counter: " + widgetType})
					return
				}
			}
			leastRecentValueTrend = valueOutputCounter
			trend, errorString = computeSnapshotTrendValue(mostRecentValueTrend, leastRecentValueTrend, counterData)
			fmt.Println("mostRecentValueTrend, leastRecentValueTrend, trend", mostRecentValueTrend, leastRecentValueTrend, trend) // todo del
			counterData.Value = finalValueOutputCounter
		}
		counterData.Trend = trend
		trendSeries = utils.ReverseSliceFloat(trendSeries)
		trendSeriesJson := widgets.TrendSeries{"trendSeries", trendSeries}
		counterData.TrendSeries = []widgets.TrendSeries{trendSeriesJson}
		trendCategories = utils.ReverseSliceString(trendCategories)
		counterData.TrendCategories = trendCategories
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
		tableUi := utils.MapTableToTableUI(tableData)
		// tableData.Rows = rowsOutputTable // todo del
		widget = tableUi
	case "list":
		listData.Data = make([]widgets.ListElem, 0)

		for key, value := range listMap {
			listElem := widgets.ListElem{key, value}
			listData.Data = append(listData.Data, listElem)
		}
		sort.Slice(listData.Data, func(i, j int) bool {
			return listData.Data[i].Count > listData.Data[j].Count
		})
		widget = listData
	}

	c.JSON(http.StatusOK, gin.H{
		"widget": widget,
	})
}

// func aggregateCounter(filePaths []string, mostRecentValueTrend float64, index int, trendSeries []float64) utils.HttpError {
// 	for ; index >= 0; index-- {
// 		var lastValidWidgetData map[string]interface{}

// 		filePath := filePaths[index]
// 		widgetData, openError, err := utils.ReadJsonIgnoreOpenError(filePath)
// 		if err != nil {
// 			if openError {
// 				widgetData = lastValidWidgetData
// 				// skip to next most recent widget file
// 				continue
// 			} else {
// 				return &utils.HttpError{http.StatusInternalServerError, err.Error()}
// 			}
// 		}
// 		lastValidWidgetData = widgetData
// 		mapstructure.Decode(widgetData, &counterData)

// 		if counterData.AggregationType == "sum" {
// 			valueOutputCounter += counterData.Value
// 		} else if counterData.aggregationType == "snapshot" {
// 			// trend series
// 			trendSeries = append(trendSeries, counterData.Value)

// 			trend, errorString = computeSnapshotTrendValue(mostRecentValueTrend, counterData, filePaths)
// 			if errorString != "" {
// 				c.JSON(http.StatusInternalServerError, gin.H{"message": errorString})
// 				return
// 			}
// 			fmt.Println("snapshot trend", trend) // todo del
// 			counterData.Trend = trend
// 		}
// 	}

// }

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
			m.Name = f.Name()
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
		minerName := miner.Name
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
