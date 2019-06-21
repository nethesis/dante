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
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/nethesis/dante/virgilio/widgets"
)

func ContainsString(stringSlice []string, searchString string) bool {
	for _, value := range stringSlice {
		if value == searchString {
			return true
		}
	}
	return false
}

func MapChartToPieChart(chart widgets.Chart) widgets.PieChart {
	var pieChart widgets.PieChart
	pieChart.Type = chart.Type
	pieChart.ChartType = chart.ChartType
	pieChart.Title = chart.Title
	pieChart.MinerId = chart.MinerId
	pieChart.Snapshot = chart.Snapshot
	pieChart.Labels = chart.Categories
	// charts of type "pie" always have only one series
	pieChartSeries := make([]float64, len(chart.Series[0].Data))

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
