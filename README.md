# 4swap ccxt proxy

This is an proxy API server for 4swap implementation of ccxt, it would be almost impossible to maintain the code if we add everything to ccxt. So this is a separated proxy server using the existing SDKs to handle the 4swap memo generation, the initiation of mixin safe transaction and etc.

## Usage

Modify `HOST` and `PORT` in `main.go` to set the correct host and port, then run:
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
- client_id: The client ID of the bot.

#### Test using curl

```
curl -i -X GET -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" http://127.0.0.1:8080/4swap/preorder?payAssetId="4d8c508b-91c5-375b-92b0-ee702ed2dac5"&fillAssetId="31d2ea9c-95eb-3355-b65b-ba096853bc18"&payAmount="0.01"
```

### POST /mixin/encodetx

#### Request

Parameters:
- tx: The string of the transaction.
- sigs: The signatures of the transaction. (empty by default)

#### Response

- raw: The raw of the transaction.

## Example

$ mixin-cli -f keystore.json sign /me --exp 24h
sign GET /me with request id aaa3fbf0-b639-40b9-a18b-a5672a9dc4c9 & exp 24h0m0s

eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjY3NjMxNzMsImlhdCI6MTcyNjY3Njc3MywianRpIjoiYWFhM2ZiZjAtYjYzOS00MGI5LWExOGItYTU2NzJhOWRjNGM5Iiwic2NwIjoiRlVMTCIsInNpZCI6IjFjZDdiMjQ3LTIwOTItNDRhNS1iNTUzLTM5ZDY5NmJkN2I2YSIsInNpZyI6IjVlNmI1OGZmYTEwYjNiYzUxNzI0ZmYwYmJkMmFmYjkxYzQ3NzFlZTM0MGY1ZDY4NTM0MGRmYTRjODU0YmFmYmEiLCJ1aWQiOiI1MTE4NmQ3ZS1kNDg4LTQxN2QtYTAzMS1iNGUzNGY0ZmRmODYifQ.eQezhHjzI3lPi93NJhxGaF7BZ1G1o6eH5kHjXL_owecBch2SUJFnXUuVcNpnkinD4nH-Ym3Qguw8tBdfFd6rAQ%                                     

$ export TOKEN=eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjY3NjMxNzMsImlhdCI6MTcyNjY3Njc3MywianRpIjoiYWFhM2ZiZjAtYjYzOS00MGI5LWExOGItYTU2NzJhOWRjNGM5Iiwic2NwIjoiRlVMTCIsInNpZCI6IjFjZDdiMjQ3LTIwOTItNDRhNS1iNTUzLTM5ZDY5NmJkN2I2YSIsInNpZyI6IjVlNmI1OGZmYTEwYjNiYzUxNzI0ZmYwYmJkMmFmYjkxYzQ3NzFlZTM0MGY1ZDY4NTM0MGRmYTRjODU0YmFmYmEiLCJ1aWQiOiI1MTE4NmQ3ZS1kNDg4LTQxN2QtYTAzMS1iNGUzNGY0ZmRmODYifQ.eQezhHjzI3lPi93NJhxGaF7BZ1G1o6eH5kHjXL_owecBch2SUJFnXUuVcNpnkinD4nH-Ym3Qguw8tBdfFd6rAQ

$ curl -i -X GET -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" http://127.0.0.1:8080/4swap/preorder\?payAssetId\="4d8c508b-91c5-375b-92b0-ee702ed2dac5"\&fillAssetId\="31d2ea9c-95eb-3355-b65b-ba096853bc18"\&payAmount\="0.01"
HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 18 Sep 2024 16:26:32 GMT
Content-Length: 134

{"client_id":"51186d7e-d488-417d-a031-b4e34f4fdf86","follow_id":"e3356eec-cff7-4476-a18a-4f37e98683ce","memo":"AgEB4zVu7M/3RHahik836YaDzgADMdLqnJXrM1W2W7oJaFO8GAFkAQABAAAAAAAPOBZ9j8eF"}