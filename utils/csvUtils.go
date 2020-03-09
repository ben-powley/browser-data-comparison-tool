package utils

import (
	"errors"
	"strconv"
	"strings"

	"github.com/ben-powley/browser-data-comparison-tool/models"
)

func ConvertCSVLinesToBrowserData(csvLines [][]string) []models.BrowserData {
	browserData := []models.BrowserData{}

	for i, csvLine := range csvLines {
		bd := models.BrowserData{
			ID:                 i,
			Date:               csvLine[1],
			ClientID:           csvLine[2],
			DeviceCategory:     csvLine[3],
			Browser:            csvLine[4],
			BrowserVersion:     csvLine[5],
			Sessions:           csvLine[6],
			Transactions:       csvLine[7],
			TransactionRevenue: csvLine[8],
		}

		date := bd.Date
		splitDate := strings.Split(date, "-")

		if len(splitDate) == 3 {
			year, yearErr := strconv.Atoi(splitDate[0])
			month, monthErr := strconv.Atoi(splitDate[1])
			day, dayErr := strconv.Atoi(splitDate[2])

			if yearErr != nil || monthErr != nil || dayErr != nil {
				panic(errors.New("Error converting date"))
			}

			formattedDate := Date(year, month, day)

			bd.FormattedDate = formattedDate
		}

		browserData = append(browserData, bd)
	}

	return browserData
}
