package auth

import (
	"errors"
	"regexp"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/alkuinvito/chat-client/internal/chat"
	"github.com/alkuinvito/chat-client/pkg/views"
)

func RegisterView(v *views.View) fyne.CanvasObject {
	alphanum := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	username := widget.NewEntry()

	username.Validator = func(s string) error {
		if len(s) < 4 || len(s) > 16 {
			return errors.New("username must be between 4 and 16 characters")
		}

		if !alphanum.MatchString(s) {
			return errors.New("username must only contains alphanumeric")
		}

		return nil
	}

	form := &widget.Form{
		Items: []*widget.FormItem{
			widget.NewFormItem("Username", username),
		},
		SubmitText: "Register",
		OnSubmit: func() {
			// RegisterUser(v.DB(), username.Text)
			v.Store().Set("username", username.Text)
			v.Render(chat.ChatView)
		}}

	gridCt := container.NewGridWrap(fyne.NewSize(400, form.MinSize().Height), form)
	ct := container.NewCenter(gridCt)

	return ct
}
