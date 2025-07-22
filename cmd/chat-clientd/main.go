package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/alkuinvito/chat-client/internal/auth"
	"github.com/alkuinvito/chat-client/pkg/db"
	"github.com/alkuinvito/chat-client/pkg/views"
)

func main() {
	db := db.NewDB()

	a := app.New()
	view := views.NewView(a, "P2P Chat", db)
	view.Window().Resize(fyne.NewSize(900, 600))

	// register new user as entrypoint
	view.Render(auth.RegisterView)

	view.Run()
}
