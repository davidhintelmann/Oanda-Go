// This package is for connecting to a local instance
// of PostgreSQL while acquiring FOREX data from Oanda.
package connect

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

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

// Connect to local instance of PostgreSQL.
//
// Requires username, password, database name, and if ssl mode should use encryption
// or not ("enable" or "disable")
func ConnectPSQL(ctx context.Context, user string, password string, dbname string, sslmode string) (*pgxpool.Pool, error) {
	// fmt.Print("Connecting to postgresql...\n")
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", user, password, dbname, sslmode)
	dbpool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		log.Printf("error during intial connection: %v\n", err)
		return nil, err
	}
	// need to remove line below to prevent
	// error occuring in main func
	// defer dbpool.Close()
	return dbpool, nil
}

// Get JSON Stream for Pricing endpoint - returns live Bid/Ask.
//   - Parameters requires list of instruments, token, and id
//
// See [Pricing - stream endpoint]
//
// [Pricing - stream endpoint]: https://developer.oanda.com/rest-live-v20/pricing-ep/
func GetStreamPSQL(ctx context.Context, conn *pgxpool.Pool, password string, instrument string, token string, id string, display bool) {
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
	// fmt.Println(req.URL.String())

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

	dec := json.NewDecoder(response.Body)

	for {
		var tick Stream
		if err := dec.Decode(&tick); err != nil {
			log.Fatal(err)
		}
		fmt.Println(tick.Instrument)
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
		} else if tick.Type == "PRICE" && tick.Instrument != "" && display == false {
			// Insert into stream.live Table in the GinTest database
			query := `INSERT INTO stream.live VALUES ('%s', '%s', '%s', %d, '%s', %d, '%s', '%s', '%s', %t, '%s');`
			tsql := fmt.Sprintf(query, tick.Type, tick.Time, tick.Bids[0].Price,
				tick.Bids[0].Liquidity, tick.Asks[0].Price, tick.Asks[0].Liquidity,
				tick.CloseOutBid, tick.CloseOutAsk, tick.Status, tick.Tradeable, tick.Instrument)
			// Execute query
			_ = conn.QueryRow(ctx, tsql)
		} else if tick.Type == "HEARTBEAT" && display == false {
			// Insert into stream.heartbeats Table in the GinTest database
			fmt.Printf("Type: %s, Time: %s\n", tick.Type, tick.Time)
			query := `INSERT INTO stream.heartbeats VALUES ('%s', '%s');`
			tsql := fmt.Sprintf(query, tick.Type, tick.Time)
			// Execute query
			_ = conn.QueryRow(ctx, tsql)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
