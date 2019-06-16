package ui

import (
	"github.com/mattn/go-gtk/gtk"
)

func (pCtx *PaysimUiContext) createButtonsAndBox() {

	pCtx.btnLd = gtk.NewButtonWithLabel(" Load Trace ")
	pCtx.btnLd.ModifyFontEasy("Dejavu Sans 9")
	pCtx.btnAssemble = gtk.NewButtonWithLabel(" Assemble Trace ")
	pCtx.btnSend = gtk.NewButtonWithLabel(" Send  ")
	pCtx.btnBox = gtk.NewHBox(false, 5)
	pCtx.btnBox.PackStart(pCtx.btnLd, false, false, 1)
	pCtx.btnBox.PackStart(pCtx.btnAssemble, false, false, 1)
	pCtx.btnBox.PackStart(pCtx.btnSend, false, false, 1)
}

func (pCtx *PaysimUiContext) ButtonBox() *gtk.HBox {
	return pCtx.btnBox
}
