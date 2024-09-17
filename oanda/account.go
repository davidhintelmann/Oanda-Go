package oanda

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

/*
struct for unmarshalling json from [Account Endpoints] for which one is authorized with the provided token.

endpoint: /v3/accounts

[Account Endpoints]: https://developer.oanda.com/rest-live-v20/account-ep/
*/
type AccountEndpoint struct {
	Account []AuthAcc `json:"accounts"`
}

/*
embedded struct for AccountEndpoint

endpoint: /v3/accounts
*/
type AuthAcc struct {
	ID   string   `json:"id"`
	Tags []string `json:"tags"`
}

/*
struct for unmarshalling json from [Account Endpoints] to get full details for an account.

endpoint: /v3/accounts/{accountID}

[Account Endpoints]: https://developer.oanda.com/rest-live-v20/account-ep/
*/
type AccountID struct {
	Account           IdDetails `json:"account"`
	LastTransactionID string    `json:"lastTransactionID"`
}

/*
embedded struct for AccountID

endpoint: /v3/accounts/{accountID}
*/
type IdDetails struct {
	GuaranteedStopLossOrderMode string      `json:"guaranteedStopLossOrderMode"`
	HedgingEnabled              bool        `json:"hedgingEnabled"`
	ID                          string      `json:"id"`
	CreatedTime                 string      `json:"createdTime"`
	Currency                    string      `json:"currency"`
	CreatedByUserID             int         `json:"createdByUserID"`
	Alias                       string      `json:"alias"`
	MarginRate                  string      `json:"marginRate"`
	LastTransactionID           string      `json:"lastTransactionID"`
	Balance                     string      `json:"balance"`
	OpenTradeCount              int         `json:"openTradeCount"`
	OpenPositionCount           int         `json:"openPositionCount"`
	PendingOrderCount           int         `json:"pendingOrderCount"`
	PL                          string      `json:"pl"`
	ResettablePL                string      `json:"resettablePL"`
	ResettablePLTime            string      `json:"resettablePLTime"`
	Financing                   string      `json:"financing"`
	Commission                  string      `json:"commission"`
	DividendAdjustment          string      `json:"dividendAdjustment"`
	GuaranteedExecutionFees     string      `json:"guaranteedExecutionFees"`
	Orders                      []string    `json:"orders"`
	Positions                   []Positions `json:"positions"`
	Trades                      []string    `json:"trades"`
	UnrealizedPL                string      `json:"unrealizedPL"`
	NAV                         string      `json:"NAV"`
	MarginUsed                  string      `json:"marginUsed"`
	MarginAvailable             string      `json:"marginAvailable"`
	PositionValue               string      `json:"positionValue"`
	MarginCloseoutUnrealizedPL  string      `json:"marginCloseoutUnrealizedPL"`
	MarginCloseoutNAV           string      `json:"marginCloseoutNAV"`
	MarginCloseoutMarginUsed    string      `json:"marginCloseoutMarginUsed"`
	MarginCloseoutPositionValue string      `json:"marginCloseoutPositionValue"`
	MarginCloseoutPercent       string      `json:"marginCloseoutPercent"`
	WithdrawalLimit             string      `json:"withdrawalLimit"`
	MarginCallMarginUsed        string      `json:"marginCallMarginUsed"`
	MarginCallPercent           string      `json:"marginCallPercent"`
}

/*
embedded struct for AccountID

endpoint: /v3/accounts/{accountID}
*/
type Positions struct {
	Instrument string `json:"instrument"`
	Long       Long   `json:"long"`
	Short      Short  `json:"short"`
}

/*
embedded struct for AccountID

endpoint: /v3/accounts/{accountID}
*/
type Long struct {
	Instrument              string `json:"instrument"`
	Units                   string `json:"units"`
	PL                      string `json:"pl"`
	ResettablePL            string `json:"resettablePL"`
	Financing               string `json:"financing"`
	DividendAdjustment      string `json:"dividendAdjustment"`
	GuaranteedExecutionFees string `json:"guaranteedExecutionFees"`
	UnrealizedPL            string `json:"unrealizedPL"`
}

/*
embedded struct for AccountID

endpoint: /v3/accounts/{accountID}
*/
type Short struct {
	Instrument              string `json:"instrument"`
	Units                   string `json:"units"`
	PL                      string `json:"pl"`
	ResettablePL            string `json:"resettablePL"`
	Financing               string `json:"financing"`
	DividendAdjustment      string `json:"dividendAdjustment"`
	GuaranteedExecutionFees string `json:"guaranteedExecutionFees"`
	UnrealizedPL            string `json:"unrealizedPL"`
}

/*
struct for unmarshalling json from [Account Endpoints] to get a summary for a single account that a client has access to.

endpoint: /v3/accounts/{accountID}/summary

[Account Endpoints]: https://developer.oanda.com/rest-live-v20/account-ep/
*/
type AccountSummary struct {
	Account           SummaryDetails `json:"account"`
	LastTransactionID string         `json:"lastTransactionID"`
}

/*
embedded struct for AccountSummary

endpoint: /v3/accounts/{accountID}/summary
*/
type SummaryDetails struct {
	NAV                         string `json:"account"`
	Alias                       string `json:"alias"`
	Balance                     string `json:"balance"`
	CreatedByUserID             int    `json:"createdByUserID"`
	CreatedTime                 string `json:"createdTime"`
	Currency                    string `json:"currency"`
	HedgingEnabled              bool   `json:"hedgingEnabled"`
	ID                          string `json:"id"`
	LastTransactionID           string `json:"lastTransactionID"`
	MarginAvailable             string `json:"marginAvailable"`
	MarginCloseoutMarginUsed    string `json:"marginCloseoutMarginUsed"`
	MarginCloseoutNAV           string `json:"marginCloseoutNAV"`
	MarginCloseoutPercent       string `json:"marginCloseoutPercent"`
	MarginCloseoutPositionValue string `json:"marginCloseoutPositionValue"`
	MarginCloseoutUnrealizedPL  string `json:"marginCloseoutUnrealizedPL"`
	MarginRate                  string `json:"marginRate"`
	MarginUsed                  string `json:"marginUsed"`
	OpenPositionCount           int    `json:"openPositionCount"`
	OpenTradeCount              int    `json:"openTradeCount"`
	PendingOrderCount           int    `json:"pendingOrderCount"`
	PL                          string `json:"pl"`
	PositionValue               string `json:"positionValue"`
	ResettablePL                string `json:"resettablePL"`
	UnrealizedPL                string `json:"unrealizedPL"`
	WithdrawalLimit             string `json:"withdrawalLimit"`
}

/*
GetAccounts function will get a list of all accounts one is authorized to use with the provided token.

For more info go to Oandas documentation for [Account Endpoints].

[Account Endpoints]: https://developer.oanda.com/rest-live-v20/account-ep/
*/
func GetAccounts(token string) (*AccountEndpoint, error) {
	url := "https://api-fxpractice.oanda.com/v3/accounts"

	// declare http client request, set to timeout after 10 seconds
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// preapre http get request with url
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err.Error())
	}

	// set headers for get request
	// check Oandas Best Practices for guidance https://developer.oanda.com/rest-live-v20/best-practices/
	// and check their instrument endpoint for candles https://developer.oanda.com/rest-live-v20/instrument-ep/
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept-Datetime-Format", "RFC3339")
	// req.URL.RawQuery = req.URL.Query().Encode()
	req.URL.Query().Encode()

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err.Error())
	} else if response.StatusCode == 400 {
		return nil, fmt.Errorf("400 error: %d", response.StatusCode)
	}
	defer response.Body.Close()

	// response body is []byte
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var account AccountEndpoint
	err = json.Unmarshal(body, &account)

	if err != nil {
		return nil, fmt.Errorf("error unmarshaling json: %s", err.Error())
	}

	return &account, nil
}

/*
GetAccountID function will return full details given an accountID for which one is authorized to use with a valid token. This includes full pending order, open trade, and open position representations are provided.

For more info go to Oandas documentation for [Account Endpoints].

[Account Endpoints]: https://developer.oanda.com/rest-live-v20/account-ep/
*/
func GetAccountID(id string, token string) (*AccountID, error) {
	url := "https://api-fxpractice.oanda.com/v3/accounts/" + id

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept-Datetime-Format", "RFC3339")
	req.URL.Query().Encode()

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err.Error())
	} else if response.StatusCode == 400 {
		return nil, fmt.Errorf("400 error: %d", response.StatusCode)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var accountid AccountID
	err = json.Unmarshal(body, &accountid)

	if err != nil {
		return nil, fmt.Errorf("error unmarshaling json: %s", err.Error())
	}

	return &accountid, nil
}

/*
GetAccountID function will return full details given an accountID for which one is authorized to use with a valid token. This includes full pending order, open trade, and open position representations are provided.

For more info go to Oandas documentation for [Account Endpoints].

[Account Endpoints]: https://developer.oanda.com/rest-live-v20/account-ep/
*/
func GetAccountSummary(id string, token string) (*AccountSummary, error) {
	url := "https://api-fxpractice.oanda.com/v3/accounts/" + id + "/summary"

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept-Datetime-Format", "RFC3339")
	req.URL.Query().Encode()

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err.Error())
	} else if response.StatusCode == 400 {
		return nil, fmt.Errorf("400 error: %d", response.StatusCode)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var accountsummary AccountSummary
	err = json.Unmarshal(body, &accountsummary)

	if err != nil {
		return nil, fmt.Errorf("error unmarshaling json: %s", err.Error())
	}

	return &accountsummary, nil
}
