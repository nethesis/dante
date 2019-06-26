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

type Layout struct {
	Widgets []Widget `json:"layout"`
	Default bool     `json:"default"`
}

type Widget struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	I      int     `json:"i"`
	Id     string  `json:"id"`
	Type   string  `json:"type"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	W      float64 `json:"w"`
	H      float64 `json:"h"`
	Text   string  `json:"text"`
}

type Label struct {
	Type    string `json:"type"`
	MinerId string `json:"minerId"`
	Value   string `json:"value"`
}

type Counter struct {
	Type            string        `json:"type"`
	MinerId         string        `json:"minerId"`
	Value           float64       `json:"value"`
	AggregationType string        `json:"aggregationType"`
	Unit            string        `json:"unit"`
	TrendType       string        `json:"trendType"`
	Trend           float64       `json:"trend"`
	TrendSeries     []TrendSeries `json:"trendSeries"`
	TrendCategories []string      `json:"trendCategories"`
}

type TrendSeries struct {
	Name string    `json:"name"`
	Data []float64 `json:"data"`
}

type Chart struct {
	Type            string   `json:"type"`
	ChartType       string   `json:"chartType"`
	MinerId         string   `json:"minerId"`
	Categories      []string `json:"categories"`
	Series          []Series `json:"series"`
	AggregationType string   `json:"aggregationType"`
	Unit            string   `json:"unit"`
}

type PieChart struct {
	Type            string    `json:"type"`
	ChartType       string    `json:"chartType"`
	MinerId         string    `json:"minerId"`
	Labels          []string  `json:"labels"`
	Series          []float64 `json:"series"`
	AggregationType string    `json:"aggregationType"`
	Unit            string    `json:"unit"`
}

type Table struct {
	Type            string      `json:"type"`
	MinerId         string      `json:"minerId"`
	Unit            string      `json:"unit"`
	ColumnHeader    []string    `json:"columnHeader"`
	RowHeader       []string    `json:"rowHeader"`
	Rows            [][]float64 `json:"rows"`
	AggregationType string      `json:"aggregationType"`
}

type TableUI struct {
	Type            string     `json:"type"`
	MinerId         string     `json:"minerId"`
	Unit            string     `json:"unit"`
	ColumnHeader    []string   `json:"columnHeader"`
	RowHeader       bool       `json:"rowHeader"`
	Rows            [][]string `json:"rows"`
	AggregationType string     `json:"aggregationType"`
}

type Series struct {
	Name string    `json:"name"`
	Data []float64 `json:"data"`
}

type List struct {
	Type            string     `json:"type"`
	MinerId         string     `json:"minerId"`
	Data            []ListElem `json:"data"`
	AggregationType string     `json:"aggregationType"`
	Unit            string     `json:"unit"`
}

type ListElem struct {
	Name  string  `json:"name"`
	Count float64 `json:"count"`
}
