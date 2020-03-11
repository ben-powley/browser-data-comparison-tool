package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ahmetb/go-linq"
	"github.com/ben-powley/browser-data-comparison-tool/models"
	"github.com/ben-powley/browser-data-comparison-tool/utils"
	csvreader "github.com/ben-powley/go-csv-reader"
)

var _totalSessions int
var _totalTransactions int

func main() {
	start := time.Now()

	fmt.Println("Start time: ", start.String())

	if _, err := os.Stat("data"); os.IsNotExist(err) {
		panic(errors.New("Data folder not found - Please add a folder named 'data' with csv files"))
	}

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

	exportValues(outputModels, browserData, &start)

	end := time.Since(start)

	fmt.Println("Execution time: ", end)
}

func exportValues(outputModels []models.OutputModel, browserData []models.BrowserData, start *time.Time) {
	browserNames := []string{
		"chrome",
		"safari",
		"safari 12+",
		"internet explorer",
		"edge",
		"firefox",
		"samsung internet",
		"opera",
		"amazon silk",
		"mozilla compatible agent",
		"android webview",
		"safari (in-app)",
		"uc browser",
		"e.ventures investment crawler",
		"(not set)",
		`''`,
	}
	countModels := []models.CountModel{}

	totalUsers := 0
	totalReturningSessions := 0
	totalReturningTransactions := 0

	for _, browser := range browserNames {
		shouldFilter := false

		if browser == "safari 12+" {
			shouldFilter = true
		}

		totalUsersForBrowser := getTotalUsersForBrowser(browser, browserData, shouldFilter)
		outputModelsForBrowser := getOutputModelsForBrowser(browser, outputModels, shouldFilter)
		averageDaysForBrowser := getAverageDaysForOutputModels(outputModelsForBrowser)
		sessionsForBrowser := getSessionsForBrowser(outputModelsForBrowser)
		transactionsForBrowser := getTransactionsForBrowser(outputModelsForBrowser)
		averageReturnsForBrowser := getAverageReturnsForBrowser(outputModelsForBrowser)
		browserCountModel := models.CountModel{
			Browser:                 strings.ToTitle(browser),
			TotalUsers:              totalUsersForBrowser,
			AverageDaysBetweenVisit: averageDaysForBrowser,
			ReturningUserTotal:      len(outputModelsForBrowser),
			SessionsTotal:           sessionsForBrowser,
			TransactionsTotal:       transactionsForBrowser,
			AverageReturns:          averageReturnsForBrowser,
		}

		totalUsers += totalUsersForBrowser
		totalReturningSessions += sessionsForBrowser
		totalReturningTransactions += transactionsForBrowser

		countModels = append(countModels, browserCountModel)
	}

	fmt.Println("--- --- --- --- --- ---")

	fmt.Println("Total records: ", len(browserData))
	fmt.Println("Total users: ", totalUsers)
	fmt.Println("Total sessions: ", _totalSessions)
	fmt.Println("Total returning users: ", len(outputModels))
	fmt.Println("Total transactions: ", _totalTransactions)
	fmt.Println("Total returning user sessions: ", totalReturningSessions)
	fmt.Println("Total returning user transactions: ", totalReturningTransactions)

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

func getTotalUsersForBrowser(browser string, browserData []models.BrowserData, filterBySafari12 bool) int {
	if filterBySafari12 {
		browser = "safari"
	}

	totalForBrowser := linq.From(browserData).Where(func(bd interface{}) bool {
		bdConverted := bd.(models.BrowserData)

		if strings.ToLower(bdConverted.Browser) == browser {
			if browser == "safari" && filterBySafari12 {
				browserVersion := bdConverted.BrowserVersion

				if strings.HasPrefix(browserVersion, "12") || strings.HasPrefix(browserVersion, "13") {
					return true
				}
			} else {
				return true
			}
		}

		return false
	}).Count()

	return totalForBrowser
}

func getAverageReturnsForBrowser(outputModels []models.OutputModel) int {
	averageReturns := 0

	if len(outputModels) > 0 {
		for _, om := range outputModels {
			averageReturns += om.Returns
		}

		averageReturns = averageReturns / len(outputModels)
	}

	return averageReturns
}

func getAverageDaysForOutputModels(outputModels []models.OutputModel) int {
	averageDays := 0

	if len(outputModels) == 0 {
		for _, om := range outputModels {
			averageDays += om.AverageDaysBetweenVisit
		}

		if averageDays > 0 {
			averageDays = averageDays / len(outputModels)
		}
	}

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

	fmt.Println(browser+" total users: ", countModel.TotalUsers)
	fmt.Println(browser+" returning users: ", countModel.ReturningUserTotal)
	fmt.Println(browser+" average days between visits: ", countModel.AverageDaysBetweenVisit)
	fmt.Println(browser+" total returning user sessions: ", countModel.SessionsTotal)
	fmt.Println(browser+" total returning user transactions: ", countModel.TransactionsTotal)
	fmt.Println(browser+" average number of returns: ", countModel.AverageReturns)

	fmt.Println("--- --- --- --- --- ---")
}

func getOuputModels(groupQuery []linq.Group) []models.OutputModel {
	outputModels := []models.OutputModel{}

	for _, gq := range groupQuery {
		totalSessions := 0
		totalTransactions := 0
		returns := 0
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

		returns += groupLength

		if len(days) == 0 {
			averageDays = 0
		} else if len(days) == 1 {
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
			Returns:                 returns,
		})
	}

	return outputModels
}

func filterBrowserData(browserData []models.BrowserData) []linq.Group {
	var groupQuery []linq.Group

	for _, bd := range browserData {
		sessions, _ := strconv.Atoi(bd.Sessions)
		transactions, _ := strconv.Atoi(bd.Transactions)

		_totalSessions += sessions
		_totalTransactions += transactions
	}

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
