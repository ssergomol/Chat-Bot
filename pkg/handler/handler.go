package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

// Handles incoming update from the Telegram webhook
func parseTelegramRequest(r *http.Request) (*Update, error) {
	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode incoming update %s", err.Error())
		return nil, err
	}
	return &update, nil
}

func sendTextToTelegramChat(chatId int, text string) (string, error) {
	log.Printf("Sending %s to chat_id: %d", text, chatId)

	var telegramApi string = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage"

	params := SendMessageParams{
		ChatId: strconv.Itoa(chatId),
		Text:   text,
	}

	form := url.Values{}

	form.Add("chat_id", params.ChatId)
	form.Add("text", params.Text)

	response, err := http.PostForm(telegramApi, form)
	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())
		return "", err
	}
	bodyString := string(bodyBytes)
	log.Printf("Body of Telegram Response: %s", bodyString)

	return bodyString, nil
}

func createButtons(chatId int) error {
	// Set up the URL for the Telegram API
	var telegramApi string = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage"

	keyboard := [][]KeyboardButton{
		{
			{Text: "Button 1.1"},
			{Text: "Button 1.2"},
		},
		{
			{Text: "Button 2.1"},
			{Text: "Button 2.2"},
			{Text: "Button 2.3"},
			{Text: "Button 2.4"},
		},
	}

	replyMarkup := &ReplyKeyboardMarkup{
		Keyboard:       keyboard,
		ResizeKeyboard: true,
	}

	params := SendMessageParams{
		ChatId:      strconv.Itoa(chatId),
		Text:        "Ready to output your messages, Sir!",
		ReplyMarkup: replyMarkup,
	}

	replyMarkupJson, err := json.Marshal(params.ReplyMarkup)
	if err != nil {
		log.Println(err)
		return err
	}

	form := url.Values{}
	form.Add("chat_id", params.ChatId)
	form.Add("text", params.Text)
	form.Add("reply_markup", string(replyMarkupJson))

	response, err := http.PostForm(telegramApi, form)
	if err != nil {
		log.Printf("/start: error when posting text to the chat: %s", err.Error())
		return err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())
		return err
	}
	bodyString := string(bodyBytes)
	log.Printf("Body of Telegram Response: %s", bodyString)

	return nil
}
