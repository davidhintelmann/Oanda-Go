package restful

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Account struct {
	Account struct {
		ID    string `json:"id"`
		Token string `json:"token"`
	} `json:"primary"`
}

type Metadata struct {
	Instrument  string `json:"instrument"`
	Granularity string `json:"granularity"`
	Candles     OHLC   `json:"candles"`
}

type OHLC []struct {
	Complete bool   `json:"complete"`
	Volume   int    `json:"volume"`
	Time     string `json:"time"`
	Bid      struct {
		O string `json:"o"`
		H string `json:"h"`
		L string `json:"l"`
		C string `json:"c"`
	} `json:"bid"`
	Ask struct {
		O string `json:"o"`
		H string `json:"h"`
		L string `json:"l"`
		C string `json:"c"`
	} `json:"ask"`
}

// GetIdToken function will return id & token
// first you must enter your ID and Token into
// res.json file, one can get these at
// https://fxtrade.oanda.com/your_account/fxtrade/register/gate?utm_source=oandaapi&utm_medium=link&utm_campaign=devportaldocs_demo
func GetIdToken(file_path string, display bool) (*Account, error) {
	jsonFile, err := os.Open(file_path)

	if err != nil {
		log.Fatal("error opening json file: ", err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Account variable
	var account Account

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'account' which we defined above
	err = json.Unmarshal(byteValue, &account)

	if err != nil {
		log.Fatal("error unmarshaling json: ", err)
	}
	// Print the account ID and Token to the console
	// if display parameter is true
	if display {
		fmt.Printf("ID: %s\n", account.Account.ID)
		fmt.Printf("Token: %s\n", account.Account.Token)
	}

	return &account, err
}

// Get Request for Instrument endpoint - Returns OHLC Bid/Ask
// https://developer.oanda.com/rest-live-v20/instrument-ep/
func GetCandlesBA(instrument, granularity, token string, display bool) (*Metadata, error) {
	url := "https://api-fxpractice.oanda.com/v3/instruments/" + instrument + "/candles"

	// declare http client request, set to timeout after 10 seconds
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// preapre http get request with url
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("Got error %s", err.Error())
	}

	// set headers for get request
	// check Oandas Best Practices for guidance https://developer.oanda.com/rest-live-v20/best-practices/
	// and check their instrument endpoint for candles https://developer.oanda.com/rest-live-v20/instrument-ep/
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept-Datetime-Format", "RFC3339")

	// query parameters for get request
	q := req.URL.Query()
	q.Add("granularity", granularity)
	q.Add("price", "BA")
	// encore the url and print it
	req.URL.RawQuery = q.Encode()

	// print string to console for debugging
	// fmt.Println(req.URL.String())

	response, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("Got error %s", err.Error())
	}
	defer response.Body.Close()

	// response body is []byte
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatalln(err)
	}

	// unmarshal the json data from get response
	var candles Metadata
	err = json.Unmarshal(body, &candles)

	if err != nil {
		return nil, fmt.Errorf("error unmarshaling json: %s", err.Error())
	}

	// if display parameter is true for GetCandlesBA() func
	// then print the get response
	if display {
		candle_count := len(candles.Candles)
		most_recent_candle := &candles.Candles[candle_count-1]
		fmt.Println()
		fmt.Printf("Instrument: \t\t%s\n", candles.Instrument)
		fmt.Printf("Granularity: \t\t%s\n", candles.Granularity)
		fmt.Printf("Candles - Count: \t%v\n", candle_count)
		fmt.Printf("Candles - Complete: \t%t\n", most_recent_candle.Complete)
		fmt.Printf("Candles - Volume: \t%v\n", most_recent_candle.Volume)
		fmt.Printf("Candles - Time: \t%s\n", most_recent_candle.Time)
		fmt.Println("\t- Bid:")
		fmt.Printf("\t\tOpen: \t%s\n", most_recent_candle.Bid.O)
		fmt.Printf("\t\tHigh: \t%s\n", most_recent_candle.Bid.H)
		fmt.Printf("\t\tLow: \t%s\n", most_recent_candle.Bid.L)
		fmt.Printf("\t\tClose: \t%s\n", most_recent_candle.Bid.C)
		fmt.Println("\t- Ask:")
		fmt.Printf("\t\tOpen: \t%s\n", most_recent_candle.Ask.O)
		fmt.Printf("\t\tHigh: \t%s\n", most_recent_candle.Ask.H)
		fmt.Printf("\t\tLow: \t%s\n", most_recent_candle.Ask.L)
		fmt.Printf("\t\tClose: \t%s\n", most_recent_candle.Ask.C)
	}

	return &candles, err
}
