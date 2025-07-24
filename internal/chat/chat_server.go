package chat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alkuinvito/chat-client/internal/discovery"
)

func ServeChat(msgStream chan *ChatMessage) {
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload ChatMessage
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&payload)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		msgStream <- &payload

		fmt.Printf("%s: %s\n", r.RemoteAddr, payload.Message)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"received"}`))
	})

	go func() {
		fmt.Printf("Server is running on: 127.0.0.1:%d\n", discovery.SVC_PORT)
		err := http.ListenAndServe(fmt.Sprintf(":%d", discovery.SVC_PORT), nil)
		if err != nil {
			panic(err)
		}
	}()
}
