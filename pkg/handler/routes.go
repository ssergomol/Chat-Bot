package handler

import (
	"log"
	"net/http"
)

func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {
	log.Println("Got request\nTrying to parse...")
	var update, err = parseTelegramRequest(r)
	if err != nil {
		log.Printf("error parsing update, %s\n", err.Error())
		return
	}
	log.Println("Successfully parsed")

	if update.Message.Text == "/start" {
		if err = createButtons(update.Message.Chat.Id); err != nil {
			log.Printf("error creating buttons, %s\n", err.Error())
			return
		}
		log.Printf("Buttons successfully created\n")
	} else {

		outputMessage := "Your message: " + update.Message.Text

		// Send the punchline back to Telegram
		var telegramResponseBody, errTelegram = sendTextToTelegramChat(update.Message.Chat.Id, outputMessage)
		if errTelegram != nil {
			log.Printf("got error %s from telegram, reponse body is %s", errTelegram.Error(), telegramResponseBody)
		} else {
			log.Printf("message \"%s\" successfully sent to chat id %d\n", outputMessage, update.Message.Chat.Id)
		}
	}
}
