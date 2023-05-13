package restful

import (
	"context"
	"database/sql"
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

type HeartBeat struct {
	Type string `json:"type"`
	Time string `json:"time"`
}

// GetIdToken function will return id & token
// first you must enter your ID and Token into
// res.json file, one can get these at
// https://fxtrade.oanda.com/your_account/fxtrade/register/gate?utm_source=oandaapi&utm_medium=link&utm_campaign=devportaldocs_demo
func GetIdToken(file_path string, display bool) (*Account, error) {
	log.SetFlags(log.Ldate | log.Lshortfile)
	jsonFile, err := os.Open(file_path)

	// change log.Fatal to log.Print for resful_test.go to work
	// specifically the TestGetIdTokenInvalidPath() test function
	if err != nil {
		log.Print("error opening json file: ", err)
		return nil, err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Print("error during ioutil.ReadAll(jsonFile): ", err)
		return nil, err
	}

	// we initialize our Account variable
	var account Account

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'account' which we defined above
	err = json.Unmarshal(byteValue, &account)

	// change log.Fatal to log.Print for resful_test.go to work
	// specifically the TestGetIdTokenInvalidPath() test function
	if err != nil {
		log.Print("error unmarshaling json: ", err)
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
	// encore the url
	req.URL.RawQuery = q.Encode()

	// print string to console for debugging
	// fmt.Println(req.URL.String())

	// time duration of request
	queryStart := time.Now()
	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Got error %s\n", err.Error())
	} else if response.StatusCode == 400 {
		return nil, fmt.Errorf("400 error %s\n", err.Error())
	}
	queryEnd := time.Now()
	queryDuration := queryEnd.Sub(queryStart)
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
		mostRecentCandle := &candles.Candles[candle_count-1]
		fmt.Printf("Get Request Duration: %v\n", queryDuration)
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

// Get JSON Stream for Pricing endpoint
// Parameters requires list of instruments, token, and id
// Returns Bid/Ask
// https://developer.oanda.com/rest-live-v20/pricing-ep/
func GetStreamMSSQL(ctx context.Context, conn *sql.DB, instrument string, token string, id string, display bool) {
	streamUrl := fmt.Sprintf("https://stream-fxpractice.oanda.com/v3/accounts/%s/pricing/stream", id)

	// declare http client request
	// no timeout due to endpoint being a data stream
	// unless idle connection
	tr := &http.Transport{
		MaxConnsPerHost: 2,
		MaxIdleConns:    2,
		IdleConnTimeout: 30 * time.Second,
	}
	client := &http.Client{Transport: tr}

	// preapre http get request with url
	req, err := http.NewRequest("GET", streamUrl, nil)

	if err != nil {
		log.Fatalln(err)
	}

	// set headers for get request
	// check Oandas Best Practices for guidance https://developer.oanda.com/rest-live-v20/best-practices/
	// and check their pricing endpoint for streaming https://developer.oanda.com/rest-live-v20/pricing-ep/
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept-Datetime-Format", "RFC3339")
	req.Header.Add("Connection", "Keep-Alive")

	// query parameters for get request
	q := req.URL.Query()
	q.Add("instruments", instrument)
	// Flag that enables/disables the sending of a
	// pricing snapshot when initially connecting to the stream.
	// [default=True]
	q.Add("snapshot", "True")
	// Flag that enables the inclusion of the
	// homeConversions field in the returned response.
	// An entry will be returned for each currency
	// in the set of all base and quote currencies
	// present in the requested instruments list.
	// [default=False]
	q.Add("includeHomeConversions", "False")
	// encore the url and print it
	req.URL.RawQuery = q.Encode()

	// print string to console for debugging
	fmt.Println(req.URL.String())

	// Send the GET Request
	response, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	// response body is []byte
	//body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatalln(err)
	}

	// dec := json.NewDecoder(strings.NewReader(body))
	dec := json.NewDecoder(response.Body)
	//dec := json.NewDecoder(bytes.NewReader(body))

	for {
		var tick Stream
		if err := dec.Decode(&tick); err != nil {
			log.Fatal(err)
		}
		if tick.Type != "PRICE" && tick.Instrument != "" && display == true {
			fmt.Printf("Type: %s\n", tick.Type)
			fmt.Printf("Time: %s\n", tick.Time)
			fmt.Println("Bids:")
			fmt.Printf("\tPrice: %s\n", tick.Bids[0].Price)
			fmt.Printf("\tLiquidity: %d\n", tick.Bids[0].Liquidity)
			fmt.Println("Ask:")
			fmt.Printf("\tPrice: %s\n", tick.Asks[0].Price)
			fmt.Printf("\tLiquidity: %d\n", tick.Asks[0].Liquidity)
			fmt.Printf("Close Out Bid: %s\n", tick.CloseOutBid)
			fmt.Printf("Close Out Ask: %s\n", tick.CloseOutAsk)
			fmt.Printf("Status: %s\n", tick.Status)
			fmt.Printf("Tradeable: %t\n", tick.Tradeable)
			fmt.Printf("Instrument: %s\n", tick.Instrument)
			// err = json.Unmarshal(body, &candles)
		} else if tick.Type == "HEARTBEAT" && display == true {
			fmt.Printf("Type: %s, Time: %s\n", tick.Type, tick.Time)
		}

		if tick.Type == "PRICE" && tick.Instrument != "" && display == false {
			// Microsoft SQL does not allow boolean values
			// convert to 0 or 1 (bit type) instead
			// where 1=true, and 0=false
			var tradeable int
			if tick.Tradeable == true {
				tradeable = 1
			} else {
				tradeable = 0
			}
			// Insert into PostNames Table in the Gin-Test database
			query := `INSERT INTO [Oanda-Stream].[dbo].[Stream] VALUES ('%s', '%s', '%s', %d, '%s', %d, '%s', '%s', '%s', %d, '%s');`
			tsql := fmt.Sprintf(query, tick.Type, tick.Time, tick.Bids[0].Price,
				tick.Bids[0].Liquidity, tick.Asks[0].Price, tick.Asks[0].Liquidity,
				tick.CloseOutBid, tick.CloseOutAsk, tick.Status, tradeable, tick.Instrument)
			// Execute query
			_, err = conn.QueryContext(ctx, tsql)
		} else if tick.Type == "HEARTBEAT" && display == false {
			fmt.Printf("%s, Time: %s\n", tick.Type, tick.Time)
			query := `INSERT INTO [Oanda-Stream].[dbo].[Heartbeats] VALUES ('%s', '%s');`
			tsql := fmt.Sprintf(query, tick.Type, tick.Time)
			// Execute query
			_, err = conn.QueryContext(ctx, tsql)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
