package main

import (
	"fmt"
	"sort"
	"context"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/fox-one/mixin-sdk-go/v2/mixinnet"
)

// https://github.com/fox-one/mixin-sdk-go/blob/d1b04c07e694f0798eebb27524334506dea12d91/transaction_input.go

type TransactionBuilder struct {
	*mixinnet.TransactionInput
	addr *MixAddress
}

type TransactionOutput struct {
	Address *MixAddress     `json:"address,omitempty"`
	Amount  decimal.Decimal `json:"amount,omitempty"`
}

func NewSafeTransactionBuilder(utxos []*mixin.SafeUtxo) *TransactionBuilder {
	b := &TransactionBuilder{
		TransactionInput: &mixinnet.TransactionInput{
			TxVersion: mixinnet.TxVersion,
			Hint:      newUUID(),
			Inputs:    make([]*mixinnet.InputUTXO, len(utxos)),
		},
	}

	for i, utxo := range utxos {
		b.Inputs[i] = &mixinnet.InputUTXO{
			Input: mixinnet.Input{
				Hash:  &utxo.TransactionHash,
				Index: utxo.OutputIndex,
			},
			Asset:  utxo.KernelAssetID,
			Amount: utxo.Amount,
		}

		addr, err := NewMixAddress(utxo.Receivers, utxo.ReceiversThreshold)
		if err != nil {
			panic(err)
		}

		if i == 0 {
			b.addr = addr
		} else if b.addr.String() != addr.String() {
			panic("invalid utxos")
		}
	}

	return b
}


func MakeTransaction(ctx context.Context, b *TransactionBuilder, outputs []*TransactionOutput) (*mixinnet.Transaction, error) {
	remain := b.TotalInputAmount()
	for _, output := range outputs {
		remain = remain.Sub(output.Amount)
	}

	if remain.IsPositive() {
		outputs = append(outputs, &TransactionOutput{
			Address: b.addr,
			Amount:  remain,
		})
	}
	
	if err := AppendOutputsToInput(ctx, b, outputs); err != nil {
		fmt.Printf("AppendOutputsToInput err: %+v\n", err)
		return nil, err
	}

	tx, err := b.Build()
	if err != nil {	
		return nil, err
	}

	return tx, nil
}

func AppendOutputsToInput(ctx context.Context, b *TransactionBuilder, outputs []*TransactionOutput) error {
	var (
		ghostInputs  []*mixin.GhostInput
		ghostOutputs []*mixinnet.Output
	)

	for _, output := range outputs {
		txOutput := &mixinnet.Output{
			Type:   mixinnet.OutputTypeScript,
			Amount: mixinnet.IntegerFromDecimal(output.Amount),
			Script: mixinnet.NewThresholdScript(output.Address.Threshold),
		}

		index := uint8(len(b.Outputs))
		if len(output.Address.uuidMembers) > 0 {
			ghostInputs = append(ghostInputs, &mixin.GhostInput{
				Receivers: output.Address.Members(),
				Index:     index,
				Hint:      uuidHash([]byte(fmt.Sprintf("hint:%s;index:%d", b.Hint, index))),
			})

			ghostOutputs = append(ghostOutputs, txOutput)
		}

		b.Outputs = append(b.Outputs, txOutput)
	}

	if len(ghostInputs) > 0 {
		keys, err := createGhostKeys(ctx, ghostInputs, b.addr.Members())
		if err != nil {
			return err
		}

		for i, key := range keys {
			output := ghostOutputs[i]
			output.Keys = key.Keys
			output.Mask = key.Mask
		}
	}

	return nil
}

func createGhostKeys(ctx context.Context, inputs []*mixin.GhostInput, senders []string) ([]*mixin.GhostKeys, error) {
	// sort receivers
	for _, input := range inputs {
		sort.Strings(input.Receivers)
	}

	return PostCreateSafeGhostKeys(ctx, inputs, senders...)
}

func newUUID() string {
	return uuid.Must(uuid.NewV4()).String()
}