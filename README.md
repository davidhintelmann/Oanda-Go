# Query Oanda's REST-V20 API with Go

This repo contains go code for querying Oandas API for Forex price info. One could also set up a demo account to trade

***WARNING:*** This is for educational purposes only.

Taken from [Oanda's REST V-20 API](https://developer.oanda.com/rest-live-v20/introduction/) Introduction page:

    The Oanda's REST V-20 API provides programmatic access to Oandas's next generation v20 trading engine. To use this API you must have a v20 trading account, which is available to all divisions except OANDA Global Markets.

    What can I do with the OANDA REST-v20 API?

    - Get real-time rates for all tradeable pairs 24 hours a day
    - Access historical pricing information dating back to 2005
    - Place, modify, close orders
    - Manage your account settings
    - Access your account/trading history


## Installation

1. Requires [Go](https://go.dev/dl/) 1.20.3 or greater
2. Register [demo account](https://fxtrade.oanda.com/your_account/fxtrade/register/gate?utm_source=oandaapi&utm_medium=link&utm_campaign=devportaldocs_demo) from Oanda to obtain an API key
3. Modify the `res_edit.json` file in this repo's root directory with `ID` and `Token` obtained in the second step
   - rename `res_edit.json` to `res.json` for go code to work correctly

Once the above is satisfied you can get this repo with:

    $ go get github.com/davidhintelmann/Oanda-Go

and then install with:

    $ go install github.com/davidhintelmann/Oanda-Go