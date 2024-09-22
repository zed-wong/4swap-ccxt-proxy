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
	// take jwt token from ccxt, sign /safe/outputs&asset_id=xxx&state=unspent

	// step 1: get unspent outputs
	utxos, err := GetUTXOWithToken(ctx, token, assetId)
	if err != nil {
		return nil, fmt.Errorf("failed to get with bearer token: %w", err)
	}
	
	b := NewSafeTransactionBuilder(utxos)

	// step 2: build transaction
	mtgGroup := ctx.Value(MTG_GROUP).(*fswap.Group)
	mixAddress, err := MixAddressFromString(mtgGroup.MixAddress)
	if err != nil {
		mixAddress = RequireNewMixAddress(mtgGroup.Members, mtgGroup.Threshold)
	}

	amountDec, err := decimal.NewFromString(amount)
	if err != nil {
		return nil, fmt.Errorf("failed to parse amount: %w", err)
	}

	tx, err := MakeTransaction(ctx, b, []*TransactionOutput{
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
	request, err := PostCreateSafeTransaction(ctx, token, []*mixin.SafeTransactionRequestInput{
		{
			RequestID: uuidHash([]byte(utxos[0].OutputID + ":SafeCreateTransactionRequest")),
			RawTransaction: raw,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create safe transaction request: %w", err)
	}
	fmt.Printf("request: %+v\n", request)

	// step 4: sign transaction	& step 5: submit transaction
	// handled by ccxt

	return nil, nil
}

func uuidHash(b []byte) string {
	h := md5.New()
	h.Write(b)
	sum := h.Sum(nil)
	sum[6] = (sum[6] & 0x0f) | 0x30
	sum[8] = (sum[8] & 0x3f) | 0x80
	return uuid.FromBytesOrNil(sum).String()
}






















































func BuildSafeTxBeforeGhostKey(utxos []*mixin.SafeUtxo) {
	b := NewSafeTransactionBuilder(utxos)

	mtgGroup := ctx.Value(MTG_GROUP).(*fswap.Group)
	mixAddress, err := MixAddressFromString(mtgGroup.MixAddress)
	if err != nil {
		mixAddress = RequireNewMixAddress(mtgGroup.Members, mtgGroup.Threshold)
	}

	amountDec, err := decimal.NewFromString(amount)
	if err != nil {
		return nil, fmt.Errorf("failed to parse amount: %w", err)
	}

	return b, []*TransactionOutput{
		{
			Address: mixAddress,
			Amount:  amountDec,
		},
	}
}

func MakeSafeTxWithGhostKey() {

}