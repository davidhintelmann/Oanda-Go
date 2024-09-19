If you are looking for the original repo please checkout `og` branch.

# Query Oanda's REST-V20 API with Go

This repo contains go code for querying Oanda's API for Forex price info. One could also set up a demo account to trade.

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

1. Requires [Go](https://go.dev/dl/) 1.22.5 or greater
2. Register [demo account](https://fxtrade.oanda.com/your_account/fxtrade/register/gate?utm_source=oandaapi&utm_medium=link&utm_campaign=devportaldocs_demo) from Oanda to obtain an API key
3. Modify the `res_edit.json` file in this repo's root directory with `ID` and `Token` obtained in the second step
   - rename `res_edit.json` to `res.json` for go code to work correctly

Once the above is satisfied you can get the functions in this repo with:

    go get github.com/davidhintelmann/Oanda-Go/oanda

Then import in your `main.go` (or any go file) with 

    import "github.com/davidhintelmann/Oanda-Go/oanda"

## Endpoints

The following list are the endpoints one can reach using this package.

### Account

[Account Endpoints](https://developer.oanda.com/rest-live-v20/account-ep/) for Oanda's REST V-20 API.

#### GET
- [x] `accounts` Get a list of Orders for an Account
- [x] `accountID` Get the full details for a single Account that a client has access to. Full pending Order, open Trade and open Position representations are provided.
- [x] `summary` Get a summary for a single Account that a client has access to.
- [x] `instruments` Get the list of tradeable instruments for the given Account. The list of tradeable instruments is dependent on the regulatory division that the Account is located in, thus should be the same for all Accounts owned by a single user.
- [ ] `changes` Endpoint used to poll an Account for its current state and changes since a specified TransactionID.
#### PATCH
- [ ] `configuration` Set the client-configurable portions of an Account.

### Instrument

[Instrument Endpoints](https://developer.oanda.com/rest-live-v20/instrument-ep/) for Oanda's REST V-20 API.

#### GET
- [ ] `candles` Fetch candlestick data for an instrument.
- [ ] `orderBook` Fetch an order book for an instrument.
- [ ] `positionBook` Fetch a position book for an instrument.

### Order

[Order Endpoints](https://developer.oanda.com/rest-live-v20/order-ep/) for Oanda's REST V-20 API.

#### GET
- [ ] `orders` Get a list of Orders for an Account
- [ ] `pendingOrders` List all pending Orders in an Account
- [ ] `orders/{orderSpecifier}` Get details for a single Order in an Account
#### POST
- [ ] `orders` Create an Order for an Account
#### PUT
- [ ] `orders/{orderSpecifier}` Replace an Order in an Account by simultaneously cancelling it and creating a replacement Order
- [ ] `orders/{orderSpecifier}/cancel` Cancel a pending Order in an Account
- [ ] `orders/{orderSpecifier}/clientExtensions` Update the Client Extensions for an Order in an Account. Do not set, modify, or delete clientExtensions if your account is associated with MT4.

### Trade

[Trade Endpoints](https://developer.oanda.com/rest-live-v20/trade-ep/) for Oanda's REST V-20 API.

#### GET
- [ ] `trades` Get a list of Trades for an Account
- [ ] `openTrades` Get the list of open Trades for an Account
- [ ] `trades/{tradeSpecifier}` Get the details of a specific Trade in an Account
#### PUT
- [ ] `trades/{tradeSpecifier}/close` Close (partially or fully) a specific open Trade in an Account
- [ ] `trades/{tradeSpecifier}/clientExtensions` Update the Client Extensions for a Trade. Do not add, update, or delete the Client Extensions if your account is associated with MT4.
- [ ] `trades/{tradeSpecifier}/orders` Create, replace and cancel a Trade’s dependent Orders (Take Profit, Stop Loss and Trailing Stop Loss) through the Trade itself

### Position

[Position Endpoints](https://developer.oanda.com/rest-live-v20/position-ep/) for Oanda's REST V-20 API.

#### GET
- [ ] `positions` List all Positions for an Account. The Positions returned are for every instrument that has had a position during the lifetime of an the Account.
- [ ] `openPositions` List all open Positions for an Account. An open Position is a Position in an Account that currently has a Trade opened for it.
- [ ] `positions/{instrument}` Get the details of a single Instrument’s Position in an Account. The Position may by open or not.
#### PUT
- [ ] `positions/{instrument}/close` Closeout the open Position for a specific instrument in an Account.

### Transaction

[Transaction Endpoints](https://developer.oanda.com/rest-live-v20/transaction-ep/) for Oanda's REST V-20 API.

#### GET
- [ ] `transactions` Get a list of Transactions pages that satisfy a time-based Transaction query.
- [ ] `transactions/{transactionID}` Get the details of a single Account Transaction.
- [ ] `transactions/idrange` Get a range of Transactions for an Account based on the Transaction IDs.
- [ ] `transactions/sinceid` Get a range of Transactions for an Account starting at (but not including) a provided Transaction ID.
- [ ] `transactions/stream` Get a stream of Transactions for an Account starting from when the request is made. **Note:** This endpoint is served by the streaming URLs.

### Pricing

[Pricing Endpoints](https://developer.oanda.com/rest-live-v20/pricing-ep/) for Oanda's REST V-20 API.

#### GET
- [ ] `candles/latest` Get dancing bears and most recently completed candles within an Account for specified combinations of instrument, granularity, and price component.
- [ ] `pricing` Get pricing information for a specified list of Instruments within an Account.
- [ ] `pricing/stream` Get a stream of Account Prices starting from when the request is made.
This pricing stream does not include every single price created for the Account, but instead will provide at most 4 prices per second (every 250 milliseconds) for each instrument being requested.
If more than one price is created for an instrument during the 250 millisecond window, only the price in effect at the end of the window is sent. This means that during periods of rapid price movement, subscribers to this stream will not be sent every price.
Pricing windows for different connections to the price stream are not all aligned in the same way (i.e. they are not all aligned to the top of the second). This means that during periods of rapid price movement, different subscribers may observe different prices depending on their alignment. **Note:** This endpoint is served by the streaming URLs.
- [x] `instruments/{instrument}/candles` Fetch candlestick data for an instrument.