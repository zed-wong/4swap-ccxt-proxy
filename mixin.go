package main

import (
	"fmt"
	"context"
	"crypto/md5"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	"github.com/fox-one/mixin-sdk-go/v2"
	fswap "github.com/fox-one/4swap-sdk-go/v2"
)

/* 
  Steps required for mixin safe transfer:
  - Step 1: get unspent outputs
  - Step 2: build transaction
  - Step 3: verify transaction
  - Step 4: sign transaction
  - Step 5: submit transaction

	Step 1-3 will be handled by MixinTransferInit
	Step 4-5 will be implemented and handled by ccxt
*/

// Parameters:
// token: jwt token from ccxt
// assetId: asset id
// amount: amount to transfer
// memo: memo for the transaction
func MixinTransferInit(ctx context.Context, token, assetId, amount, memo string) (*mixin.SafeTransactionRequest, error) {
	// take jwt token from ccxt, sign /safe/outputs
	client := mixin.NewFromAccessToken(token)

	// step 0: get user info from jwt token
	me, err := client.UserMe(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user me: %w", err)
	}
	fmt.Printf("me: %+v\n", me)

	// step 1: get unspent outputs
	utxos, err := client.SafeListUtxos(ctx, mixin.SafeListUtxoOption{
		Members: []string{me.UserID},
		Limit:   1,
		State:   mixin.SafeUtxoStateUnspent,
	})
	b := mixin.NewSafeTransactionBuilder(utxos)

	// step 2: build transaction
	mtgGroup := ctx.Value(MTG_GROUP).(*fswap.Group)
	fmt.Printf("mtgGroup: %+v\n", mtgGroup)
	
	mixAddress, err := mixin.MixAddressFromString(mtgGroup.MixAddress)
	if err != nil {
		mixAddress = mixin.RequireNewMixAddress(mtgGroup.Members, mtgGroup.Threshold)
	}
	fmt.Printf("mixAddress: %+v\n", mixAddress)

	amountDec, err := decimal.NewFromString(amount)
	if err != nil {
		return nil, fmt.Errorf("failed to parse amount: %w", err)
	}
	tx, err := client.MakeTransaction(ctx, b, []*mixin.TransactionOutput{
		{
			Address: mixAddress,
			Amount:  amountDec,
		},
	})
	raw, err := tx.Dump()
	if err != nil {
		return nil, fmt.Errorf("failed to dump transaction: %w", err)
	}
	fmt.Println("raw:", raw)

	// step 3: verify transaction
	request, err := client.SafeCreateTransactionRequest(ctx, &mixin.SafeTransactionRequestInput{
		RequestID: uuidHash([]byte(utxos[0].OutputID + ":SafeCreateTransactionRequest")),
		RawTransaction: raw,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create safe transaction request: %w", err)
	}
	fmt.Printf("request: %+v\n", request)

	// step 4: sign transaction	& step 5: submit transaction
	// handled by ccxt

	return request, nil
}

func uuidHash(b []byte) string {
	h := md5.New()
	h.Write(b)
	sum := h.Sum(nil)
	sum[6] = (sum[6] & 0x0f) | 0x30
	sum[8] = (sum[8] & 0x3f) | 0x80
	return uuid.FromBytesOrNil(sum).String()
}