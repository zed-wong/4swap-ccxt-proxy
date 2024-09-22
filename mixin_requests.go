package main

import (
	"fmt"
	"bytes"
	"context"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/fox-one/mixin-sdk-go/v2"
)

type SafeUtxoResponse struct {
	Data []*mixin.SafeUtxo `json:"data"`
}

func GetUTXOWithToken(ctx context.Context, token, assetId string) ([]*mixin.SafeUtxo, error) {
	var result []*mixin.SafeUtxo
	const LIMIT = 10000
	client := &http.Client{}
	
	req, err := http.NewRequest("GET", mixin.DefaultApiHost+"/safe/outputs", nil)
	if err != nil {
		return nil, err
	}

	// Add query parameters to the request
	q := req.URL.Query()
	q.Add("asset_id", assetId)
	q.Add("state", "unspent")
	q.Add("limit", fmt.Sprintf("%d", LIMIT))
	req.URL.RawQuery = q.Encode()

	// Set the Authorization header with the provided token
	req.Header.Set("Authorization", "Bearer " + token)

	// Execute the HTTP request
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	defer resp.Body.Close()

	var response SafeUtxoResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	result = response.Data
	return result, nil
}

func PostCreateSafeTransaction(ctx context.Context, token string, inputs []*mixin.SafeTransactionRequestInput) ([]*mixin.SafeTransactionRequest, error) {
	requestBody, err := json.Marshal(inputs)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal inputs: %w", err)
	}
	req, err := http.NewRequest("POST", mixin.DefaultApiHost+"/safe/transaction/requests", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + token)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", body)
	}

	// Unmarshal response
	var response []*mixin.SafeTransactionRequest
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, nil
}

func PostCreateSafeGhostKeys(ctx context.Context, token string, inputs []*mixin.GhostInput, senders ...string) ([]*mixin.GhostKeys, error) {
	var body interface{} = inputs

	if len(senders) > 0 {
		body = map[string]interface{}{
			"keys":    inputs,
			"senders": senders,
		}
	}

	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal inputs: %w", err)
	}

	req, err := http.NewRequest("POST", mixin.DefaultApiHost+"/safe/keys", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + token)

	client := &http.Client{}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", bodyBytes)
	}

	var response []*mixin.GhostKeys
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, nil
}