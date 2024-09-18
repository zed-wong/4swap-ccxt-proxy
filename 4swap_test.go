package main

import (
	"context"
	"testing"

	fswap "github.com/fox-one/4swap-sdk-go/v2"
	"github.com/shopspring/decimal"
)

// MockClient is a mock implementation of the fswap.Client interface
type MockClient struct {
	ListPairsFunc func(ctx context.Context) ([]*fswap.Pair, error)
}

func (m *MockClient) ListPairs(ctx context.Context) ([]*fswap.Pair, error) {
	return m.ListPairsFunc(ctx)
}

func TestFswapPreOrder(t *testing.T) {
	ctx := context.Background()
	mockClient := &MockClient{
		ListPairsFunc: func(ctx context.Context) ([]*fswap.Pair, error) {
			return []*fswap.Pair{
				// Add mock pairs as needed
			}, nil
		},
	}
	ctx = context.WithValue(ctx, "client", mockClient)

	payAssetId := "4d8c508b-91c5-375b-92b0-ee702ed2dac5"
	fillAssetId := "31d2ea9c-95eb-3355-b65b-ba096853bc18"
	followID := "609910ae-a2c0-4d42-8966-94071eb904ad"
	payAmount := "0.01"

	// Mock PreOrderWithPairs
	originalPreOrderWithPairs := fswap.PreOrderWithPairs
	defer func() { fswap.PreOrderWithPairs = originalPreOrderWithPairs }()
	fswap.PreOrderWithPairs = func(pairs []*fswap.Pair, req *fswap.PreOrderReq) (*fswap.PreOrder, error) {
		expectedReq := &fswap.PreOrderReq{
			PayAssetID:  payAssetId,
			FillAssetID: fillAssetId,
			PayAmount:   decimal.NewFromFloat(0.01),
		}
		if req.PayAssetID != expectedReq.PayAssetID || req.FillAssetID != expectedReq.FillAssetID || !req.PayAmount.Equal(expectedReq.PayAmount) {
			t.Errorf("unexpected PreOrderReq: got %v, want %v", req, expectedReq)
		}
		return &fswap.PreOrder{
			FillAmount: decimal.NewFromFloat(0.0099),
			Paths:      []string{"path1", "path2"},
		}, nil
	}

	// Test the function
	memo, err := FswapPreOrder(ctx, payAssetId, fillAssetId, payAmount, followID)
	if err != nil {
		t.Fatalf("FswapPreOrder returned an error: %v", err)
	}
	if memo == "" {
		t.Fatalf("FswapPreOrder returned an empty memo")
	}
}
