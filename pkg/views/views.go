package views

import (
	"fyne.io/fyne/v2"
	"gorm.io/gorm"
)

type View struct {
	a  fyne.App
	w  fyne.Window
	db *gorm.DB
}

func NewView(a fyne.App, title string, db *gorm.DB) *View {
	w := a.NewWindow(title)
	return &View{a, w, db}
}

func (v *View) DB() *gorm.DB {
	return v.db
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
