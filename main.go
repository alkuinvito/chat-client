package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/alkuinvito/chat-client/internal/auth"
	"github.com/alkuinvito/chat-client/internal/chat"
	"github.com/alkuinvito/chat-client/internal/discovery"
	"github.com/alkuinvito/chat-client/pkg/db"
	"github.com/alkuinvito/chat-client/pkg/views"
)

func main() {
	db := db.NewDB()

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload chat.ChatMessage
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&payload)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

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

	a := app.New()
	view := views.NewView(a, "P2P Chat", db)
	view.Window().Resize(fyne.NewSize(900, 600))

	// register new user as entrypoint
	view.Render(auth.RegisterView)

	view.Run()
}
