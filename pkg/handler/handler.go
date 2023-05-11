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

type SendMessageParams struct {
	ChatId      string               `json:"chat_id"`
	Text        string               `json:"text"`
	ReplyMarkup *ReplyKeyboardMarkup `json:"reply_markup,omitempty"`
}

type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard bool               `json:"one_time_keyboard,omitempty"`
	Selective       bool               `json:"selective,omitempty"`
}

type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact,omitempty"`
	RequestLocation bool   `json:"request_location,omitempty"`
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

// sendTextToTelegramChat sends a text message to the Telegram chat identified by its chat Id
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
		Text:        "Some text",
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

	// Create a slice of slices to represent the inline keyboard
	// buttons := [][]string{
	// 	{"Button 1", "Button 2"},
	// 	{"Button 3", "Button 4"},
	// }

	// opts := map[string]interface{}{
	// 	"reply_markup": map[string]interface{}{
	// 		"inline_keyboard": [][]map[string]interface{}{
	// 			{
	// 				{
	// 					"text": "A",
	// 				},
	// 				{
	// 					"text": "B",
	// 				},
	// 			},
	// 		},
	// 	},
	// }
	// optsJSON, err := json.Marshal(opts)
	// if err != nil {
	// 	log.Println("Error:", err)
	// 	return err
	// }
	// fmt.Println(string(optsJSON))

	// response, err := http.PostForm(
	// 	telegramApi,
	// 	url.Values{
	// 		"chat_id":      {strconv.Itoa(chatId)},
	// 		"text":         {"Some text"},
	// 		"reply_markup": {string(optsJSON)},
	// 	})

	// Convert the keyboard to JSON format
	// keyboard, err := json.Marshal(map[string]interface{}{
	// 	"inline_keyboard": buttons,
	// })
	// if err != nil {
	// 	log.Println(err)
	// }

	// Create the message payload
	// payload, err := json.Marshal(map[string]interface{}{
	// 	"chat_id":      strconv.Itoa(chatId),
	// 	"text":         "Please select an option:",
	// 	"reply_markup": keyboard,
	// })
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }

	// Send the message to the Telegram API
	// resp, err := http.Post(telegramApi, "application/json", bytes.NewBuffer(payload))
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }

	// Parse the response from the Telegram API
	// var telegramResp telegramResponse
	// err = json.NewDecoder(resp.Body).Decode(&telegramResp)
	// if err != nil {
	// 	return err
	// }

	// var bodyBytes, errRead = ioutil.ReadAll(resp.Body)
	// if errRead != nil {
	// 	log.Printf("error in parsing telegram answer %s", errRead.Error())
	// 	return err
	// }

	// bodyString := string(bodyBytes)
	// log.Printf("Body of Telegram Response to \\start : %s\n", bodyString)

	// // Check if the response was successful
	// if !telegramResp.Ok {
	// 	log.Fatalf("Error sending message: %s", resp.Status)
	// }

	// log.Println("Message sent!")

	// response, err := http.PostForm(
	// 	telegramApi,
	// 	url.Values{
	// 		"chat_id": {strconv.Itoa(chatId)},
	// 		"text":    {text},
	// 	})

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
