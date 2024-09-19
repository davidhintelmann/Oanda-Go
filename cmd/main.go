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

// Oanda-Go package is a wrapper API for [Oanda-V20] RESTful API.
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
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/davidhintelmann/Oanda-Go/oanda"
)

// must include ID and Token into
// res.json file, one can get these at
// https://fxtrade.oanda.com/your_account/fxtrade/register/gate?utm_source=oandaapi&utm_medium=link&utm_campaign=devportaldocs_demo
const accountJSON string = "../res.json"

func main() {
	// set flag for account name in res.json file
	accountName := flag.String("account", "primary", "Account name found in res.json file")
	flag.Parse()

	// set log flags for date and script file with line number for where the error occurred
	log.SetFlags(log.Ldate | log.Llongfile)

	// Get ID and Token for Oanda Account
	creds, err := oanda.GetAllIdToken(accountJSON, false)
	if err != nil {
		log.Fatalf("error during GetIdToken(): %v", err)
	}
	token := creds.Account[*accountName].Token

	// GetCandlesBA function sends a GET request to Oanda's API
	// set the display parameter to true to output OHLC data to the console
	_, err = oanda.GetCandlesBA("USD_CAD", "S5", token, false)
	if err != nil {
		log.Fatalf("error during GetCandlesBA(): %v", err)
	}

	accounts, err := oanda.GetAccounts(token)
	if err != nil {
		log.Fatalf("error during GetAccounts(): %v", err)
	} else {
		fmt.Println(token)
		fmt.Println(accounts.Account[0].ID)
	}

	accountID, err := oanda.GetAccountID(accounts.Account[0].ID, token)
	if err != nil {
		log.Fatalf("error during GetAccountID(): %v", err)
	} else {
		fmt.Println(accountID.Account.Balance)
	}

	summary, err := oanda.GetAccountSummary(accounts.Account[0].ID, token)
	if err != nil {
		log.Fatalf("error during GetAccountSummary(): %v", err)
	} else {
		fmt.Println(summary.Account.Currency)
	}

	instruments, err := oanda.GetAccountInstru(accounts.Account[0].ID, token)
	if err != nil {
		log.Fatalf("error during GetAccountInstru(): %v", err)
	} else {
		for _, instru := range instruments.List {
			fmt.Println(instru.Name)
		}
		fmt.Printf("There are %d instruments to choose from.\n", len(instruments.List))
	}
}
