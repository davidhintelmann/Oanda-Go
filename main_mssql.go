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

package OandaGo

import (
	"context"
	"database/sql"
	"log"

	"github.com/davidhintelmann/Oanda-Go/restful"
	_ "github.com/jackc/pgx/v5/pgxpool"
)

// must include ID and Token into
// res.json file, one can get these at
// https://fxtrade.oanda.com/your_account/fxtrade/register/gate?utm_source=oandaapi&utm_medium=link&utm_campaign=devportaldocs_demo
const account_json_path string = "./res.json"

// server, database, driver configuration
var server, database, driver = "lpc:localhost", "Oanda-Stream", "mssql" // "sqlserver" or "mssql"

// trusted connection, and encryption configuraiton
var trusted_connection, encrypt = true, true

// db is global variable to pass between functions
var conn *sql.DB

// Use background context globally to pass between functions
var ctx_mssql = context.Background()

func main_mssql() {
	// set log flags for date and script file with line number for where the error occurred
	log.SetFlags(log.Ldate | log.Lshortfile)

	// Get ID and Token for Oanda Account
	idToken, err := restful.GetIdToken(account_json_path, false)
	id, token := idToken.Account.ID, idToken.Account.Token
	if err != nil {
		log.Fatalf("error during GetIdToken(): %v", err)
	}

	// GetCandlesBA function sends a GET request to Oanda's API
	// set the display parameter to true to output OHLC data to the console
	_, err = restful.GetCandlesBA("USD_CAD", "S5", token, true)
	if err != nil {
		log.Fatalf("error during GetCandlesBA(): %v", err)
	}

	// connect to local instance of Microsoft SQL
	conn, err := restful.ConnectMSSQL(ctx_mssql, conn, driver, server, database, trusted_connection, encrypt)

	if err != nil {
		log.Fatalf("error during ConnectMSSQL(): %v", err)
	}
	defer conn.Close()

	instrumentList := "USD_CAD,USD_JPY,USD_CHF,USD_HKD,USD_SGD,GBP_USD,NZD_USD,EUR_CAD,EUR_USD,EUR_GBP,EUR_AUD,EUR_JPY,AUD_CAD,AUD_USD,AUD_NZD,AUD_JPY,AUD_HKD,CAD_HKD,CAD_CHF,CAD_JPY,CAD_SGD"
	restful.GetStreamMSSQL(ctx_mssql, conn, instrumentList, token, id, false)
}
