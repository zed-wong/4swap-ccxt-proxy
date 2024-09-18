package main

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	fswap "github.com/fox-one/4swap-sdk-go/v2"
)

// Parameters:
// payAssetId: The ID of the asset to be paid
// fillAssetId: The ID of the asset to be received
// payAmount: The amount of the asset to be paid
// followID: An optional unique identifier for tracking the order
func FswapPreOrder(ctx context.Context, payAssetId, fillAssetId, payAmount, followID string) (string, error) {
	client := ctx.Value("client").(*fswap.Client)
	pairs, err := client.ListPairs(ctx)
	if err != nil {
		return "", err
	}
	payAmountt, err := decimal.NewFromString(payAmount)
	if err != nil {
			return "", err
	}
	PreOrderReq := &fswap.PreOrderReq{
		PayAssetID:  payAssetId,
		FillAssetID: fillAssetId,
		PayAmount:   payAmountt,
	}
	preOrder, err := fswap.PreOrderWithPairs(pairs, PreOrderReq)
	if err != nil {
		return "", err
	}
	if (followID == "") {
		followID = uuid.Must(uuid.NewV4()).String()
	}
	minAmount := preOrder.FillAmount.Mul(decimal.NewFromFloat(0.99)).Truncate(8)
	memo := fswap.BuildSwap(followID, fillAssetId, preOrder.Paths, minAmount)
	return memo, nil
}