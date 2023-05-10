package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ssergomol/Chat-Bot/pkg/handler"
)

func main() {
	log.Println("Starting telegram bot...")
	http.HandleFunc("/", handler.HandleTelegramWebHook)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
