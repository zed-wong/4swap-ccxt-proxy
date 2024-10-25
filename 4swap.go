package main

import (
	"fmt"
	"time"
	"context"
	"strconv"
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

func FswapAddLiquidity(followID, oppositeAsset, slippage, expireDuration string) (string, error) {
	if followID == "" {
		followID = uuid.Must(uuid.NewV4()).String()
	}
	slippageDecimal, err := decimal.NewFromString(slippage)
	if err != nil {
		return "", err
	}
	expireDurationInt, err := strconv.Atoi(expireDuration)
	if err != nil {
		return "", err
	}
	expireDurationX := time.Second * time.Duration(expireDurationInt)
	memo := fswap.BuildAdd(followID, oppositeAsset, slippageDecimal, expireDurationX)
	return memo, nil
}

func FswapRemoveLiquidity(followID string) (string, error) {
	if followID == "" {
		followID = uuid.Must(uuid.NewV4()).String()
	}
	memo := fswap.BuildRemove(followID)
	return memo, nil
}

// I was planned to write a function to combine /cmc/pairs and /pairs
// Because /cmc/pairs doesn't return amount in the pool 
// which will be needed for calculating the amount needed for add_liquidity
// But turns out /pairs/{base_id}/{quote_id} is implemented in sdk but not written in api doc
// So this function is unnessary now, but can be used for simplify the create add_liqudiity order process
// Right now we need to call /pairs/{base_id}/{quote_id} to get the amount in the pool
// If we have this function, we can get the amount in the pool by calling fetchMarket method in ccxt
func FswapCombinedPairs(ctx context.Context) (string, error) {
	client := ctx.Value("client").(*fswap.Client)
	combinedPairs, err := client.ListPairs(ctx)
	fmt.Println(combinedPairs)
	if err != nil {
		return "", err
	}
	return "", nil
}
