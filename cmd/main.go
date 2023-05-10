package main

import (
	"log"
	"net/http"

	handler "github.com/ssergomol/Chat-Bot/pkg/handler"
)

func main() {
	http.HandleFunc("/", handler.HandleTelegramWebHook)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
