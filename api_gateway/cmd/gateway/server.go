package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/customer/payment/authorize", customerPaymentAuthorize)
	http.HandleFunc("/customer/payment/capture", customerPaymentCapture)
}

func login(w http.ResponseWriter, r *http.Request) {}
func customerPaymentAuthorize(w http.ResponseWriter, r *http.Request) {
}
func customerPaymentCapture(w http.ResponseWriter, r *http.Request) {}
