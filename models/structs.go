package models

import (
	"time"

	. "github.com/ahmetb/go-linq"
)

type BrowserData struct {
	ID                 int
	Date               string
	ClientID           string
	DeviceCategory     string
	Browser            string
	BrowserVersion     string
	Sessions           string
	Transactions       string
	TransactionRevenue string
	FormattedDate      time.Time
}

type ClientModel struct {
	ClientID string
	Browser  string
}

type CSVRow struct {
	Title string
	Count string
}

type OutputModel struct {
	Browser                 string
	Group                   Group
	AverageDaysBetweenVisit int
}

type CountModel struct {
	Browser                 string
	ReturningUserTotal      int
	AverageDaysBetweenVisit int
}
