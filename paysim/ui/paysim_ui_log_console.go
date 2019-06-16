package ui

import (
	_ "github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/gtk"
)

type PaysimUiConsole struct {
	txtBuf  *gtk.TextBuffer
	txtView *gtk.TextView
}

var init_ bool = false
var PaysimConsole *PaysimUiConsole

func Init() {
	if !init_ {
		PaysimConsole = newConsole()
		init_ = true
	}
}

func newConsole() *PaysimUiConsole {

	self := new(PaysimUiConsole)
	self.txtBuf = gtk.NewTextBuffer(gtk.NewTextTagTable())
	self.txtView = gtk.NewTextViewWithBuffer(*self.txtBuf)
	self.txtView.SetEditable(false)
	self.txtView.ModifyFontEasy("consolas 8")
	return self
}

func (pc *PaysimUiConsole) TextView() *gtk.TextView {
	return pc.txtView
}

func (pc *PaysimUiConsole) Log(log string) {
	var endIter gtk.TextIter
	pc.txtBuf.GetEndIter(&endIter)
	pc.txtBuf.Insert(&endIter, log)
}
