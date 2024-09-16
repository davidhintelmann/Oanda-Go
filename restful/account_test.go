package restful_test

import (
	"fmt"
	"log"

	"github.com/davidhintelmann/Oanda-Go/restful"
)

// must include ID and Token into
// res.json file, one can get these at
// https://fxtrade.oanda.com/your_account/fxtrade/register/gate?utm_source=oandaapi&utm_medium=link&utm_campaign=devportaldocs_demo
const accountJSON string = "../res_edit.json"

func Example_account() {
	// Get ID and Token for Oanda Account
	creds, err := restful.GetAllIdToken(accountJSON, false)
	_, token := creds.Account["primary"].ID, creds.Account["primary"].Token
	if err != nil {
		log.Fatalf("error during GetIdToken(): %v", err)
	}

	acc, err := restful.GetAccounts(token)
	if err != nil {
		log.Fatalf("error during GetCandlesBA(): %v", err)
	}

	fmt.Println(token)
	fmt.Println(acc.Account)
	// Output:
	// XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
	// []
}
