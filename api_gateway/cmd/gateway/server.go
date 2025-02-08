package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	authpb "github.com/jasmine-nguyen/go-microservices/api_gateway/auth"
	mmpb "github.com/jasmine-nguyen/go-microservices/api_gateway/money_movement"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var authClient authpb.AuthServiceClient
var mmClient mmpb.MoneyMovementServiceClient

func main() {
	authConn, err := grpc.NewClient("auth:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := authConn.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	authClient = authpb.NewAuthServiceClient(authConn)

	mmConn, err := grpc.NewClient("money_movement:7000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := mmConn.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	mmClient = mmpb.NewMoneyMovementServiceClient(mmConn)
	log.Printf("---money movement client created: %v", mmClient)
	http.HandleFunc("/login", login)
	http.HandleFunc("/customer/payment/authorize", customerPaymentAuthorize)
	http.HandleFunc("/customer/payment/capture", customerPaymentCapture)

	fmt.Println("listening on port: 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	userName, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	ctx := context.Background()
	token, err := authClient.GetToken(ctx, &authpb.Credentials{UserName: userName, Password: password})
	if err != nil {
		if err.Error() == "invalid credentials" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		_, writeErr := w.Write([]byte(err.Error()))
		if writeErr != nil {
			log.Println(writeErr)
		}
		return
	}

	_, err = w.Write([]byte(token.Jwt))
	if err != nil {
		log.Println(err)
		return
	}
}

func customerPaymentAuthorize(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	log.Printf("---auth header: %s", authHeader)
	if authHeader == "" {
		log.Println("---empty header")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		log.Println("---missing prefix")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	log.Printf("---token after trim: %s", token)
	ctx := context.Background()
	_, err := authClient.ValidateToken(ctx, &authpb.Token{
		Jwt: token,
	})
	if err != nil {
		log.Println("---invalid token")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	type authorizePayload struct {
		CustomerWalletUserId string `json:"customer_wallet_user_id"`
		MerchantWalletUserId string `json:"merchant_wallet_user_id"`
		Cents                int64  `json:"cents"`
		Currency             string `json:"currency"`
	}

	var payload authorizePayload
	body, err := io.ReadAll(r.Body)
	log.Printf("---body: %s", string(body))
	if err != nil {
		log.Printf("---error reading request body: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &payload)
	log.Printf("---payload: %v", payload)
	if err != nil {
		log.Printf("---error unmarshaling payload: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	ctx = context.Background()
	ar, err := mmClient.Authorize(ctx, &mmpb.AuthorizeRequest{
		CustomerWalletUserId: payload.CustomerWalletUserId,
		MerchantWalletUserId: payload.MerchantWalletUserId,
		Cents:                payload.Cents,
		Currency:             payload.Currency,
	})
	log.Printf("---authorize_payload: %v", ar)
	if err != nil {
		log.Printf("---authorize payload error: %s", err.Error())
		_, writeErr := w.Write([]byte(err.Error()))
		if writeErr != nil {
			log.Println(writeErr)
			return
		}
	}

	type response struct {
		Pid string `json:"pid"`
	}

	resp := &response{
		Pid: ar.Pid,
	}
	log.Printf("---resp: %v", resp)
	responseJSON, err := json.Marshal(resp)
	log.Printf("---responseJSON: %v", responseJSON)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		log.Println(err)
		return
	}
}

func customerPaymentCapture(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	ctx := context.Background()
	_, err := authClient.ValidateToken(ctx, &authpb.Token{
		Jwt: token,
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	type capturePayload struct {
		Pid string `json:"pid"`
	}

	var payload capturePayload
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	ctx = context.Background()
	_, err = mmClient.Capture(ctx, &mmpb.CaptureRequest{
		Pid: payload.Pid,
	})
	if err != nil {
		_, writeErr := w.Write([]byte(err.Error()))
		if writeErr != nil {
			log.Println(writeErr)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
