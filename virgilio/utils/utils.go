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

package utils

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/nethesis/dante/virgilio/widgets"
)

type HttpError struct {
	Code        int
	ErrorString string
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.ErrorString)
}

func ContainsString(stringSlice []string, searchString string) bool {
	for _, value := range stringSlice {
		if value == searchString {
			return true
		}
	}
	return false
}

func MapTableToTableUI(table widgets.Table) widgets.TableUI {
	var tableUi widgets.TableUI
	tableUi.Type = table.Type
	tableUi.MinerId = table.MinerId
	tableUi.Unit = table.Unit
	tableUi.AggregationType = table.AggregationType

	tableUi.Rows = make([][]string, len(table.Rows))

	if len(table.RowHeader) > 0 {
		tableUi.RowHeader = true
		if len(table.ColumnHeader) > 0 {
			// put an empty string at the beginning of the column header (top-left cell)
			tableUi.ColumnHeader = make([]string, len(table.ColumnHeader)+1)
			copy(tableUi.ColumnHeader, table.ColumnHeader)
			copy(tableUi.ColumnHeader[1:], tableUi.ColumnHeader[:])
			tableUi.ColumnHeader[0] = ""
		} else {
			tableUi.ColumnHeader = make([]string, 0)
		}

		// initialize tableUi.Rows
		for i := range tableUi.Rows {
			tableUi.Rows[i] = make([]string, len(table.Rows[0])+1)
		}

		// put every row header at the beginning of every row
		for i := 0; i < len(table.RowHeader); i++ {
			tableUi.Rows[i][0] = table.RowHeader[i]

			for j := 0; j < len(table.Rows[0]); j++ {
				tableUi.Rows[i][j+1] = fmt.Sprintf("%g", table.Rows[i][j])
			}
		}
	} else {
		// no row headers
		tableUi.ColumnHeader = table.ColumnHeader

		// initialize tableUi.Rows
		for i := range tableUi.Rows {
			tableUi.Rows[i] = make([]string, len(table.Rows[0]))
		}

		for i := 0; i < len(table.Rows); i++ {
			for j := 0; j < len(table.Rows[0]); j++ {
				tableUi.Rows[i][j] = fmt.Sprintf("%g", table.Rows[i][j])
			}
		}
	}
	return tableUi
}

func MapChartToPieChart(chart widgets.Chart) widgets.PieChart {
	var pieChart widgets.PieChart
	pieChart.Type = chart.Type
	pieChart.ChartType = chart.ChartType
	pieChart.MinerId = chart.MinerId
	pieChart.AggregationType = chart.AggregationType
	pieChart.Labels = chart.Categories
	// charts of type "pie" always have only one series
	pieChartSeries := make([]float64, len(chart.Series[0].Data))
	pieChart.Unit = chart.Unit

	for i, value := range chart.Series[0].Data {
		pieChartSeries[i] = value
	}
	pieChart.Series = pieChartSeries
	return pieChart
}

func ReadJson(filePath string) (map[string]interface{}, error) {
	var mapData map[string]interface{}

	jsonFile, err := os.Open(filePath)
	defer jsonFile.Close()
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(bytes, &mapData)
	if err != nil {
		return nil, err
	}
	return mapData, nil
}

func ReadJsonIgnoreOpenError(filePath string) (map[string]interface{}, bool, error) {
	var mapData map[string]interface{}

	jsonFile, err := os.Open(filePath)
	defer jsonFile.Close()
	if err != nil {
		return nil, true, err
	}
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, false, err
	}
	json.Unmarshal(bytes, &mapData)
	if err != nil {
		return nil, false, err
	}
	return mapData, false, nil
}

func ReverseSliceFloat(slice []float64) []float64 {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

func ReverseSliceString(slice []string) []string {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

func GetDateStringFromWidgetFilePath(filePath string) string {
	filePathTokens := strings.Split(filePath, "/")
	day := filePathTokens[len(filePathTokens)-2]
	month := filePathTokens[len(filePathTokens)-3]
	year := filePathTokens[len(filePathTokens)-4]
	return year + "-" + month + "-" + day
}

func Anonymize(value string) string {
	h := md5.New()
	h.Write([]byte(value))
	bs := h.Sum(nil)
	anonymizedString := fmt.Sprintf("%x", bs)
	return anonymizedString[:8]
}
