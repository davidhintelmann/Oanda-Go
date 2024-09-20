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
struct for unmarshalling json from [Account Endpoints] for which one is
authorized with the provided token.

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
struct for unmarshalling json from [Account Endpoints] to get full details for
an account.

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
	GuaranteedStopLossOrderMode string        `json:"guaranteedStopLossOrderMode"`
	HedgingEnabled              bool          `json:"hedgingEnabled"`
	ID                          string        `json:"id"`
	CreatedTime                 string        `json:"createdTime"`
	Currency                    string        `json:"currency"`
	CreatedByUserID             int           `json:"createdByUserID"`
	Alias                       string        `json:"alias"`
	MarginRate                  string        `json:"marginRate"`
	LastTransactionID           string        `json:"lastTransactionID"`
	Balance                     string        `json:"balance"`
	OpenTradeCount              int           `json:"openTradeCount"`
	OpenPositionCount           int           `json:"openPositionCount"`
	PendingOrderCount           int           `json:"pendingOrderCount"`
	PL                          string        `json:"pl"`
	ResettablePL                string        `json:"resettablePL"`
	ResettablePLTime            string        `json:"resettablePLTime"`
	Financing                   string        `json:"financing"`
	Commission                  string        `json:"commission"`
	DividendAdjustment          string        `json:"dividendAdjustment"`
	GuaranteedExecutionFees     string        `json:"guaranteedExecutionFees"`
	Orders                      []string      `json:"orders"`
	Positions                   []PositionsID `json:"positions"`
	Trades                      []string      `json:"trades"`
	UnrealizedPL                string        `json:"unrealizedPL"`
	NAV                         string        `json:"NAV"`
	MarginUsed                  string        `json:"marginUsed"`
	MarginAvailable             string        `json:"marginAvailable"`
	PositionValue               string        `json:"positionValue"`
	MarginCloseoutUnrealizedPL  string        `json:"marginCloseoutUnrealizedPL"`
	MarginCloseoutNAV           string        `json:"marginCloseoutNAV"`
	MarginCloseoutMarginUsed    string        `json:"marginCloseoutMarginUsed"`
	MarginCloseoutPositionValue string        `json:"marginCloseoutPositionValue"`
	MarginCloseoutPercent       string        `json:"marginCloseoutPercent"`
	WithdrawalLimit             string        `json:"withdrawalLimit"`
	MarginCallMarginUsed        string        `json:"marginCallMarginUsed"`
	MarginCallPercent           string        `json:"marginCallPercent"`
}

/*
embedded struct for AccountID

endpoint: /v3/accounts/{accountID}
*/
type PositionsID struct {
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
struct for unmarshalling json from [Account Endpoints] to get a summary for
a single account that a client has access to.

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
struct for unmarshalling json from [Account Endpoints] to get a list of tradeable
instruments for the given account. The list of tradeable instruments is dependent
on the regulatory division that the account is located in, thus should be the same
for all accounts owned by a single user.

endpoint: /v3/accounts/{accountID}/instruments

[Account Endpoints]: https://developer.oanda.com/rest-live-v20/account-ep/
*/
type AccountInstru struct {
	List              []InstruDetails `json:"instruments"`
	LastTransactionID string          `json:"lastTransactionID"`
}

/*
embedded struct for AccountInstru

endpoint: /v3/accounts/{accountID}/instruments
*/
type InstruDetails struct {
	Name                        string       `json:"name"`
	Type                        string       `json:"type"`
	DisplayName                 string       `json:"displayName"`
	PipLocation                 int          `json:"pipLocation"`
	DisplayPrecision            int          `json:"displayPrecision"`
	TradeUnitsPrecision         int          `json:"tradeUnitsPrecision"`
	MinimumTradeSize            string       `json:"minimumTradeSize"`
	MaximumTrailingStopDistance string       `json:"maximumTrailingStopDistance"`
	MinimumTrailingStopDistance string       `json:"minimumTrailingStopDistance"`
	MaximumPositionSize         string       `json:"maximumPositionSize"`
	MaximumOrderUnits           string       `json:"maximumOrderUnits"`
	MarginRate                  string       `json:"marginRate"`
	GuaranteedStopLossOrderMode string       `json:"guaranteedStopLossOrderMode"`
	Tags                        []InstruTags `json:"tags"`
}

/*
embedded struct for AccountInstru

endpoint: /v3/accounts/{accountID}/instruments
*/
type InstruTags struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

/*
embedded struct for AccountInstru

endpoint: /v3/accounts/{accountID}/instruments
*/
type InstruFinancing struct {
	LongRate            string             `json:"longRate"`
	ShortRate           string             `json:"shortRate"`
	FinancingDaysOfWeek []InstruDaysOfWeek `json:"financingDaysOfWeek"`
}

/*
embedded struct for AccountInstru

endpoint: /v3/accounts/{accountID}/instruments
*/
type InstruDaysOfWeek struct {
	DayOfWeek   string `json:"dayOfWeek"`
	DaysCharged string `json:"daysCharged"`
}

/*
struct for unmarshalling json from [Account Endpoints] which can be used
to poll an Account for its current state and changes since a specified TransactionID.

endpoint: /v3/accounts/{accountID}/changes

[Account Endpoints]: https://developer.oanda.com/rest-live-v20/account-ep/
*/
type AccountChange struct {
	Changes           Changes `json:"changes"`
	State             State   `json:"state"`
	LastTransactionID string  `json:"lastTransactionID"`
}

/*
embedded struct for AccountChange

endpoint: /v3/accounts/{accountID}/changes
*/
type Changes struct {
	OrdersCancelled []string           `json:"ordersCancelled,omitempty"` // incomplete, wrong type
	OrdersCreated   []string           `json:"ordersCreated,omitempty"`   // incomplete, wrong type
	OrdersFilled    []OrdersFilled     `json:"ordersFilled,omitempty"`
	OrdersTriggered []string           `json:"ordersTriggered,omitempty"` // incomplete, wrong type
	Positions       []ChangesPositions `json:"positions,omitempty"`
	TradesClosed    []string           `json:"tradesClosed,omitempty"` // incomplete, wrong type
	TradesOpened    []struct {
		CurrentUnits string `json:"currentUnits"`
		Financing    string `json:"financing"`
		ID           string `json:"id"`
		InitialUnits string `json:"initialUnits"`
		Instrument   string `json:"instrument"`
		OpenTime     string `json:"openTime"`
		Price        string `json:"price"`
		RealizedPL   string `json:"realizedPL"`
		State        string `json:"state"`
	} `json:"tradesOpened,omitempty"`
	TradesReduced []string `json:"tradesReduced,omitempty"` // incomplete, wrong type
	Transactions  []struct {
		AccountBalance string `json:"accountBalance,omitempty"`
		AccountID      string `json:"accountID,omitempty"`
		BatchID        string `json:"batchID,omitempty"`
		Financing      string `json:"financing,omitempty"`
		ID             string `json:"id,omitempty"`
		Instrument     string `json:"instrument,omitempty"`
		PositionFill   string `json:"positionFill,omitempty"`
		Reason         string `json:"reason,omitempty"`
		Time           string `json:"time,omitempty"`
		TimeInForce    string `json:"timeInForce,omitempty"`
		TradeOpened    []struct {
			TradeID string `json:"tradeID"`
			Units   string `json:"units"`
		} `json:"tradeOpened,omitempty"`
		Type   string `json:"type"`
		Units  string `json:"units"`
		UserID string `json:"userID"`
	} `json:"transactions,omitempty"`
}

/*
embedded struct for AccountChange

endpoint: /v3/accounts/{accountID}/changes
*/
type ChangesPositions struct {
	Instrument string `json:"instrument"`
	Long       struct {
		AveragePrice string   `json:"averagePrice,omitempty"`
		PL           string   `json:"pl"`
		ResettablePL string   `json:"resettablePL"`
		TradeIDs     []string `json:"tradeIDs,omitempty"`
		Units        string   `json:"units"`
	} `json:"long,omitempty"`
	PL           string `json:"pl"`
	ResettablePL string `json:"resettablePL"`
	Short        struct {
		AveragePrice string   `json:"averagePrice,omitempty"`
		PL           string   `json:"pl"`
		ResettablePL string   `json:"resettablePL"`
		TradeIDs     []string `json:"tradeIDs,omitempty"`
		Units        string   `json:"units"`
	} `json:"short,omitempty"`
}

/*
embedded struct for AccountChange

endpoint: /v3/accounts/{accountID}/changes
*/
type OrdersFilled struct {
	CreateTime           string `json:"createTime"`
	FilledTime           string `json:"filledTime"`
	FillingTransactionID string `json:"fillingTransactionID"`
	ID                   string `json:"id"`
	Instrument           string `json:"instrument"`
	PositionFill         string `json:"positionFill"`
	State                string `json:"state"`
	TimeInForce          string `json:"timeInForce"`
	TradeOpenedID        string `json:"tradeOpenedID"`
	Type                 string `json:"type"`
	Units                string `json:"units"`
}

/*
embedded struct for AccountChange

endpoint: /v3/accounts/{accountID}/changes
*/
type State struct {
	NAV                        string           `json:"NAV"`
	MarginAvailable            string           `json:"marginAvailable"`
	MarginCloseoutMarginUsed   string           `json:"marginCloseoutMarginUsed"`
	MarginCloseoutNAV          string           `json:"marginCloseoutNAV"`
	MarginCloseoutPercent      string           `json:"marginCloseoutPercent"`
	MarginCloseoutUnrealizedPL string           `json:"marginCloseoutUnrealizedPL"`
	MarginUsed                 string           `json:"marginUsed"`
	Orders                     []string         `json:"orders,omitempty"` // incomplete, wrong type
	PositionValue              string           `json:"positionValue"`
	Positions                  []StatePositions `json:"positions,omitempty"` // incomplete, wrong type
	Trades                     []StateTrades    `json:"trades,omitempty"`    // incomplete, wrong type
	UnrealizedPL               string           `json:"unrealizedPL"`
	WithdrawalLimit            string           `json:"withdrawalLimit"`
}

/*
embedded struct for AccountChange

endpoint: /v3/accounts/{accountID}/changes
*/
type StatePositions struct {
	Instrument        string `json:"instrument"`
	LongUnrealizedPL  string `json:"longUnrealizedPL"`
	NetUnrealizedPL   string `json:"netUnrealizedPL"`
	ShortUnrealizedPL string `json:"shortUnrealizedPL"`
}

/*
embedded struct for AccountChange

endpoint: /v3/accounts/{accountID}/changes
*/
type StateTrades struct {
	ID           string `json:"id"`
	UnrealizedPL string `json:"unrealizedPL"`
}

/*
struct for unmarshalling error messages returned from Oanda's REST-V20 API
*/
type ErrorMsg struct {
	ErrorMessage string `json:"errorMessage"`
}

func (m *ErrorMsg) Error() string {
	return m.ErrorMessage
}

func (m *ErrorMsg) Empty() bool {
	if len(m.ErrorMessage) > 0 {
		return false
	} else {
		return true
	}
}

/*
GetAccounts function will get a list of all accounts one is authorized to use with
the provided token.

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
GetAccountID function will return full details given an accountID for which one is
authorized to use with a valid token. This includes full pending order, open trade,
and open position representations are provided.

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
GetAccountID function will return Get a summary for a single account that a
client has access to.

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

/*
GetAccountID function will return a list of tradeable instruments for the given account.
The list of tradeable instruments is dependent on the regulatory division that the account
is located in, thus should be the same for all accounts owned by a single user.

For more info go to Oandas documentation for [Account Endpoints].

[Account Endpoints]: https://developer.oanda.com/rest-live-v20/account-ep/
*/
func GetAccountInstru(id string, token string) (*AccountInstru, error) {
	url := "https://api-fxpractice.oanda.com/v3/accounts/" + id + "/instruments"

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

	var accountInstru AccountInstru
	err = json.Unmarshal(body, &accountInstru)

	if err != nil {
		return nil, fmt.Errorf("error unmarshaling json: %s", err.Error())
	}

	return &accountInstru, nil
}

/*
GetAccountID function will return an account for its current state and changes since a specified transactionID.

For more info go to Oandas documentation for [Account Endpoints].

[Account Endpoints]: https://developer.oanda.com/rest-live-v20/account-ep/
*/
func GetAccountChanges(id string, transactionID string, token string) (*AccountChange, error) {
	url := "https://api-fxpractice.oanda.com/v3/accounts/" + id + "/changes"

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
	// query parameters for get request
	q := req.URL.Query()
	q.Add("sinceTransactionID", transactionID)
	req.URL.RawQuery = q.Encode()

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err.Error())
	} else if response.StatusCode == 400 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatalln(err)
		}

		var errorMsg ErrorMsg
		err = json.Unmarshal(body, &errorMsg)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling json: %s", err.Error())
		}

		return nil, &ErrorMsg{ErrorMessage: errorMsg.ErrorMessage}
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var accountChange AccountChange
	err = json.Unmarshal(body, &accountChange)

	if err != nil {
		return nil, fmt.Errorf("error unmarshaling json: %s", err.Error())
	}

	return &accountChange, nil
}
