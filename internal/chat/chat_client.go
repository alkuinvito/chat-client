package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func SendMessage(port int, room ChatRoom, message ChatMessage) ([]byte, error) {
	payload := map[string]string{
		"sender":  message.Sender,
		"message": message.Message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	target := fmt.Sprintf("http://%s:%d/chat", room.IP.String(), port)
	resp, err := http.Post(target, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
