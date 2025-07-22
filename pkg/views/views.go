package views

import (
	"fyne.io/fyne/v2"
)

type View struct {
	a fyne.App
	w fyne.Window
}

func NewView(a fyne.App, title string) *View {
	w := a.NewWindow(title)
	return &View{a, w}
}

func (v *View) Render(view func(*View) fyne.CanvasObject) {
	v.w.SetContent(view(v))
}

func (v *View) Run() {
	v.w.ShowAndRun()
}

func (v *View) Window() fyne.Window {
	return v.w
}
