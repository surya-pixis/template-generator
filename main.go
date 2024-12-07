package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type GlobalTemplateConfig struct {
	Global Global `json:"global"`
}

type Global struct {
	TemplateID      string            `json:"template_id"`
	TemplateName    string            `json:"template_name"`
	TemplateConfigs []TemplateConfigs `json:"template_configs"`
}

type TemplateConfigs struct {
	TemplateConfigName string `json:"template_config_name"`
	TemplateType       string `json:"template_type"`
	BoardType          string `json:"board_type"`
	TemplateConfigID   string `json:"template_config_id"`
	Tabs               []Tab  `json:"tabs"`
}

type Tab struct {
	Title         string `json:"title"`
	SubTitle      string `json:"sub_title"`
	TemplateTabID string `json:"template_tab_id"`
	Grids         []Grid `json:"grids"`
}

type Grid struct {
	Title          string      `json:"title"`
	Position       int         `json:"position"`
	SubTitle       string      `json:"sub_title"`
	TemplateGridID string      `json:"template_grid_id"`
	Styling        GridStyling `json:"styling"`
	Charts         []Chart     `json:"charts"`
}

type GridStyling struct {
	TitleStyle    GridFontStyle `json:"titleStyle"`
	SubTitleStyle GridFontStyle `json:"subTitleStyle"`
}

type GridFontStyle struct {
	Font       string   `json:"font"`
	Color      string   `json:"color"`
	FontSize   int      `json:"font_size"`
	FontFormat []string `json:"font_format"`
}

type Chart struct {
	ChartType       string       `json:"chart_type"`
	Source          string       `json:"source"`
	Title           string       `json:"title"`
	TemplateChartID string       `json:"template_chart_id"`
	LeftMetrics     []Metric     `json:"left_metrics"`
	RightMetrics    []Metric     `json:"right_metrics,omitempty"`
	Dimensions      []Metric     `json:"dimensions,omitempty"`
	GridPosition    GridPos      `json:"grid_position"`
	Styling         ChartStyling `json:"styling"`
}

type ChartStyling struct {
	Palette        int              `json:"palette"`
	TitleStyle     ChartFontStyle   `json:"titleStyle"`
	TableStyle     TableTypeStyle   `json:"tableStyle"`
	LegendStyle    InsideTableStyle `json:"legendStyle"`
	LegendPosition string           `json:"legendPosition"`
}

type ChartFontStyle struct {
	Font       string   `json:"font"`
	Color      string   `json:"color"`
	FontSize   int      `json:"fontSize"`
	FontFormat []string `json:"fontFormat"`
	Alignment  string   `josn:"alignment"`
}

type TableTypeStyle struct {
	TableHeader  InsideTableStyle `json:"tableHeader"`
	TableContent InsideTableStyle `json:"tableContent"`
}

type InsideTableStyle struct {
	Font     string `json:"font"`
	FontSize int    `json:"fontSize"`
}

type Metric struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Path              string `json:"path"`
	Type              string `json:"type"`
	Group             string `json:"group"`
	Category          string `json:"category"`
	DataType          string `json:"dataType"`
	MetricType        string `json:"metricType"`
	Description       string `json:"description"`
	DivideByMillion   bool   `json:"divideByMillion"`
	AggregationMethod string `json:"aggregationMethod"`
}

type GridPos struct {
	H    int `json:"h"`
	W    int `json:"w"`
	X    int `json:"x"`
	Y    int `json:"y"`
	MaxH int `json:"maxH"`
	MinH int `json:"minH"`
	MinW int `json:"minW"`
}

func main() {
	// Google Sheet ID and credentials
	sheetID := "1ktTQ1scbWywG8oJZoLuqvHrZXZhtHwv3QtU92hv021w"
	credentialsFile := "credentials.json"

	// Initialize Sheets service
	sheetsService, err := sheets.NewService(context.Background(), option.WithCredentialsFile(credentialsFile))
	if err != nil {
		log.Fatalf("Unable to create Sheets service: %v", err)
	}

	// Initialize GlobalTemplateConfig
	finalTemplateConfig := GlobalTemplateConfig{
		Global: Global{
			TemplateID:   uuid.New().String(),
			TemplateName: "Generated Template",
		},
	}

	// Read data from Google Sheet
	readRange := "Sheet1!A4:J" // Adjust range if necessary
	data, err := sheetsService.Spreadsheets.Values.Get(sheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from Google Sheet: %v", err)
	}

	var dashboardTemplateConfigArray []TemplateConfigs
	dashboardTemplateConfigNumber := -1
	var reportTemplateConfigArray []TemplateConfigs
	reportTemplateConfigNumber := -1
	// var dashboardTabsArray []Tab
	// dashboardTabNumber := 0
	// var reportTabsArray []Tab
	// reportTabNumber := 0

	currentDashboardTemplateConfig := &TemplateConfigs{}
	currentReportTemplateConfig := &TemplateConfigs{}
	currentTab := &Tab{}
	currentGrid := &Grid{}
	currentChart := &Chart{}

	boardOfType := ""

	// Parse sheet data
	for _, row := range data.Values {
		if len(row) == 0 {
			continue
		}

		//intCheck := 0

		if row[0] != "" {
			if boardOfType == "" { //this condition is for starting the process and setting the context to the first board type found
				if !strings.Contains(fmt.Sprint(row[0]), "Report") {
					boardOfType = "DASHBOARD"
					currentDashboardTemplateConfig = &TemplateConfigs{
						BoardType:        boardOfType,
						TemplateConfigID: uuid.New().String(),
						TemplateType:     "TAB_GRID_CHART",
					}
					dashboardTemplateConfigArray = append(dashboardTemplateConfigArray, *currentDashboardTemplateConfig)
					dashboardTemplateConfigNumber++
				} else if strings.Contains(fmt.Sprint(row[0]), "Report") {
					boardOfType = "REPORT"
					currentReportTemplateConfig = &TemplateConfigs{
						BoardType:        boardOfType,
						TemplateConfigID: uuid.New().String(),
						TemplateType:     "TAB_CHART",
					}
					reportTemplateConfigArray = append(reportTemplateConfigArray, *currentReportTemplateConfig)
					reportTemplateConfigNumber++
				}

			} else if boardOfType == "DASHBOARD" && !strings.Contains(fmt.Sprint(row[0]), "Report") { // this check is to continue the dashboard context
				//we append the tab without changing the context
				fmt.Println("Same DASHBOARD Template Continuation")
				fmt.Println("Appending Dashboard Tab")
				dashboardTemplateConfigArray[dashboardTemplateConfigNumber].Tabs = append(dashboardTemplateConfigArray[dashboardTemplateConfigNumber].Tabs, *currentTab)
				//dashboardTabNumber++
			} else if boardOfType == "REPORT" && strings.Contains(fmt.Sprint(row[0]), "Report") { // this check is to continue the report context
				//we append the tab without changing the context
				fmt.Println("Same REPORT Template Continuation")
				fmt.Println("Appending Report Tab")
				reportTemplateConfigArray[reportTemplateConfigNumber].Tabs = append(reportTemplateConfigArray[reportTemplateConfigNumber].Tabs, *currentTab)
				//reportTabNumber++
			} else if boardOfType == "DASHBOARD" && strings.Contains(fmt.Sprint(row[0]), "Report") { // this check is to change the dashboard context to report context
				//we need to change the context here
				fmt.Println("Changing from DASHBOARD Template to REPORT Template")
				fmt.Println("Appending Dashboard Tab")
				dashboardTemplateConfigArray[dashboardTemplateConfigNumber].Tabs = append(dashboardTemplateConfigArray[dashboardTemplateConfigNumber].Tabs, *currentTab)
				//changing the context here
				boardOfType = "REPORT"
				currentReportTemplateConfig = &TemplateConfigs{
					BoardType:        boardOfType,
					TemplateConfigID: uuid.New().String(),
					TemplateType:     "TAB_CHART",
				}
				reportTemplateConfigArray = append(reportTemplateConfigArray, *currentReportTemplateConfig)
				reportTemplateConfigNumber++
			} else if boardOfType == "REPORT" && !strings.Contains(fmt.Sprint(row[0]), "Report") { // this check is to change the report context to dashboard context
				//we need to change the context here
				fmt.Println("Changing from REPORT Template to DASHBOARD Template")
				fmt.Println("Appending Report Tab")
				reportTemplateConfigArray[reportTemplateConfigNumber].Tabs = append(reportTemplateConfigArray[reportTemplateConfigNumber].Tabs, *currentTab)
				//changing the context here
				boardOfType = "DASHBOARD"
				currentDashboardTemplateConfig = &TemplateConfigs{
					BoardType:        boardOfType,
					TemplateConfigID: uuid.New().String(),
					TemplateType:     "TAB_GRID_CHART",
				}
				dashboardTemplateConfigArray = append(dashboardTemplateConfigArray, *currentDashboardTemplateConfig)
				dashboardTemplateConfigNumber++
			}

			//creating a new tab whenever we come accross it in column 1
			newTab := &Tab{
				Title:         fmt.Sprint(row[0]),
				TemplateTabID: uuid.New().String(),
			}
			currentTab = newTab

		}

		//handling the grids when we come across it

		//when we get a new grid at the start
		if len(row) < 2 {
			continue
		}
		if row[1] != "" {

			currentTab.Grids = append(currentTab.Grids, *currentGrid)

			newGrid := &Grid{
				Title:          fmt.Sprint(row[1]),
				TemplateGridID: uuid.New().String(),
			}
			currentGrid = newGrid
			//handling charts
			if len(row) < 3 {
				continue
			}
			if row[2] != "" {

				currentGrid.Charts = append(currentGrid.Charts, *currentChart)

				newChart := &Chart{
					TemplateChartID: uuid.New().String(),
					ChartType:       fmt.Sprint(row[2]),
					Title:           fmt.Sprint(row[3]),
				}
				currentChart = newChart
				// handling first chart dimension
				if row[4] != "" {
					dimension := &Metric{
						Name: fmt.Sprint(row[4]),
						ID:   fmt.Sprint(row[5]),
					}
					currentChart.Dimensions = append(currentChart.Dimensions, *dimension)
				}

				//handling chart first left/right metric
				if row[6] != "" {
					metric := &Metric{
						Name: fmt.Sprint(row[6]),
						ID:   fmt.Sprint(row[7]),
					}
					currentChart.LeftMetrics = append(currentChart.LeftMetrics, *metric)
				}

			} else if row[2] == "" {
				if row[4] != "" {
					dimension := &Metric{
						Name: fmt.Sprint(row[4]),
						ID:   fmt.Sprint(row[5]),
					}
					currentChart.Dimensions = append(currentChart.Dimensions, *dimension)
				}
				if row[6] != "" && currentChart.ChartType == "Line" {
					metric := &Metric{
						Name: fmt.Sprint(row[6]),
						ID:   fmt.Sprint(row[7]),
					}
					currentChart.RightMetrics = append(currentChart.RightMetrics, *metric)

				} else if row[6] != "" && currentChart.ChartType != "Line" {
					metric := &Metric{
						Name: fmt.Sprint(row[6]),
						ID:   fmt.Sprint(row[7]),
					}
					currentChart.LeftMetrics = append(currentChart.LeftMetrics, *metric)
				}
			}

		} else if row[1] == "" {
			//handling charts
			if row[2] != "" {

				currentGrid.Charts = append(currentGrid.Charts, *currentChart)

				newChart := &Chart{
					TemplateChartID: uuid.New().String(),
					ChartType:       fmt.Sprint(row[2]),
					Title:           fmt.Sprint(row[3]),
				}
				currentChart = newChart
				// handling first chart dimension
				if row[4] != "" {
					dimension := &Metric{
						Name: fmt.Sprint(row[4]),
						ID:   fmt.Sprint(row[5]),
					}
					currentChart.Dimensions = append(currentChart.Dimensions, *dimension)
				}

				//handling chart first left/right metric
				if row[6] != "" {
					metric := &Metric{
						Name: fmt.Sprint(row[6]),
						ID:   fmt.Sprint(row[7]),
					}
					currentChart.LeftMetrics = append(currentChart.LeftMetrics, *metric)
				}

			} else if row[2] == "" {
				if row[4] != "" {
					dimension := &Metric{
						Name: fmt.Sprint(row[4]),
						ID:   fmt.Sprint(row[5]),
					}
					currentChart.Dimensions = append(currentChart.Dimensions, *dimension)
				}
				if row[6] != "" && currentChart.ChartType == "Line" {
					metric := &Metric{
						Name: fmt.Sprint(row[6]),
						ID:   fmt.Sprint(row[7]),
					}
					currentChart.RightMetrics = append(currentChart.RightMetrics, *metric)

				} else if row[6] != "" && currentChart.ChartType != "Line" {
					metric := &Metric{
						Name: fmt.Sprint(row[6]),
						ID:   fmt.Sprint(row[7]),
					}
					currentChart.LeftMetrics = append(currentChart.LeftMetrics, *metric)
				}
			}
		}

		//intCheck++

	}

	//we added this piece of code at the end to add the last tab which may be a report or a dashboard
	if boardOfType == "DASHBOARD" { // appending the last tab with respect to context
		fmt.Println("Appending Last Tab which is a DASHBOARD")
		dashboardTemplateConfigArray[dashboardTemplateConfigNumber].Tabs = append(dashboardTemplateConfigArray[dashboardTemplateConfigNumber].Tabs, *currentTab)
	} else if boardOfType == "REPORT" { // appending the last tab with respect to context
		fmt.Println("Appending Last Tab which is a REPORT")
		reportTemplateConfigArray[reportTemplateConfigNumber].Tabs = append(reportTemplateConfigArray[reportTemplateConfigNumber].Tabs, *currentTab)
	}

	finalTemplateConfig.Global.TemplateConfigs = append(finalTemplateConfig.Global.TemplateConfigs, dashboardTemplateConfigArray...)
	finalTemplateConfig.Global.TemplateConfigs = append(finalTemplateConfig.Global.TemplateConfigs, reportTemplateConfigArray...)

	// Write JSON to file
	outputFile, err := os.Create("output_template.json")
	if err != nil {
		log.Fatalf("Unable to create output file: %v", err)
	}
	defer outputFile.Close()

	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(finalTemplateConfig); err != nil {
		log.Fatalf("Unable to write JSON to file: %v", err)
	}

	fmt.Println("Template JSON generated successfully!")
}

// above is the old logic

// func main() {
// 	// Google Sheet ID and credentials
// 	sheetID := "1ktTQ1scbWywG8oJZoLuqvHrZXZhtHwv3QtU92hv021w"
// 	credentialsFile := "credentials.json"

// 	// Initialize Sheets service
// 	sheetsService, err := sheets.NewService(context.Background(), option.WithCredentialsFile(credentialsFile))
// 	if err != nil {
// 		log.Fatalf("Unable to create Sheets service: %v", err)
// 	}

// 	// Initialize GlobalTemplateConfig
// 	finalTemplateConfig := GlobalTemplateConfig{
// 		Global: Global{
// 			TemplateID:   uuid.New().String(),
// 			TemplateName: "Generated Template",
// 		},
// 	}

// 	// Read data from Google Sheet
// 	readRange := "Sheet1!A4:J" // Adjust range if necessary
// 	data, err := sheetsService.Spreadsheets.Values.Get(sheetID, readRange).Do()
// 	if err != nil {
// 		log.Fatalf("Unable to retrieve data from Google Sheet: %v", err)
// 	}

// 	var dashboardTemplateConfigArray []TemplateConfigs
// 	dashboardTemplateConfigNumber := -1
// 	var reportTemplateConfigArray []TemplateConfigs
// 	reportTemplateConfigNumber := -1
// 	// var dashboardTabsArray []Tab
// 	// dashboardTabNumber := 0
// 	// var reportTabsArray []Tab
// 	// reportTabNumber := 0

// 	currentDashboardTemplateConfig := &TemplateConfigs{}
// 	currentReportTemplateConfig := &TemplateConfigs{}
// 	currentTab := &Tab{}
// 	currentGrid := &Grid{}
// 	currentChart := &Chart{}

// 	boardOfType := ""

// 	// 	var let i = 0;
// 	// for (i; i < len(rows); i++) {
// 	//     if (row[i][0] != "") {
// 	//         tab = []
// 	//         let j = i;
// 	//         for (j; j < len(rows); j++) {
// 	//             if (row[j][1] != "") {
// 	//                 grid = []
// 	//                 for (k = j; k < len(rows); k++) {
// 	//                     if (row[k][2] != "" && row[k][0] == "" && row[k][1] = "") {
// 	//                         grid.charts.append(chart)
// 	//                     } else {
// 	//                         j = k
// 	//                         // new gird or tab to be created
// 	//                     }
// 	//                 }
// 	//             } else {
// 	//                 tab.grid.append(grid)
// 	//             }
// 	//         } else {
// 	//             config.tabs.append(tab)
// 	//         }
// 	//     }
// 	// }

// 	spreadSheet := data.Values

// 	// Parse sheet data
// 	for i := 0; i < len(spreadSheet); i++ {
// 		if len(spreadSheet[i]) == 0 {
// 			continue
// 		}

// 		if boardOfType == "" { //this condition is for starting the process and setting the context to the first board type found
// 			if !strings.Contains(fmt.Sprint(row[0]), "Report") {
// 				boardOfType = "DASHBOARD"
// 				currentDashboardTemplateConfig = &TemplateConfigs{
// 					BoardType:        boardOfType,
// 					TemplateConfigID: uuid.New().String(),
// 					TemplateType:     "TAB_GRID_CHART",
// 				}
// 				dashboardTemplateConfigArray = append(dashboardTemplateConfigArray, *currentDashboardTemplateConfig)
// 				dashboardTemplateConfigNumber++
// 			} else if strings.Contains(fmt.Sprint(row[0]), "Report") {
// 				boardOfType = "REPORT"
// 				currentReportTemplateConfig = &TemplateConfigs{
// 					BoardType:        boardOfType,
// 					TemplateConfigID: uuid.New().String(),
// 					TemplateType:     "TAB_CHART",
// 				}
// 				reportTemplateConfigArray = append(reportTemplateConfigArray, *currentReportTemplateConfig)
// 				reportTemplateConfigNumber++
// 			}

// 		} else if boardOfType == "DASHBOARD" && !strings.Contains(fmt.Sprint(row[0]), "Report") { // this check is to continue the dashboard context
// 			//we append the tab without changing the context
// 			fmt.Println("Same DASHBOARD Template Continuation")
// 			fmt.Println("Appending Dashboard Tab")
// 			dashboardTemplateConfigArray[dashboardTemplateConfigNumber].Tabs = append(dashboardTemplateConfigArray[dashboardTemplateConfigNumber].Tabs, *currentTab)
// 			//dashboardTabNumber++
// 		} else if boardOfType == "REPORT" && strings.Contains(fmt.Sprint(row[0]), "Report") { // this check is to continue the report context
// 			//we append the tab without changing the context
// 			fmt.Println("Same REPORT Template Continuation")
// 			fmt.Println("Appending Report Tab")
// 			reportTemplateConfigArray[reportTemplateConfigNumber].Tabs = append(reportTemplateConfigArray[reportTemplateConfigNumber].Tabs, *currentTab)
// 			//reportTabNumber++
// 		} else if boardOfType == "DASHBOARD" && strings.Contains(fmt.Sprint(row[0]), "Report") { // this check is to change the dashboard context to report context
// 			//we need to change the context here
// 			fmt.Println("Changing from DASHBOARD Template to REPORT Template")
// 			fmt.Println("Appending Dashboard Tab")
// 			dashboardTemplateConfigArray[dashboardTemplateConfigNumber].Tabs = append(dashboardTemplateConfigArray[dashboardTemplateConfigNumber].Tabs, *currentTab)
// 			//changing the context here
// 			boardOfType = "REPORT"
// 			currentReportTemplateConfig = &TemplateConfigs{
// 				BoardType:        boardOfType,
// 				TemplateConfigID: uuid.New().String(),
// 				TemplateType:     "TAB_CHART",
// 			}
// 			reportTemplateConfigArray = append(reportTemplateConfigArray, *currentReportTemplateConfig)
// 			reportTemplateConfigNumber++
// 		} else if boardOfType == "REPORT" && !strings.Contains(fmt.Sprint(row[0]), "Report") { // this check is to change the report context to dashboard context
// 			//we need to change the context here
// 			fmt.Println("Changing from REPORT Template to DASHBOARD Template")
// 			fmt.Println("Appending Report Tab")
// 			reportTemplateConfigArray[reportTemplateConfigNumber].Tabs = append(reportTemplateConfigArray[reportTemplateConfigNumber].Tabs, *currentTab)
// 			//changing the context here
// 			boardOfType = "DASHBOARD"
// 			currentDashboardTemplateConfig = &TemplateConfigs{
// 				BoardType:        boardOfType,
// 				TemplateConfigID: uuid.New().String(),
// 				TemplateType:     "TAB_GRID_CHART",
// 			}
// 			dashboardTemplateConfigArray = append(dashboardTemplateConfigArray, *currentDashboardTemplateConfig)
// 			dashboardTemplateConfigNumber++
// 		}

// 		//handling the grids when we come across it

// 		//when we get a new grid at the start
// 		if len(row) < 2 {
// 			continue
// 		}
// 		if row[1] != "" {

// 			currentTab.Grids = append(currentTab.Grids, *currentGrid)

// 			newGrid := &Grid{
// 				Title:          fmt.Sprint(row[1]),
// 				TemplateGridID: uuid.New().String(),
// 			}
// 			currentGrid = newGrid
// 			//handling charts
// 			if len(row) < 3 {
// 				continue
// 			}
// 			if row[2] != "" {

// 				currentGrid.Charts = append(currentGrid.Charts, *currentChart)

// 				newChart := &Chart{
// 					TemplateChartID: uuid.New().String(),
// 					ChartType:       fmt.Sprint(row[2]),
// 					Title:           fmt.Sprint(row[3]),
// 				}
// 				currentChart = newChart
// 				// handling first chart dimension
// 				if row[4] != "" {
// 					dimension := &Metric{
// 						Name: fmt.Sprint(row[4]),
// 						ID:   fmt.Sprint(row[5]),
// 					}
// 					currentChart.Dimensions = append(currentChart.Dimensions, *dimension)
// 				}

// 				//handling chart first left/right metric
// 				if row[6] != "" {
// 					metric := &Metric{
// 						Name: fmt.Sprint(row[6]),
// 						ID:   fmt.Sprint(row[7]),
// 					}
// 					currentChart.LeftMetrics = append(currentChart.LeftMetrics, *metric)
// 				}

// 			} else if row[2] == "" {
// 				if row[4] != "" {
// 					dimension := &Metric{
// 						Name: fmt.Sprint(row[4]),
// 						ID:   fmt.Sprint(row[5]),
// 					}
// 					currentChart.Dimensions = append(currentChart.Dimensions, *dimension)
// 				}
// 				if row[6] != "" && currentChart.ChartType == "Line" {
// 					metric := &Metric{
// 						Name: fmt.Sprint(row[6]),
// 						ID:   fmt.Sprint(row[7]),
// 					}
// 					currentChart.RightMetrics = append(currentChart.RightMetrics, *metric)

// 				} else if row[6] != "" && currentChart.ChartType != "Line" {
// 					metric := &Metric{
// 						Name: fmt.Sprint(row[6]),
// 						ID:   fmt.Sprint(row[7]),
// 					}
// 					currentChart.LeftMetrics = append(currentChart.LeftMetrics, *metric)
// 				}
// 			}

// 		} else if row[1] == "" {
// 			//handling charts
// 			if row[2] != "" {

// 				currentGrid.Charts = append(currentGrid.Charts, *currentChart)

// 				newChart := &Chart{
// 					TemplateChartID: uuid.New().String(),
// 					ChartType:       fmt.Sprint(row[2]),
// 					Title:           fmt.Sprint(row[3]),
// 				}
// 				currentChart = newChart
// 				// handling first chart dimension
// 				if row[4] != "" {
// 					dimension := &Metric{
// 						Name: fmt.Sprint(row[4]),
// 						ID:   fmt.Sprint(row[5]),
// 					}
// 					currentChart.Dimensions = append(currentChart.Dimensions, *dimension)
// 				}

// 				//handling chart first left/right metric
// 				if row[6] != "" {
// 					metric := &Metric{
// 						Name: fmt.Sprint(row[6]),
// 						ID:   fmt.Sprint(row[7]),
// 					}
// 					currentChart.LeftMetrics = append(currentChart.LeftMetrics, *metric)
// 				}

// 			} else if row[2] == "" {
// 				if row[4] != "" {
// 					dimension := &Metric{
// 						Name: fmt.Sprint(row[4]),
// 						ID:   fmt.Sprint(row[5]),
// 					}
// 					currentChart.Dimensions = append(currentChart.Dimensions, *dimension)
// 				}
// 				if row[6] != "" && currentChart.ChartType == "Line" {
// 					metric := &Metric{
// 						Name: fmt.Sprint(row[6]),
// 						ID:   fmt.Sprint(row[7]),
// 					}
// 					currentChart.RightMetrics = append(currentChart.RightMetrics, *metric)

// 				} else if row[6] != "" && currentChart.ChartType != "Line" {
// 					metric := &Metric{
// 						Name: fmt.Sprint(row[6]),
// 						ID:   fmt.Sprint(row[7]),
// 					}
// 					currentChart.LeftMetrics = append(currentChart.LeftMetrics, *metric)
// 				}
// 			}
// 		}

// 		//intCheck++

// 	}

// 	//we added this piece of code at the end to add the last tab which may be a report or a dashboard
// 	if boardOfType == "DASHBOARD" { // appending the last tab with respect to context
// 		fmt.Println("Appending Last Tab which is a DASHBOARD")
// 		dashboardTemplateConfigArray[dashboardTemplateConfigNumber].Tabs = append(dashboardTemplateConfigArray[dashboardTemplateConfigNumber].Tabs, *currentTab)
// 	} else if boardOfType == "REPORT" { // appending the last tab with respect to context
// 		fmt.Println("Appending Last Tab which is a REPORT")
// 		reportTemplateConfigArray[reportTemplateConfigNumber].Tabs = append(reportTemplateConfigArray[reportTemplateConfigNumber].Tabs, *currentTab)
// 	}

// 	finalTemplateConfig.Global.TemplateConfigs = append(finalTemplateConfig.Global.TemplateConfigs, dashboardTemplateConfigArray...)
// 	finalTemplateConfig.Global.TemplateConfigs = append(finalTemplateConfig.Global.TemplateConfigs, reportTemplateConfigArray...)

// 	// Write JSON to file
// 	outputFile, err := os.Create("output_template.json")
// 	if err != nil {
// 		log.Fatalf("Unable to create output file: %v", err)
// 	}
// 	defer outputFile.Close()

// 	encoder := json.NewEncoder(outputFile)
// 	encoder.SetIndent("", "  ")
// 	if err := encoder.Encode(finalTemplateConfig); err != nil {
// 		log.Fatalf("Unable to write JSON to file: %v", err)
// 	}

// 	fmt.Println("Template JSON generated successfully!")
// }
