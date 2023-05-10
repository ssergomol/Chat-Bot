package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

// Update is a Telegram object that the handler receives every time an user interacts with the bot.
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

// Message is a Telegram object that can be found in an update.
type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

// A Telegram Chat indicates the conversation to which the message belongs.
type Chat struct {
	Id int `json:"id"`
}

// Define a struct to represent the JSON response from the Telegram API
type telegramResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageId int `json:"message_id"`
		Chat      struct {
			Id int `json:"id"`
		} `json:"chat"`
	} `json:"result"`
}

// parseTelegramRequest handles incoming update from the Telegram web hook
func parseTelegramRequest(r *http.Request) (*Update, error) {
	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode incoming update %s", err.Error())
		return nil, err
	}
	return &update, nil
}

func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {
	// Parse incoming request
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
	}

	outputMessage := "Your message: " + update.Message.Text

	// Send the punchline back to Telegram
	var telegramResponseBody, errTelegram = sendTextToTelegramChat(update.Message.Chat.Id, outputMessage)
	if errTelegram != nil {
		log.Printf("got error %s from telegram, reponse body is %s", errTelegram.Error(), telegramResponseBody)
	} else {
		log.Printf("message \"%s\" successfully sent to chat id %d\n", outputMessage, update.Message.Chat.Id)
	}
}

// sendTextToTelegramChat sends a text message to the Telegram chat identified by its chat Id
func sendTextToTelegramChat(chatId int, text string) (string, error) {
	log.Printf("Sending %s to chat_id: %d", text, chatId)

	var telegramApi string = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage"
	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		})

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
	apiUrl := "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage"

	// Create a slice of slices to represent the inline keyboard
	buttons := [][]string{
		{"Button 1", "Button 2"},
		{"Button 3", "Button 4"},
	}

	// Convert the keyboard to JSON format
	keyboard, err := json.Marshal(map[string]interface{}{
		"inline_keyboard": buttons,
	})
	if err != nil {
		log.Println(err)
	}

	// Create the message payload
	payload, err := json.Marshal(map[string]interface{}{
		"chat_id":      strconv.Itoa(chatId),
		"text":         "Please select an option:",
		"reply_markup": keyboard,
	})
	if err != nil {
		log.Println(err)
		return err
	}

	// Send the message to the Telegram API
	resp, err := http.Post(apiUrl, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Println(err)
		return err
	}

	// Parse the response from the Telegram API
	var telegramResp telegramResponse
	err = json.NewDecoder(resp.Body).Decode(&telegramResp)
	if err != nil {
		return err
	}

	log.Println(resp.Body)

	// Check if the response was successful
	if !telegramResp.Ok {
		log.Fatalf("Error sending message: %s", resp.Status)
	}

	log.Println("Message sent!")
	return nil
}
