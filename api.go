package main

import (
	"fmt"
	"strings"
	"context"
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
		payAssetId := r.FormValue("payAssetId")
		fillAssetId := r.FormValue("fillAssetId")
		payAmount := r.FormValue("payAmount")
		followID := r.FormValue("followID")
		if (followID == "") {
			followID = uuid.Must(uuid.NewV4()).String()
		}

		if payAssetId == "" {
			http.Error(w, "Missing required parameter: payAssetId", http.StatusBadRequest)
			return
		}
		if fillAssetId == "" {
			http.Error(w, "Missing required parameter: fillAssetId", http.StatusBadRequest)
			return
		}
		if payAmount == "" {
			http.Error(w, "Missing required parameter: payAmount", http.StatusBadRequest)
			return
		}
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

func StartAPIServer(ctx context.Context) {
	group := ctx.Value(MTG_GROUP)
	if group == nil {
		panic("4swap MTG group not found in context")
	}
	http.HandleFunc("/4swap/preorder", fswapPreOrderHandler(ctx))
	http.HandleFunc("/mixin/encodetx", mixinEncodeHandler(ctx))
	
	host := ctx.Value(HOST_KEY).(string)
	port := ctx.Value(PORT_KEY).(int)
	fmt.Printf("\n\033[1;34mStarting API server on \033[1;32m%s:%d\033[0m\n", host, port)
	fmt.Printf("\033[1;33m[POST] \033[1;36m/4swap/preorder\033[0m - Endpoint to create a preorder for 4swap transactions (auth required)\n")
	fmt.Printf("\033[1;33m[POST] \033[1;36m/mixin/encodetx\033[0m - Endpoint to encode a Mixin transaction\n")
	address := fmt.Sprintf("%s:%d", host, port)
	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}
