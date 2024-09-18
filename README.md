# 4swap ccxt proxy

This is an proxy API server for 4swap implementation of ccxt, it would be almost impossible to maintain the code if we add everything to ccxt. So this is a separated proxy server using the existing SDKs to handle the 4swap memo generation, the initiation of mixin safe transaction and etc.

## Usage

Modify `config.json` to set the correct host and port, then run:
```
go mod tidy
go build
./4swap-ccxt-proxy
```

## API

### POST /4swap/preorder

#### Request

Headers:
- Authorization: The JWT token requested from mixin api by sign `/me`.

Parameters:
- payAssetId: The ID of the asset to be paid.
- fillAssetId: The ID of the asset to be received.
- payAmount: The amount of the asset to be paid.
- followID: An optional unique identifier for tracking the transaction.

#### Response

- memo: The memo to be used for the transaction.
- followID: The followID used for tracking the transaction.


#### Test using curl

```
curl -i -X GET -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" http://127.0.0.1:8080/4swap/preorder?payAssetId="4d8c508b-91c5-375b-92b0-ee702ed2dac5"&fillAssetId="31d2ea9c-95eb-3355-b65b-ba096853bc18"&payAmount="0.01"
```

### POST /mixin/transfer

#### Request

Parameters:
- asset_id: The ID of the asset to be transferred.
- amount: The amount of the asset to be transferred.
- trace_id: The trace ID of the transaction.

#### Response
