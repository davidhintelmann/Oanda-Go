package restful

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

/*
struct for unmarshalling [Account Endpoints] which one is authorized for with the provided token.

[Account Endpoints]: https://developer.oanda.com/rest-live-v20/account-ep/
*/
type AccountEndpoint struct {
	Account []AuthAcc `json:"accounts"`
}

type AuthAcc struct {
	ID   string   `json:"id"`
	Tags []string `json:"tags"`
}

/*
GetAccounts function will get a list of all accounts authorized for with the provided token.

For more info go to Oandas documentation for [Account Endpoints].

[Account Endpoints]: https://developer.oanda.com/rest-live-v20/account-ep/
*/
func GetAccounts(token string) (*AccountEndpoint, error) {
	url := "https://api-fxpractice.oanda.com/v3/accounts"

	// declare http client request, set to timeout after 10 seconds
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// preapre http get request with url
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err.Error())
	}

	// set headers for get request
	// check Oandas Best Practices for guidance https://developer.oanda.com/rest-live-v20/best-practices/
	// and check their instrument endpoint for candles https://developer.oanda.com/rest-live-v20/instrument-ep/
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept-Datetime-Format", "RFC3339")

	// query parameters for get request
	q := req.URL.Query()
	// q.Add("granularity", granularity)
	// q.Add("price", "BA")
	// encore the url
	req.URL.RawQuery = q.Encode()

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err.Error())
	} else if response.StatusCode == 400 {
		return nil, fmt.Errorf("400 error: %d", response.StatusCode)
	}
	defer response.Body.Close()

	// response body is []byte
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var account AccountEndpoint
	err = json.Unmarshal(body, &account)

	if err != nil {
		return nil, fmt.Errorf("error unmarshaling json: %s", err.Error())
	}

	return &account, nil
}
