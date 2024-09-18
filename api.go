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
)

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

		if ENABLE_USERME {
			user, err := mixin.UserMe(ctx, token)
			if err != nil || user.UserID == "" {
				http.Error(w, "Invalid access token", http.StatusInternalServerError)
				return
			}
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
			// "code": "",
			// "code_url": "",
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

// Parameters:
// assetId: The ID of the asset to be transferred
// amount: The amount of the asset to be transferred
// memo: The memo for the transaction (generated from /4swap/preorder)
//
// note: The JWT token must be signed with /safe/outputs
func mixinTransferHandler(ctx context.Context) http.HandlerFunc {
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
	
		assetId := r.FormValue("assetId")
		amount := r.FormValue("amount")
		memo := r.FormValue("memo")
		if assetId == "" {
			http.Error(w, "Missing required parameter: assetId", http.StatusBadRequest)
			return
		}
		if amount == "" {
			http.Error(w, "Missing required parameter: amount", http.StatusBadRequest)
			return
		}
		if memo == "" {
			http.Error(w, "Missing required parameter: memo", http.StatusBadRequest)
			return
		}
		
		request, err := MixinTransferInit(ctx, token, assetId, amount, memo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response := map[string]interface{}{
			"request": request,
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
	fmt.Printf("MTG Group loaded: %+v\n", group)
	http.HandleFunc("/4swap/preorder", fswapPreOrderHandler(ctx))
	http.HandleFunc("/mixin/transfer", mixinTransferHandler(ctx))

	host := ctx.Value(HOST_KEY).(string)
	port := ctx.Value(PORT_KEY).(int)
	address := fmt.Sprintf("%s:%d", host, port)
	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}
