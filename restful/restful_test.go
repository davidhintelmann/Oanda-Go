package restful

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetIdTokenValidPath(t *testing.T) {
	// test 1 - check path is relative
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf(`File Path for GetIdToken(file_path, display) needs to be a relative path to this projects root directory pointing to res.json file.
Parameter 'file_path' has not been entered correctly, please check your syntax.`)
	}
	file_path := filepath.Join("../", currentDir)
	if filepath.IsAbs(file_path) {
		t.Fatalf(`Parameter 'file_path' has not been entered correctly: file Path for GetIdToken(file_path, display) should be relative.
It is currently set as an absolute path. Keep res.json in projects root directory.

Error: %v`, err)
	}

	// test 2 - check res.json path works
	//
	// loop through two types of test
	//	1. test display parameter
	//	2. test correct path (see below)
	//
	// main.go, in projects root directory has global const
	// const account_json_path string = "./res.json"
	// path for this test is "../res.json"
	// ie. that is the path that needs to work
	pathNeedsToWork := "../res.json"
	for _, b := range [2]bool{true, false} {
		account, err := GetIdToken(pathNeedsToWork, b)

		if err != nil {
			t.Fatalf(`Test for GetIdToken(file_path, display) produced an error: %v
			
Should have returned id and token of type struct but path for res.json is incorrectly set.`, err)
		}

		id, token := account.Account.ID, account.Account.Token
		type_id, type_token := reflect.TypeOf(id).Kind(), reflect.TypeOf(token).Kind()

		if type_id != reflect.String || type_token != reflect.String {
			t.Fatalf(`account struct (id and token) returned by GetIdToken should unpack into two variables
ie. 'id, token := account.Account.ID, account.Account.Token'
'id' is supposed to be of type 'int' but is of type: %v
'token' is supposed to be of type 'string' but is of type: %v`, type_id, type_token)
		}
	}
}

func TestGetIdTokenInvalidPath(t *testing.T) {
	// Suppress log output from GetIdToken() function
	// errors are logged from each function
	null, _ := os.Open(os.DevNull)
	defer log.SetOutput(null)
	// test 1 - check path is relative
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf(`File Path for GetIdToken(file_path, display) needs to be a relative path to this projects root directory pointing to res.json file.
Parameter 'file_path' has not been entered correctly, please check your syntax.`)
	}
	file_path := filepath.Join("../", currentDir)
	if filepath.IsAbs(file_path) {
		t.Fatalf(`Parameter 'file_path' has not been entered correctly: file Path for GetIdToken(file_path, display) should be relative.
It is currently set as an absolute path. Keep res.json in projects root directory.

Error: %v`, err)
	}

	// test 2 - check any path other than res.json does not work
	//
	// main.go, in projects root directory has global const
	// const account_json_path string = "./res.json"
	// check 5 paths that all need to fail
	// last element ("../LICENSE") in InvalidPaths works but unmarshaling json will fail
	invalidPaths := [6]string{"res.json", "../../res.json", "../r.json", "../res.txt", "../LICENSE"}
	for _, v := range invalidPaths {
		_, err := GetIdToken(v, false)

		if err == nil {
			t.Fatalf(`Test for GetIdToken(file_path, display) produced an error: %v
			
'file_path' parameter in GetIdToken(file_path, display) finds correct path to open and unmarshal json but should fail`, err)
		}
	}
}

// must include ID and Token into
// res.json file, one can get these at
// https://fxtrade.oanda.com/your_account/fxtrade/register/gate?utm_source=oandaapi&utm_medium=link&utm_campaign=devportaldocs_demo
const account_path string = "./res.json"

func Example_candles() {
	// set log flags for date and script file with line number for where the error occurred
	log.SetFlags(log.Ldate | log.Lshortfile)

	// Get ID and Token for Oanda Account
	idToken, err := GetIdToken(account_path, false)
	_, token := idToken.Account.ID, idToken.Account.Token
	if err != nil {
		log.Fatalf("error during GetIdToken(): %v", err)
	}
	// GetCandlesBA function sends a GET request to Oanda's API
	// set the display parameter to true to output OHLC data to the console
	_, err = GetCandlesBA("USD_CAD", "S5", token, true)
	// candles := _.Candles
	if err != nil {
		log.Fatalf("error during GetCandlesBA(): %v", err)
	}

}
