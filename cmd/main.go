// MIT License

// Copyright (c) 2023 David Hintelmann

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// OandaGo package is a wrapper API for [Oanda-V20] RESTful API.
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
// Additionally this package can be used with either Micrsoft SQL or
// PostgreSQL for inserting data and querying a local database instead
// of querying Oanda's API each time one needs historical data. You can also
// insert data from the live stream into a database.
//
// Note: `main_psql.go` is for using PostgreSQL and `main_mssql.go` is
// for using Microsoft SQL as your database.
//
// [Oanda-V20]: https://developer.oanda.com/rest-live-v20/introduction/
// [Instrument - candles endpoint]: https://developer.oanda.com/rest-live-v20/instrument-ep/
// [Pricing - stream endpoint]: https://developer.oanda.com/rest-live-v20/pricing-ep/
// [Best Practices]: https://developer.oanda.com/rest-live-v20/best-practices/
package main

import (
	"log"

	"github.com/davidhintelmann/Oanda-Go/restful"
)

// must include ID and Token into
// res.json file, one can get these at
// https://fxtrade.oanda.com/your_account/fxtrade/register/gate?utm_source=oandaapi&utm_medium=link&utm_campaign=devportaldocs_demo
const accountJSON string = "./res.json"

func main() {
	// set log flags for date and script file with line number for where the error occurred
	log.SetFlags(log.Ldate | log.Lshortfile)

	// Get ID and Token for Oanda Account
	idToken, err := restful.GetIdToken(accountJSON, false)
	_, token := idToken.Account.ID, idToken.Account.Token
	if err != nil {
		log.Fatalf("error during GetIdToken(): %v", err)
	}
	// GetCandlesBA function sends a GET request to Oanda's API
	// set the display parameter to true to output OHLC data to the console
	_, err = restful.GetCandlesBA("USD_CAD", "S5", token, true)
	// candles := _.Candles
	if err != nil {
		log.Fatalf("error during GetCandlesBA(): %v", err)
	}
}
