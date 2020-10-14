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

type TotalsModel struct {
	Users        int
	Sessions     int
	Transactions int
}

type OutputModel struct {
	Browser                 string
	Group                   Group
	AverageDaysBetweenVisit int
	TotalSessions           int
	TotalTransactions       int
	Returns                 int
}

type CountModel struct {
	Browser                    string
	TotalUsers                 int
	ReturningUserTotal         int
	AverageDaysBetweenVisit    int
	SessionsTotal              int
	TransactionsTotal          int
	ReturningSessionsTotal     int
	ReturningTransactionsTotal int
	AverageReturns             int
}
