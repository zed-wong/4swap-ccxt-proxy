package main

import (
	"fmt"
	"strings"
	"context"
	"strconv"
	"net/http"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/fox-one/mixin-sdk-go/v2"
	fswap "github.com/fox-one/4swap-sdk-go/v2"
	"github.com/fox-one/mixin-sdk-go/v2/mixinnet"
)

// POST /4swap/preorder
// Parameters:
// payAssetId: The ID of the asset to be paid
// fillAssetId: The ID of the asset to be received
// payAmount: The amount of the asset to be paid
// followID: An optional unique identifier for tracking the order
//
// note: The JWT token must be signed with /me
func fswapPreOrderHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		user, err := mixin.UserMe(ctx, token)
		if err != nil || user.UserID == "" {
			http.Error(w, "Invalid access token", http.StatusInternalServerError)
			return
		}

		client, ok := ctx.Value("client").(*fswap.Client)
		if !ok || client == nil {
			http.Error(w, "Client not found in context", http.StatusInternalServerError)
			return
		}
		client.UseToken(token)
		var requestData struct {
			PayAssetID  string `json:"payAssetId"`
			FillAssetID string `json:"fillAssetId"`
			PayAmount   string `json:"payAmount"`
			FollowID    string `json:"followID"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestData); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		if requestData.FollowID == "" {
			requestData.FollowID = uuid.Must(uuid.NewV4()).String()
		}

		if requestData.PayAssetID == "" {
			http.Error(w, "Missing required parameter: payAssetId", http.StatusBadRequest)
			return
		}
		if requestData.FillAssetID == "" {
			http.Error(w, "Missing required parameter: fillAssetId", http.StatusBadRequest)
			return
		}
		if requestData.PayAmount == "" {
			http.Error(w, "Missing required parameter: payAmount", http.StatusBadRequest)
			return
		}

		payAssetId := requestData.PayAssetID
		fillAssetId := requestData.FillAssetID
		payAmount := requestData.PayAmount
		followID := requestData.FollowID
		memo, err := FswapPreOrder(ctx, payAssetId, fillAssetId, payAmount, followID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// code := mixin.GenerateCode()
		response := map[string]string{
			"memo": memo,
			"follow_id": followID,
			"client_id": user.UserID,
			// "code": "CAN BE ADDED IN FUTURE",
			// "code_url": "CAN BE ADDED IN FUTURE",
		}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}

// POST /mixin/encodetx
// Parameters:
// tx: the string of the transaction
// sigs: the signatures map
//
// Response:
// raw: the raw of the transaction
func mixinEncodeHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tx := r.FormValue("tx")
		sigss := r.FormValue("sigs")
		var sigs []map[uint16]*mixinnet.Signature
		if sigss == "" {
			sigs = []map[uint16]*mixinnet.Signature{}
		}
		if tx == "" {
			http.Error(w, "Missing required parameter: tx", http.StatusBadRequest)
			return
		}
		raw, err := EncodeSafeTx(tx, sigs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response := map[string]interface{}{
			"raw": raw,
		}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}

/*
	curl -X POST http://127.0.0.1:8080/mixin/mixaddress \
     -H "Content-Type: application/json" \
     -d '{
           "members": ["a753e0eb-3010-4c4a-a7b2-a7bda4063f62","099627f8-4031-42e3-a846-006ee598c56e","aefbfd62-727d-4424-89db-ae41f75d2e04","d68ca71f-0e2c-458a-bb9c-1d6c2eed2497","e4bc0740-f8fe-418c-ae1b-32d9926f5863"],
           "threshold": "3"
         }'
*/
func mixinMixAddressHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestBody struct {
			Members   []string `json:"members"`
			Threshold string   `json:"threshold"`
		}

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&requestBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if len(requestBody.Members) == 0 {
			http.Error(w, "Missing required parameter: members", http.StatusBadRequest)
			return
		}
		if requestBody.Threshold == "" {
			http.Error(w, "Missing required parameter: threshold", http.StatusBadRequest)
			return
		}

		thresholdUint8, err := strconv.ParseUint(requestBody.Threshold, 10, 8)
		if err != nil {
			http.Error(w, "Invalid threshold value", http.StatusBadRequest)
			return
		}

		address := mixin.RequireNewMixAddress(requestBody.Members, uint8(thresholdUint8)).String()
		response := map[string]interface{}{
			"address": address,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

func rootAlive(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"state": "alive",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}


func StartAPIServer(ctx context.Context) {
	group := ctx.Value(MTG_GROUP)
	if group == nil {
		panic("4swap MTG group not found in context")
	}
	http.HandleFunc("/4swap/preorder", fswapPreOrderHandler(ctx))
	http.HandleFunc("/mixin/encodetx", mixinEncodeHandler(ctx))
	http.HandleFunc("/mixin/mixaddress", mixinMixAddressHandler(ctx))
	http.HandleFunc("/", rootAlive)
	
	host := ctx.Value(HOST_KEY).(string)
	port := ctx.Value(PORT_KEY).(int)
	fmt.Printf("\n\033[1;34mStarting API server on \033[1;32m%s:%d\033[0m\n", host, port)
	fmt.Printf("\033[1;33m[GET] \033[1;36m/\033[0m - Endpoint to check health\n")
	fmt.Printf("\033[1;33m[POST] \033[1;36m/4swap/preorder\033[0m - Endpoint to create a preorder for 4swap transactions (sign /me auth required)\n")
	fmt.Printf("\033[1;33m[POST] \033[1;36m/mixin/encodetx\033[0m - Endpoint to encode a Mixin transaction\n")
	fmt.Printf("\033[1;33m[POST] \033[1;36m/mixin/mixaddress\033[0m - Endpoint to create a Mixin mix address\n")
	address := fmt.Sprintf("%s:%d", host, port)
	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}
