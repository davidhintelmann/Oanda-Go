package oanda_test

import (
	"fmt"
	"log"

	"github.com/davidhintelmann/Oanda-Go/oanda"
)

// must include ID and Token into
// res.json file, one can get these at
// https://fxtrade.oanda.com/your_account/fxtrade/register/gate?utm_source=oandaapi&utm_medium=link&utm_campaign=devportaldocs_demo
const accountJSON string = "../res_edit.json"

func Example_account() {
	// Get ID and Token for Oanda Account
	creds, err := oanda.GetAllIdToken(accountJSON, false)
	_, token := creds.Account["primary"].ID, creds.Account["primary"].Token
	if err != nil {
		log.Fatalf("error during GetAllIdToken(): %v", err)
	}

	acc, err := oanda.GetAccounts(token)
	if err != nil {
		log.Fatalf("error during GetAccounts(): %v", err)
	}

	// print all the accounts one is authorized to use with the provided token
	fmt.Println(acc.Account)
}

func Example_accountID() {
	// Get ID and Token for Oanda Account
	creds, err := oanda.GetAllIdToken(accountJSON, false)
	id, token := creds.Account["primary"].ID, creds.Account["primary"].Token
	if err != nil {
		log.Fatalf("error during GetAllIdToken(): %v", err)
	}

	accID, err := oanda.GetAccountID(id, token)
	if err != nil {
		log.Fatalf("error during GetAccountID(): %v", err)
	}

	// print balance for account
	fmt.Println(accID.Account.Balance)
}

func Example_accountSummary() {
	// Get ID and Token for Oanda Account
	creds, err := oanda.GetAllIdToken(accountJSON, false)
	id, token := creds.Account["primary"].ID, creds.Account["primary"].Token
	if err != nil {
		log.Fatalf("error during GetAllIdToken(): %v", err)
	}

	summary, err := oanda.GetAccountSummary(id, token)
	if err != nil {
		log.Fatalf("error during GetAccountSummary(): %v", err)
	}

	// print the base currency for this account
	fmt.Println(summary.Account.Currency)
}

func Example_accountInstruments() {
	// Get ID and Token for Oanda Account
	creds, err := oanda.GetAllIdToken(accountJSON, false)
	id, token := creds.Account["primary"].ID, creds.Account["primary"].Token
	if err != nil {
		log.Fatalf("error during GetAllIdToken(): %v", err)
	}

	instruments, err := oanda.GetAccountInstru(id, token)
	if err != nil {
		log.Fatalf("error during GetAccountInstru(): %v", err)
	}

	// print one possible instrument that you can trade on Oanda
	fmt.Println(instruments.List[0].Name)
}
