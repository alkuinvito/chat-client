package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/alkuinvito/chat-client/pkg/views"
)

func main() {
	a := app.New()
	view := views.NewView(a, "P2P Chat")

	view.Run()
}
