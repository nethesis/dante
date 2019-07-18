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

func computeTrendValue(mostRecentValueTrend float64, leastRecentValueTrend float64, counterData widgets.Counter) (float64, error) {
	var trend float64

	if counterData.TrendType == "number" {
		trend = mostRecentValueTrend - leastRecentValueTrend
		return trend, nil
	} else if counterData.TrendType == "percentage" {
		if leastRecentValueTrend != 0 {
			trend = (mostRecentValueTrend - leastRecentValueTrend) / leastRecentValueTrend * 100
		} else {
			// avoid division by zero
			trend = 0
		}
		trendRounded := math.Round(trend*100) / 100
		return trendRounded, nil
	} else {
		return 0, &utils.HttpError{http.StatusInternalServerError, "Trend type not implemented: " + counterData.TrendType}
	}
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

func readMostRecentWidgetFile(filePaths []string) (widgets.WidgetBlob, error) {
	var filePathIndex int
	var widgetData map[string]interface{}
	firstFileRead := false
	var openError bool
	var err error
	aggregate := false
	var valueOutputCounter float64
	var lastValidWidgetData map[string]interface{}
	var labelData widgets.Label
	var counterData widgets.Counter
	var chartData widgets.Chart
	var tableData widgets.Table
	var listData widgets.List
	var trendSeries []float64
	var widgetType string
	var trendCategories []string
	var mostRecentValueTrend float64
	var seriesOutputChart []widgets.Series
	var rowsOutputTable [][]float64
	var listMap map[string]float64
	var widgetBlob widgets.WidgetBlob
	var widgetInfo widgets.WidgetInfo
	var labelInfo widgets.LabelInfo
	var counterInfo widgets.CounterInfo
	var chartInfo widgets.ChartInfo
	var tableInfo widgets.TableInfo
	var listInfo widgets.ListInfo

	for filePathIndex = len(filePaths) - 1; filePathIndex >= 0 && !firstFileRead; filePathIndex-- {
		filePath := filePaths[filePathIndex]
		widgetData, openError, err = utils.ReadJsonIgnoreOpenError(filePath)
		if err != nil {
			if openError {
				widgetData = lastValidWidgetData
				// skip to next most recent widget file
				continue
			} else {
				return widgetBlob, &utils.HttpError{http.StatusInternalServerError, err.Error()}
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
			labelInfo.LabelData = labelData

			// assign labelInfo to widgetBlob
			widgetBlob.LabelInfo = labelInfo
		case "counter":
			mapstructure.Decode(widgetData, &counterData)
			valueOutputCounter = counterData.Value
			trendSeries = append(trendSeries, counterData.Value)
			dateString := utils.GetDateStringFromFilePath(filePath)
			trendCategories = append(trendCategories, dateString)

			if counterData.AggregationType == "snapshot" {
				// save most recent value to compute trend
				mostRecentValueTrend = counterData.Value
			}
			// even if it's snapshot, aggregation is needed to compute trend
			aggregate = true

			// assign props to counterInfo object
			counterInfo.CounterData = counterData
			counterInfo.TrendSeries = trendSeries
			counterInfo.ValueOutputCounter = valueOutputCounter
			counterInfo.TrendCategories = trendCategories
			counterInfo.MostRecentValueTrend = mostRecentValueTrend

			// assign counterInfo to widgetBlob
			widgetBlob.CounterInfo = counterInfo
		case "chart":
			chartData, seriesOutputChart = AggregateChart(widgetData, chartData, seriesOutputChart, true)

			// assign props to chartInfo object
			chartInfo.ChartData = chartData
			chartInfo.SeriesOutputChart = seriesOutputChart

			// assign chartInfo to widgetBlob
			widgetBlob.ChartInfo = chartInfo
		case "table":
			tableData, rowsOutputTable = AggregateTable(widgetData, tableData, rowsOutputTable, true)

			// assign props to tableInfo object
			tableInfo.TableData = tableData
			tableInfo.RowsOutputTable = rowsOutputTable

			// assign tableInfo to widgetBlob
			widgetBlob.TableInfo = tableInfo
		case "list":
			listData, listMap = AggregateList(widgetData, listData, listMap, true)

			// assign props to listInfo object
			listInfo.ListMap = listMap
			listInfo.ListData = listData

			// assign listInfo to widgetBlob
			widgetBlob.ListInfo = listInfo
		default:
			return widgetBlob, &utils.HttpError{http.StatusNotImplemented, "Widget type not implemented: " + widgetType}
		}
	}

	widgetInfo.WidgetData = widgetData
	widgetInfo.LastValidWidgetData = lastValidWidgetData
	widgetInfo.Aggregate = aggregate
	widgetInfo.WidgetType = widgetType
	widgetInfo.FilePaths = filePaths
	widgetInfo.FilePathIndex = filePathIndex

	widgetBlob.WidgetInfo = widgetInfo

	return widgetBlob, nil
}

func readCounterWidget(widgetInfo widgets.WidgetInfo, counterInfo widgets.CounterInfo) (widgets.Counter, error) {
	var openError bool
	var err error
	var leastRecentValueTrend float64
	var trend float64
	var lastValidWidgetDataForTrend map[string]interface{}

	counterData := counterInfo.CounterData
	aggregate := widgetInfo.Aggregate
	filePaths := widgetInfo.FilePaths
	filePathIndex := widgetInfo.FilePathIndex
	widgetData := widgetInfo.WidgetData
	lastValidWidgetData := widgetInfo.LastValidWidgetData
	widgetName := widgetInfo.WidgetName
	startDate := widgetInfo.StartDate
	deltaDays := widgetInfo.DeltaDays
	valueOutputCounter := counterInfo.ValueOutputCounter
	trendSeries := counterInfo.TrendSeries
	trendCategories := counterInfo.TrendCategories
	mostRecentValueTrend := counterInfo.MostRecentValueTrend

	if aggregate {
		for ; filePathIndex >= 0; filePathIndex-- {
			filePath := filePaths[filePathIndex]
			widgetData, openError, err = utils.ReadJsonIgnoreOpenError(filePath)
			if err != nil {
				if openError {
					widgetData = lastValidWidgetData
					// skip to next widget file
					continue
				} else {
					return counterData, &utils.HttpError{http.StatusInternalServerError, err.Error()}
				}
			}
			lastValidWidgetData = widgetData

			widgetType := widgetData["type"].(string)
			if widgetType != "counter" {
				return counterData, &utils.HttpError{http.StatusInternalServerError, "Unexpected error, widgetType should be counter, but is " + widgetType}
			}

			mapstructure.Decode(widgetData, &counterData)

			if counterData.AggregationType == "sum" {
				valueOutputCounter += counterData.Value
			}
			// trend management
			leastRecentValueTrend = counterData.Value
			trendSeries = append(trendSeries, counterData.Value)
			dateString := utils.GetDateStringFromFilePath(filePath)
			trendCategories = append(trendCategories, dateString)
		}
	}
	counterData.Value = valueOutputCounter

	if counterData.AggregationType == "snapshot" {
		trend, err = computeTrendValue(mostRecentValueTrend, leastRecentValueTrend, counterData)
		if err != nil {
			return counterData, err
		}
	} else if counterData.AggregationType == "sum" {
		// second aggregation
		finalValueOutputCounter := valueOutputCounter
		mostRecentValueTrend = valueOutputCounter
		valueOutputCounter = 0
		startDate = startDate.AddDate(0, 0, -deltaDays-1)
		filePaths = widgets.GetFileLists(widgetName, startDate, deltaDays)
		leastRecentValueTrendFirstAggregation := leastRecentValueTrend

		for filePathIndex = len(filePaths) - 1; filePathIndex >= 0; filePathIndex-- {
			filePath := filePaths[filePathIndex]
			widgetData, openError, err = utils.ReadJsonIgnoreOpenError(filePath)
			if err != nil {
				if openError {
					widgetData = lastValidWidgetDataForTrend
					// skip to next widget file
					continue
				} else {
					return counterData, &utils.HttpError{http.StatusInternalServerError, err.Error()}
				}
			}
			lastValidWidgetDataForTrend = widgetData

			mapstructure.Decode(widgetData, &counterData)
			valueOutputCounter += counterData.Value
		}
		leastRecentValueTrend = valueOutputCounter

		// if there is no data available in second aggregation, use least recent value of first aggregation
		if leastRecentValueTrend == 0 {
			leastRecentValueTrend = leastRecentValueTrendFirstAggregation
		}
		trend, err = computeTrendValue(mostRecentValueTrend, leastRecentValueTrend, counterData)
		if err != nil {
			return counterData, err
		}
		counterData.Value = finalValueOutputCounter
	}
	counterData.Trend = trend
	trendSeries = utils.ReverseSliceFloat(trendSeries)
	trendSeriesJson := widgets.TrendSeries{"trendSeries", trendSeries}
	counterData.TrendSeries = []widgets.TrendSeries{trendSeriesJson}
	trendCategories = utils.ReverseSliceString(trendCategories)
	counterData.TrendCategories = trendCategories
	return counterData, nil
}

func readChartWidget(widgetInfo widgets.WidgetInfo, chartInfo widgets.ChartInfo) (interface{}, error) {
	var openError bool
	var err error

	aggregate := widgetInfo.Aggregate
	filePaths := widgetInfo.FilePaths
	filePathIndex := widgetInfo.FilePathIndex
	widgetData := widgetInfo.WidgetData
	lastValidWidgetData := widgetInfo.LastValidWidgetData
	chartData := chartInfo.ChartData
	seriesOutputChart := chartInfo.SeriesOutputChart

	if aggregate {
		for ; filePathIndex >= 0; filePathIndex-- {
			filePath := filePaths[filePathIndex]
			widgetData, openError, err = utils.ReadJsonIgnoreOpenError(filePath)
			if err != nil {
				if openError {
					widgetData = lastValidWidgetData
					// skip to next widget file
					continue
				} else {
					return chartData, &utils.HttpError{http.StatusInternalServerError, err.Error()}
				}
			}
			lastValidWidgetData = widgetData

			widgetType := widgetData["type"].(string)
			if widgetType != "chart" {
				return chartData, &utils.HttpError{http.StatusInternalServerError, "Unexpected error, widgetType should be chart, but is " + widgetType}
			}
			chartData, seriesOutputChart = AggregateChart(widgetData, chartData, seriesOutputChart, false)
		}
	}
	chartData.Series = seriesOutputChart

	if configuration.Config.Virgilio.Anonymize && chartData.Anonymizable {
		for i := 0; i < len(chartData.Categories); i++ {
			chartData.Categories[i] = utils.Anonymize(chartData.Categories[i])
		}
	}

	if chartData.ChartType == "pie" {
		pieChart := widgets.MapChartToPieChart(chartData)
		return pieChart, nil
	}
	return chartData, nil
}

func readTableWidget(widgetInfo widgets.WidgetInfo, tableInfo widgets.TableInfo) (widgets.TableUI, error) {
	var openError bool
	var err error
	var tableUi widgets.TableUI

	aggregate := widgetInfo.Aggregate
	filePaths := widgetInfo.FilePaths
	filePathIndex := widgetInfo.FilePathIndex
	widgetData := widgetInfo.WidgetData
	lastValidWidgetData := widgetInfo.LastValidWidgetData
	tableData := tableInfo.TableData
	rowsOutputTable := tableInfo.RowsOutputTable

	if aggregate {
		for ; filePathIndex >= 0; filePathIndex-- {
			filePath := filePaths[filePathIndex]
			widgetData, openError, err = utils.ReadJsonIgnoreOpenError(filePath)
			if err != nil {
				if openError {
					widgetData = lastValidWidgetData
					// skip to next widget file
					continue
				} else {
					return tableUi, &utils.HttpError{http.StatusInternalServerError, err.Error()}
				}
			}
			lastValidWidgetData = widgetData

			widgetType := widgetData["type"].(string)
			if widgetType != "table" {
				return tableUi, &utils.HttpError{http.StatusInternalServerError, "Unexpected error, widgetType should be table, but is " + widgetType}
			}
			tableData, rowsOutputTable = AggregateTable(widgetData, tableData, rowsOutputTable, false)
		}
	}

	if configuration.Config.Virgilio.Anonymize && tableData.Anonymizable {
		for i := 0; i < len(tableData.RowHeader); i++ {
			tableData.RowHeader[i] = utils.Anonymize(tableData.RowHeader[i])
		}
	}
	tableUi = widgets.MapTableToTableUI(tableData)
	return tableUi, nil
}

func readListWidget(widgetInfo widgets.WidgetInfo, listInfo widgets.ListInfo) (widgets.List, error) {
	var openError bool
	var err error

	aggregate := widgetInfo.Aggregate
	filePaths := widgetInfo.FilePaths
	filePathIndex := widgetInfo.FilePathIndex
	widgetData := widgetInfo.WidgetData
	lastValidWidgetData := widgetInfo.LastValidWidgetData
	listData := listInfo.ListData
	listMap := listInfo.ListMap

	if aggregate {
		for ; filePathIndex >= 0; filePathIndex-- {
			filePath := filePaths[filePathIndex]
			widgetData, openError, err = utils.ReadJsonIgnoreOpenError(filePath)
			if err != nil {
				if openError {
					widgetData = lastValidWidgetData
					// skip to next widget file
					continue
				} else {
					return listData, &utils.HttpError{http.StatusInternalServerError, err.Error()}
				}
			}
			lastValidWidgetData = widgetData

			widgetType := widgetData["type"].(string)
			if widgetType != "list" {
				return listData, &utils.HttpError{http.StatusInternalServerError, "Unexpected error, widgetType should be list, but is " + widgetType}
			}
			listData, listMap = AggregateList(widgetData, listData, listMap, false)
		}
	}
	listData.Data = make([]widgets.ListElem, 0)

	for key, value := range listMap {
		listElem := widgets.ListElem{key, value}
		listData.Data = append(listData.Data, listElem)
	}
	sort.Slice(listData.Data, func(i, j int) bool {
		return listData.Data[i].Count > listData.Data[j].Count
	})
	// limit number of list entries
	if len(listData.Data) > configuration.Config.Virgilio.MaxEntries {
		listData.Data = listData.Data[0:configuration.Config.Virgilio.MaxEntries]
	}
	if configuration.Config.Virgilio.Anonymize && listData.Anonymizable {
		for key, el := range listData.Data {
			el.Name = utils.Anonymize(el.Name)
			listData.Data[key] = el
		}
	}
	return listData, nil
}

// ReadWidget parses query and validate it
// Finally, return the widget result
func ReadWidget(c *gin.Context) {
	var widget interface{}

	widgetName := c.Param("widgetName")
	startDateString := c.Query("startDate")
	endDateString := c.Query("endDate")

	message, startDate, deltaDays := validate(widgetName, startDateString, endDateString)
	if message != "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": message})
		return
	}
	filePaths := widgets.GetFileLists(widgetName, startDate, deltaDays)

	widgetBlob, err := readMostRecentWidgetFile(filePaths)
	if err != nil {
		if httpError, ok := err.(*utils.HttpError); ok {
			c.JSON(httpError.Code, gin.H{"message": httpError.ErrorString})
			return
		}
	}

	widgetInfo := widgetBlob.WidgetInfo
	widgetInfo.WidgetName = widgetName
	widgetInfo.StartDate = startDate
	widgetInfo.DeltaDays = deltaDays

	switch widgetInfo.WidgetType {
	case "label":
		widget = widgetBlob.LabelInfo.LabelData
	case "counter":
		counterData, err := readCounterWidget(widgetInfo, widgetBlob.CounterInfo)
		if err != nil {
			if httpError, ok := err.(*utils.HttpError); ok {
				c.JSON(httpError.Code, gin.H{"message": httpError.ErrorString})
				return
			}
		}
		widget = counterData
	case "chart":
		widget, err = readChartWidget(widgetInfo, widgetBlob.ChartInfo)
		if err != nil {
			if httpError, ok := err.(*utils.HttpError); ok {
				c.JSON(httpError.Code, gin.H{"message": httpError.ErrorString})
				return
			}
		}
	case "table":
		tableData, err := readTableWidget(widgetInfo, widgetBlob.TableInfo)
		if err != nil {
			if httpError, ok := err.(*utils.HttpError); ok {
				c.JSON(httpError.Code, gin.H{"message": httpError.ErrorString})
				return
			}
		}
		widget = tableData
	case "list":
		listData, err := readListWidget(widgetInfo, widgetBlob.ListInfo)
		if err != nil {
			if httpError, ok := err.(*utils.HttpError); ok {
				c.JSON(httpError.Code, gin.H{"message": httpError.ErrorString})
				return
			}
		}
		widget = listData
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

		// exclude directories and file names starting with "." (e.g. Vim swap files)
		if !f.IsDir() && string(f.Name()[0]) != "." {
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

// GetDefaultLayout returns the default layout
func GetDefaultLayout(c *gin.Context) {
	layout := widgets.ReadDefaultLayout()

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
	basePath := configuration.Config.Beatrice.BaseDirectory + "/i18n"
	i18nFile := basePath + "/" + langCode + "/" + langCode + ".json"
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
		minerI18nFile := basePath + "/" + langCode + "/" + minerName + ".json"
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
