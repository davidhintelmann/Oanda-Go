// restful package is a wrapper API for [Oanda-V20] RESTful API.
// Currently this wrapper API only covers two endpoints:
//
//  1. Get request for [Instrument - candles endpoint] which returns
//     historical OHLC Bid/Ask.
//
//     - Parameters requires instrument symbol, token, and granularity (i.e., 'S5' for 5 second candles)
//
//  2. Get JSON Stream for [Pricing - stream endpoint]
//     which returns live Bid/Ask.
//
//     - Parameters requires list of instruments, token, and id
//
// Don't forget to check Oanda's [Best Practices] before querying any
// of their endpoints.
//
// [Oanda-V20]: https://developer.oanda.com/rest-live-v20/introduction/
// [Instrument - candles endpoint]: https://developer.oanda.com/rest-live-v20/instrument-ep/
// [Pricing - stream endpoint]: https://developer.oanda.com/rest-live-v20/pricing-ep/
// [Best Practices]: https://developer.oanda.com/rest-live-v20/best-practices/
package restful

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// struct for unmarshalling primary account in `res.json` file which
// contains account information (i.e., id and token).
type PrimaryAccount struct {
	Account Account `json:"primary"`
}

type Account struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

// struct for unmarshalling all accounts in `res.json` file which
// contains account information (i.e., id and token).
type Credential struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

type Credentials struct {
	Account map[string]Credential `json:"-"` // Dynamically handle all fields with the same structure
}

func (d *Credentials) UnmarshalJSON(data []byte) error {
	// Unmarshal everything into a map[string]Credential
	var rawMap map[string]Credential
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	// Store the dynamically parsed fields
	d.Account = rawMap
	return nil
}

// struct for unmarshalling metadata from Oanda's [Instrument - candles endpoint].
//
// [Instrument - candles endpoint]: https://developer.oanda.com/rest-live-v20/instrument-ep/
type Metadata struct {
	Instrument  string `json:"instrument"`
	Granularity string `json:"granularity"`
	Candles     []OHLC `json:"candles"`
}

// struct for unmarshalling FOREX OHLC data from Oanda's [Instrument - candles endpoint].
//
// [Instrument - candles endpoint]: https://developer.oanda.com/rest-live-v20/instrument-ep/
type OHLC struct {
	Complete bool   `json:"complete"`
	Volume   int    `json:"volume"`
	Time     string `json:"time"`
	Bid      Bid    `json:"bid"`
	Ask      Ask    `json:"ask"`
}

type Bid struct {
	O string `json:"o"`
	H string `json:"h"`
	L string `json:"l"`
	C string `json:"c"`
}

type Ask struct {
	O string `json:"o"`
	H string `json:"h"`
	L string `json:"l"`
	C string `json:"c"`
}

/*
FormatTime function will format the time, as specified by input, by parsing a OHLC time string into a go lang time.Time type and then return time in string format.
*/
func (ohlc *OHLC) FormatTime(format string) string {
	timestamp, err := time.Parse(time.RFC3339, ohlc.Time)
	if err != nil {
		log.Fatalf("error parsing timestamp: %v", err)
	}
	return timestamp.Format(format)
}

// struct for unmarshalling json data from Oanda's [Pricing - stream endpoint].
//
// [Pricing - stream endpoint]: https://developer.oanda.com/rest-live-v20/pricing-ep/
type Stream struct {
	Type string `json:"type"`
	Time string `json:"time"`
	Bids [1]struct {
		Price     string `json:"price"`
		Liquidity int64  `json:"liquidity"`
	} `json:"bids"`
	Asks [1]struct {
		Price     string `json:"price"`
		Liquidity int64  `json:"liquidity"`
	} `json:"asks"`
	CloseOutBid string `json:"closeoutbid"`
	CloseOutAsk string `json:"closeoutask"`
	Status      string `json:"status"`
	Tradeable   bool   `json:"tradeable"`
	Instrument  string `json:"instrument"`
}

// struct for unmarshalling json data from Oanda's [Pricing - stream endpoint].
//
// Note: every 5 seconds a 'heartbeat' is sent from this endpoint to
// let you know the connection is still alive.
//
// [Pricing - stream endpoint]: https://developer.oanda.com/rest-live-v20/pricing-ep/
type HeartBeat struct {
	Type string `json:"type"`
	Time string `json:"time"`
}

// GetIdToken function will return id & token for primary account.
// First you must enter your ID and token into
// res.json file. You can generate these from
// Oanda's [Demo Account].
//
// [Demo Account]: https://fxtrade.oanda.com/your_account/fxtrade/register/gate?utm_source=oandaapi&utm_medium=link&utm_campaign=devportaldocs_demo
func GetIdToken(file_path string, display bool) (*PrimaryAccount, error) {
	log.SetFlags(log.Ldate | log.Lshortfile)
	jsonFile, err := os.Open(file_path)

	// change log.Fatal to log.Print for resful_test.go to work
	// specifically the TestGetIdTokenInvalidPath() test function
	if err != nil {
		// log.Print("error opening json file: ", err)
		return nil, err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Print("error during ioutil.ReadAll(jsonFile): ", err)
		return nil, err
	}

	// we initialize our Account variable
	var account PrimaryAccount

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'account' which we defined above
	err = json.Unmarshal(byteValue, &account)

	// change log.Fatal to log.Print for resful_test.go to work
	// specifically the TestGetIdTokenInvalidPath() test function
	if err != nil {
		// log.Print("error unmarshaling json: ", err)
		return nil, err
	}
	// Print the account ID and Token to the console
	// if display parameter is true
	if display {
		fmt.Printf("ID: %s\n", account.Account.ID)
		fmt.Printf("Token: %s\n", account.Account.Token)
	}

	return &account, err
}

// GetAllIdToken function will return all id & token pairs.
// First you must enter your ID and token into
// res.json file. You can generate these from
// Oanda's [Demo Account].
//
// [Demo Account]: https://fxtrade.oanda.com/your_account/fxtrade/register/gate?utm_source=oandaapi&utm_medium=link&utm_campaign=devportaldocs_demo
func GetAllIdToken(file_path string, display bool) (*Credentials, error) {
	log.SetFlags(log.Ldate | log.Lshortfile)
	jsonFile, err := os.Open(file_path)

	// change log.Fatal to log.Print for resful_test.go to work
	// specifically the TestGetIdTokenInvalidPath() test function
	if err != nil {
		log.Print("error opening json file: ", err)
		return nil, err
	}
	defer jsonFile.Close()

	var credentials Credentials
	decoder := json.NewDecoder(jsonFile)
	if err := decoder.Decode(&credentials); err != nil {
		log.Fatal(err)
	}

	// Output the dynamically captured fields
	if display {
		for key, credential := range credentials.Account {
			fmt.Printf("Account: %s, ID: %s, Token: %s\n", key, credential.ID, credential.Token)
		}
	}

	return &credentials, err
}

// Get Request for Instrument endpoint - returns historical OHLC Bid/Ask.
//   - Parameters requires instrument symbol, token, and granularity (i.e., 'S5' for 5 second candles)
//
// See [Instrument - candles endpoint]
//
// [Instrument - candles endpoint]: https://developer.oanda.com/rest-live-v20/instrument-ep/
func GetCandlesBA(instrument, granularity, token string, display bool) (*Metadata, error) {
	url := "https://api-fxpractice.oanda.com/v3/instruments/" + instrument + "/candles"

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

	// query parameters for get request
	q := req.URL.Query()
	q.Add("granularity", granularity)
	q.Add("price", "BA")
	// encore the url
	req.URL.RawQuery = q.Encode()

	// print string to console for debugging
	// fmt.Println(req.URL.String())

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
		mostRecentCandle := &candles.Candles[candle_count-1]
		fmt.Printf("Instrument: \t\t%s\n", candles.Instrument)
		fmt.Printf("Granularity: \t\t%s\n", candles.Granularity)
		fmt.Printf("Candles - Count: \t%v\n", candle_count)
		fmt.Printf("Candles - Complete: \t%t\n", mostRecentCandle.Complete)
		fmt.Printf("Candles - Volume: \t%v\n", mostRecentCandle.Volume)
		fmt.Printf("Candles - Time: \t%s\n", mostRecentCandle.Time)
		fmt.Println("\t- Bid:")
		fmt.Printf("\t\tOpen: \t%s\n", mostRecentCandle.Bid.O)
		fmt.Printf("\t\tHigh: \t%s\n", mostRecentCandle.Bid.H)
		fmt.Printf("\t\tLow: \t%s\n", mostRecentCandle.Bid.L)
		fmt.Printf("\t\tClose: \t%s\n", mostRecentCandle.Bid.C)
		fmt.Println("\t- Ask:")
		fmt.Printf("\t\tOpen: \t%s\n", mostRecentCandle.Ask.O)
		fmt.Printf("\t\tHigh: \t%s\n", mostRecentCandle.Ask.H)
		fmt.Printf("\t\tLow: \t%s\n", mostRecentCandle.Ask.L)
		fmt.Printf("\t\tClose: \t%s\n", mostRecentCandle.Ask.C)
	}

	return &candles, err
}
