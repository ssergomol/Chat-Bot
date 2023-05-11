package handler

// Data sent by user
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

type Chat struct {
	Id int `json:"id"`
}

// Parameters to send from server to Telegram Chat
type SendMessageParams struct {
	ChatId      string               `json:"chat_id"`
	Text        string               `json:"text"`
	ReplyMarkup *ReplyKeyboardMarkup `json:"reply_markup,omitempty"`
}

// Keyboard information
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
