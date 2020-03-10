package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ahmetb/go-linq"
	"github.com/ben-powley/browser-data-comparison-tool/models"
	"github.com/ben-powley/browser-data-comparison-tool/utils"
	csvreader "github.com/ben-powley/go-csv-reader"
)

func main() {
	start := time.Now()

	fmt.Println("Start time: ", start.String())

	filenames, filenamesErr := utils.GetFilenamesFromFolder("//data")

	if filenamesErr != nil {
		panic(filenamesErr)
	}

	if len(filenames) == 0 {
		panic(errors.New("Could not find any .CSV files in directory"))
	}

	csvData, csvDataErr := csvreader.ReadCSVFiles(filenames, true)

	if csvDataErr != nil {
		panic(csvDataErr)
	}

	browserData := utils.ConvertCSVLinesToBrowserData(csvData)

	if browserData == nil {
		panic(errors.New("Browser data is empty"))
	}

	groupQuery := filterBrowserData(browserData)
	outputModels := getOuputModels(groupQuery)

	if len(outputModels) == 0 {
		panic(errors.New("Output models is empty"))
	}

	exportValues(outputModels, len(browserData), &start)

	end := time.Since(start)

	fmt.Println("Execution time: ", end)
}

func exportValues(outputModels []models.OutputModel, browserDataLen int, start *time.Time) {
	fmt.Println("--- --- --- --- --- ---")

	fmt.Println("Total records: ", browserDataLen)
	fmt.Println("Total returning users: ", len(outputModels))

	fmt.Println("--- --- --- --- --- ---")

	browserNames := []string{
		"chrome",
		"safari",
		"safari 12+",
		"internet explorer",
		"edge",
		"firefox",
		"samsung internet",
	}
	countModels := []models.CountModel{}

	totalSessions := 0
	totalTransactions := 0

	for _, browser := range browserNames {
		shouldFilter := false

		if browser == "safari 12+" {
			shouldFilter = true
		}

		outputModelsForBrowser := getOutputModelsForBrowser(browser, outputModels, shouldFilter)
		averageDaysForBrowser := getAverageDaysForOutputModels(outputModelsForBrowser)
		sessionsForBrowser := getSessionsForBrowser(outputModelsForBrowser)
		transactionsForBrowser := getTransactionsForBrowser(outputModelsForBrowser)
		browserCountModel := models.CountModel{
			Browser:                 strings.ToTitle(browser),
			AverageDaysBetweenVisit: averageDaysForBrowser,
			ReturningUserTotal:      len(outputModelsForBrowser),
			SessionsTotal:           sessionsForBrowser,
			TransactionsTotal:       transactionsForBrowser,
		}

		totalSessions += sessionsForBrowser
		totalTransactions += transactionsForBrowser

		countModels = append(countModels, browserCountModel)
	}

	fmt.Println("Total sessions: ", totalSessions)
	fmt.Println("Total transactions: ", totalTransactions)

	fmt.Println("--- --- --- --- --- ---")

	for _, cm := range countModels {
		outputValuesToConsole(cm)
	}
}

func getOutputModelsForBrowser(browser string, outputModels []models.OutputModel, filterBySafari12 bool) []models.OutputModel {
	om := []models.OutputModel{}

	linq.From(outputModels).Where(func(om interface{}) bool {
		browserName := om.(models.OutputModel).Group.Group[0].(models.BrowserData).Browser

		if filterBySafari12 {
			browser = "safari"
		}

		if strings.ToLower(browserName) == browser {
			if browser == "safari" && filterBySafari12 {
				browserVersion := om.(models.OutputModel).Group.Group[0].(models.BrowserData).BrowserVersion

				if strings.HasPrefix(browserVersion, "12") || strings.HasPrefix(browserVersion, "13") {
					return true
				}
			} else {
				return true
			}
		}

		return false
	}).ToSlice(&om)

	return om
}

func getAverageDaysForOutputModels(outputModels []models.OutputModel) int {
	averageDays := 0

	for _, om := range outputModels {
		averageDays += om.AverageDaysBetweenVisit
	}

	if len(outputModels) == 0 {
		return averageDays
	}

	averageDays = averageDays / len(outputModels)

	return averageDays
}

func getSessionsForBrowser(outputModels []models.OutputModel) int {
	sessions := 0

	for _, om := range outputModels {
		sessions += om.TotalSessions
	}

	return sessions
}

func getTransactionsForBrowser(outputModels []models.OutputModel) int {
	transactions := 0

	for _, om := range outputModels {
		transactions += om.TotalTransactions
	}

	return transactions
}

func outputValuesToConsole(countModel models.CountModel) {
	browser := countModel.Browser

	fmt.Println(browser+" returning users: ", countModel.ReturningUserTotal)
	fmt.Println(browser+" average days between visits: ", countModel.AverageDaysBetweenVisit)
	fmt.Println(browser+" total sessions: ", countModel.SessionsTotal)
	fmt.Println(browser+" total transactions: ", countModel.TransactionsTotal)

	fmt.Println("--- --- --- --- --- ---")
}

func getOuputModels(groupQuery []linq.Group) []models.OutputModel {
	outputModels := []models.OutputModel{}

	for _, gq := range groupQuery {
		totalSessions := 0
		totalTransactions := 0
		averageDays := 0
		days := []int{}
		groupLength := len(gq.Group)

		for i := 0; i <= groupLength-1; i++ {
			if i > 0 {
				currentDate := gq.Group[i].(models.BrowserData).FormattedDate
				previousDate := gq.Group[i-1].(models.BrowserData).FormattedDate

				daysBetween := int(currentDate.Sub(previousDate).Hours() / 24)

				days = append(days, daysBetween)
			}

			sessions, _ := strconv.Atoi(gq.Group[i].(models.BrowserData).Sessions)
			transactions, _ := strconv.Atoi(gq.Group[i].(models.BrowserData).Transactions)

			totalSessions += sessions
			totalTransactions += transactions
		}

		if len(days) == 1 {
			averageDays = days[0]
		} else {
			sumDays := 0

			for i := range days {
				sumDays += days[i]
			}

			averageDays = sumDays / len(days)
		}

		outputModels = append(outputModels, models.OutputModel{
			Group:                   gq,
			AverageDaysBetweenVisit: averageDays,
			TotalSessions:           totalSessions,
			TotalTransactions:       totalTransactions,
		})
	}

	return outputModels
}

func filterBrowserData(browserData []models.BrowserData) []linq.Group {
	var groupQuery []linq.Group

	// Group browser data by Client Ids
	linq.From(browserData).GroupBy(func(bd interface{}) interface{} {
		return bd.(models.BrowserData).ClientID
	}, func(bd interface{}) interface{} {
		return bd.(models.BrowserData)
	}).Where(func(group interface{}) bool {
		convertedGroup := group.(linq.Group)

		if len(convertedGroup.Group) > 1 {
			return true
		}

		return false
	}).ToSlice(&groupQuery)

	return groupQuery
}
