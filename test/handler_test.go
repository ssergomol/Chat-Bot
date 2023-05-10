package test

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/ssergomol/Chat-Bot/pkg/handler"
)

func TestParseUpdateMessageWithText(t *testing.T) {
	var msg = handler.Message{
		Text: "hello world",
		Chat: handler.Chat{Id: 1},
	}

	var update = handler.Update{
		UpdateId: 1,
		Message:  msg,
	}

	requestBody, err := json.Marshal(update)
	if err != nil {
		t.Errorf("Failed to marshal update in json, got %s", err.Error())
	}
	req := httptest.NewRequest("POST", "http:/127.0.0.1:8080/", bytes.NewBuffer(requestBody))

	var updateToTest, errParse = handler.ParseTelegramRequest(req)
	if errParse != nil {
		t.Errorf("Expected a <nil> error, got %s", errParse.Error())
	}
	if *updateToTest != update {
		t.Errorf("Expected update %s, got %s", update.Message.Text, updateToTest.Message.Text)
	}

}
