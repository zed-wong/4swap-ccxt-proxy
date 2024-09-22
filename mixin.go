package main

import (
	"fmt"
	"github.com/fox-one/mixin-sdk-go/v2/mixinnet"
)

/* 
  Steps required for mixin safe transfer:
  - Step 1: get unspent outputs (ccxt)
  - Step 2: build transaction (tx-ccxt, raw-go)
  - Step 3: verify transaction (ccxt)
  - Step 4: sign transaction (ccxt)
  - Step 5: submit transaction (ccxt)
*/
// Get raw from tx
func EncodeSafeTx(tx string, sigs []map[uint16]*mixinnet.Signature) (string, error) {
	txData := []byte(tx)
	transaction, err := mixinnet.TransactionFromData(txData)
	if sigs != nil {
		transaction.Signatures = sigs
	}
	if err != nil {
		return "", fmt.Errorf("failed to parse transaction: %w", err)
	}
	
	encoder := mixinnet.NewEncoder()
	raw := encoder.EncodeTransaction(transaction)

	return string(raw), nil
}