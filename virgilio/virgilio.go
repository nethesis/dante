package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"encoding/json"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/mitchellh/mapstructure"
)

type Label struct {
	Type     string   `json:"type"`
	Title    string   `json:"title"`
	MinerId  string   `json:"minerId"`
	Value    string   `json:"value"`
	Position int      `json:"position"`
	Tags     []string `json:"tags"`
	Snapshot bool     `json:"snapshot"`
}

type Counter struct {
	Type     string   `json:"type"`
	Title    string   `json:"title"`
	MinerId  string   `json:"minerId"`
	Value    float64  `json:"value"`
	Position int      `json:"position"`
	Tags     []string `json:"tags"`
	Snapshot bool     `json:"snapshot"`
}

type Chart struct {
	Type       string   `json:"type"`
	ChartType  string   `json:"chartType"`
	Title      string   `json:"title"`
	MinerId    string   `json:"minerId"`
	Position   int      `json:"position"`
	Tags       []string `json:"tags"`
	Categories []string `json:"categories"`
	Series     []Series `json:"series"`
	Snapshot   bool     `json:"snapshot"`
}

type Table struct {
	Type         string      `json:"type"`
	Title        string      `json:"title"`
	MinerId      string      `json:"minerId"`
	Position     int         `json:"position"`
	Tags         []string    `json:"tags"`
	Unit         string      `json:"bytes"`
	ColumnHeader []string    `json:"columnHeader"`
	RowHeader    []string    `json:"rowHeader"`
	Rows         [][]float64 `json:"rows"`
	Snapshot     bool        `json:"snapshot"`
}

type Series struct {
	Name string    `json:"name"`
	I18n string    `json:"i18n"`
	Data []float64 `json:"data"`
}

func main() {
	widgetDataRoot := "/home/aleardini/Downloads/var/lib/nethserver/dante/" // todo edit
	numMaxDays := 366 // todo edit
	router := gin.Default()

	// cors
	corsConf := cors.DefaultConfig()
	corsConf.AllowOrigins = []string{"http://localhost:8080"}
	router.Use(cors.New(corsConf))

	router.GET("/widget/:widgetName", func(c *gin.Context) {
		widgetName := c.Param("widgetName")
		if widgetName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "widgetName is mandatory"})
			return
		}

		startDateString := c.Query("startDate")
		endDateString := c.Query("endDate")
		if startDateString == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "startDate is mandatory"})
			return
		}
		if endDateString == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "endDate is mandatory"})
			return
		}

		startDateTokens := strings.Split(startDateString, "-")
		endDateTokens := strings.Split(endDateString, "-")
		if len(startDateTokens) < 3 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed startDate"})
			return
		}
		if len(endDateTokens) < 3 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed endDate"})
			return
		}

		startDateYear, startDateYearErr := strconv.Atoi(startDateTokens[0])
		startDateMonth, startDateMonthErr := strconv.Atoi(startDateTokens[1])
		startDateDay, startDateDayErr := strconv.Atoi(startDateTokens[2])
		if startDateYearErr != nil || startDateMonthErr != nil || startDateDayErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "startDate must be numeric"})
			return
		}
		if startDateMonth < 1 || startDateMonth > 12 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "startDate month is invalid"})
			return
		}
		if startDateDay < 1 || startDateDay > 31 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "startDate day of month is invalid"})
			return
		}

		startDate := time.Date(startDateYear, time.Month(startDateMonth), startDateDay, 0, 0, 0, 0, time.UTC)

		endDateYear, endDateYearErr := strconv.Atoi(endDateTokens[0])
		endDateMonth, endDateMonthErr := strconv.Atoi(endDateTokens[1])
		endDateDay, endDateDayErr := strconv.Atoi(endDateTokens[2])
		if endDateYearErr != nil || endDateMonthErr != nil || endDateDayErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "endDate must be numeric"})
			return
		}
		if endDateMonth < 1 || endDateMonth > 12 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "endDate month is invalid"})
			return
		}
		if endDateDay < 1 || endDateDay > 31 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "endDate day of month is invalid"})
			return
		}

		endDate := time.Date(endDateYear, time.Month(endDateMonth), endDateDay, 0, 0, 0, 0, time.UTC)
		if endDate.Before(startDate) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "endDate can't be before startDate"})
			return
		}

		delta := endDate.Sub(startDate)

		// difference in days between startDate and endDate
		deltaDays := int(delta.Hours() / 24)
		if deltaDays > numMaxDays {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Maximum number of days exceeded"})
			return
		}

		// list of files to aggregate
		filePaths := make([]string, deltaDays+1)
		fileName := ""

		for i := 0; i <= deltaDays; i++ {
			date := startDate.AddDate(0, 0, i)
			year := strconv.Itoa(date.Year())
			month := fmt.Sprintf("%02d", int(date.Month()))
			day := fmt.Sprintf("%02d", date.Day())
			fileName = year + "/" + month + "/" + day + "/" + widgetName + ".json"
			fullPath := widgetDataRoot + fileName
			filePaths[i] = fullPath
		}
		fmt.Println("filePaths", filePaths) // todo del

		var widgetData map[string]interface{}

		var labelData Label

		var counterData Counter
		var valueOutputCounter float64

		var chartData Chart
		var numSeries int
		var numCategories int
		var seriesOutputChart []Series

		var tableData Table
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
			fmt.Println("filePath", filePath) // todo del
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

			fmt.Println("widgetData", widgetData) // todo del

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
					seriesOutputChart = make([]Series, numSeries)

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
			widget = chartData
		case "table":
			tableData.Rows = rowsOutputTable
			widget = tableData
		}

		c.JSON(http.StatusOK, gin.H{
			"widget": widget,
		})
	})
	router.Run(":8081") // listen and serve
}
