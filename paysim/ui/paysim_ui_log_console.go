package ui

import (
	_ "github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/gtk"
	
)



type PaysimUiConsole struct {
	_txt_buf  *gtk.TextBuffer
	_txt_view *gtk.TextView
}

var PaysimConsole *PaysimUiConsole=nil
var init_ bool = false

func Init() {
	if !init_ {
		PaysimConsole = _new_console()
		init_ = true
	}
}
func _new_console() *PaysimUiConsole {

	self := new(PaysimUiConsole)
	self._txt_buf = gtk.NewTextBuffer(gtk.NewTextTagTable())
	self._txt_view = gtk.NewTextViewWithBuffer(*self._txt_buf)
	self._txt_view.SetEditable(false)
	self._txt_view.ModifyFontEasy("consolas 10");
	return self
}

func (self *PaysimUiConsole) TextView() *gtk.TextView {
	return self._txt_view
}

func (self *PaysimUiConsole) Log(log string) {
	var end_iter gtk.TextIter
	self._txt_buf.GetEndIter(&end_iter)
	self._txt_buf.Insert(&end_iter, log)
}


