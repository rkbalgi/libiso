package ui

import(
		"github.com/mattn/go-gtk/gtk"
)

func (self *PaysimUiContext) create_buttons_and_box(){

	self._btn_ld = gtk.NewButtonWithLabel(" Load Trace ")
	self._btn_ld.ModifyFontEasy("Dejavu Sans 9");
	self._btn_assemble = gtk.NewButtonWithLabel(" Assemble Trace ")
	self._btn_send = gtk.NewButtonWithLabel(" Send  ")
    self._btn_box = gtk.NewHBox(false, 5)
	self._btn_box.PackStart(self._btn_ld, false, false, 1)
	self._btn_box.PackStart(self._btn_assemble, false, false, 1)
	self._btn_box.PackStart(self._btn_send, false, false, 1)
}

func (self *PaysimUiContext) ButtonBox() *gtk.HBox{
	return self._btn_box;
}