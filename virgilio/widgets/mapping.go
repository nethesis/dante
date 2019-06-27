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
	"fmt"
)

func MapTableToTableUI(table Table) TableUI {
	var tableUi TableUI
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

func MapChartToPieChart(chart Chart) PieChart {
	var pieChart PieChart
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
