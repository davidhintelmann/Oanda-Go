package main

import (
	"fmt"
	"log"

	"github.com/davidhintelmann/Oanda-Go/restful"
)

// must include ID and Token into
// res.json file, one can get these at
// https://fxtrade.oanda.com/your_account/fxtrade/register/gate?utm_source=oandaapi&utm_medium=link&utm_campaign=devportaldocs_demo
const account_json_path string = "./ress.json"

func main() {
	// set log flags for date and script file with line number for where the error occurred
	log.SetFlags(log.Ldate | log.Lshortfile)

	// Get ID and Token for Oanda Account
	id_token, err := restful.GetIdToken(account_json_path, false)
	_, token := id_token.Account.ID, id_token.Account.Token
	if err != nil {
		log.Fatal("error during GetIdToken(): ", err)
	}

	// GetCandlesBA function sends GET request to Oandas API
	metadata, err := restful.GetCandlesBA("USD_CAD", "S5", token, false)

	if err != nil {
		log.Fatal("error during GetCandlesBA(): ", err)
	}

	// print out json reponse
	most_recent_candle := metadata.Candles[0]
	fmt.Printf("Instrument: %s\nGranularity: %s\n", metadata.Instrument, metadata.Granularity)
	fmt.Printf("Candles: %v\n", metadata.Candles[0])
	fmt.Printf("Candles - Complete: \t\t%v\n", most_recent_candle.Complete)
	fmt.Printf("\t- Time: \t\t%v\n", most_recent_candle.Time)
	fmt.Println("\t\tBid:")
	fmt.Printf("\t\t\tOpen: \t%v\n", most_recent_candle.Bid.O)
	fmt.Printf("\t\t\tHigh: \t%v\n", most_recent_candle.Bid.H)
	fmt.Printf("\t\t\tLow: \t%v\n", most_recent_candle.Bid.L)
	fmt.Printf("\t\t\tClose: \t%v\n", most_recent_candle.Bid.C)
	fmt.Println("\t\tAsk:")
	fmt.Printf("\t\t\tOpen: \t%v\n", most_recent_candle.Ask.O)
	fmt.Printf("\t\t\tHigh: \t%v\n", most_recent_candle.Ask.H)
	fmt.Printf("\t\t\tLow: \t%v\n", most_recent_candle.Ask.L)
	fmt.Printf("\t\t\tClose: \t%v\n", most_recent_candle.Ask.C)
}
