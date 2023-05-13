package main

import (
	"context"
	"log"

	"github.com/davidhintelmann/Oanda-Go/connect"
	"github.com/davidhintelmann/Oanda-Go/restful"
)

// must include ID and Token into
// res.json file, one can get these at
// https://fxtrade.oanda.com/your_account/fxtrade/register/gate?utm_source=oandaapi&utm_medium=link&utm_campaign=devportaldocs_demo
const account_path string = "./res.json"

// user, password, database name for postgresql instance, and ssl mode
const user, dbname, sslmode = "david", "GinTest", "disable"

// be careful not to expose your password to the public
// modify password_edit.go file found in connect directory
// rename that file password.go and this file only has one function.
// This functions name needs to be edited by removing the underscore
// at the end of 'ImportPassword_' function.
// modify the return statement for the password of your PostgreSQL
// user that is accessing your local instance
var password = connect.ImportPassword()

// use background context globally to pass between functions
var ctx_psql = context.Background()

func main() {
	// set log flags for date and script file with line number for where the error occurred
	log.SetFlags(log.Ldate | log.Lshortfile)

	// Get ID and Token for Oanda Account
	idToken, err := restful.GetIdToken(account_path, false)
	id, token := idToken.Account.ID, idToken.Account.Token
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

	// connect to local instance of PostgreSQL
	conn, err := connect.ConnectPSQL(ctx_psql, user, password, dbname, sslmode)
	if err != nil {
		log.Fatalf("error during ConnectPSQL(): %v", err)
	}
	defer conn.Close()
	err = conn.Ping(ctx_psql)
	if err != nil {
		log.Fatalf("error pinning the local instance of PostgreSQL: %v\n", err)
	}

	instrumentList := "USD_CAD,USD_JPY,USD_CHF,USD_HKD,USD_SGD,GBP_USD,NZD_USD,EUR_CAD,EUR_USD,EUR_GBP,EUR_AUD,EUR_JPY,AUD_CAD,AUD_USD,AUD_NZD,AUD_JPY,AUD_HKD,CAD_HKD,CAD_CHF,CAD_JPY,CAD_SGD"
	connect.GetStreamPSQL(ctx_psql, conn, password, instrumentList, token, id, false)
}
