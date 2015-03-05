package main

import (
	_ "fmt"
	"github.com/rkbalgi/go/iso8583"
	_ "github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

func main() {

	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)

	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle("PaySim v1.00 - An ISO8583 Simulator")
	window.SetIconName("gtk-dialog-info")
	window.SetSizeRequest(800, 600)
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		println("got destroy!", ctx.Data().(string))
		gtk.MainQuit()
	}, "foo")

	menu_bar := gtk.NewMenuBar()
	cascade_file_mi := gtk.NewMenuItemWithMnemonic("_File")
	file_menu := gtk.NewMenu()
	cascade_file_mi.SetSubmenu(file_menu)
	menu_bar.Append(cascade_file_mi)

	open_spec_mi := gtk.NewMenuItemWithMnemonic("_Open Specs Def")
	file_menu.Add(open_spec_mi)
	open_spec_mi.Connect("activate", func() {

		filechooserdialog := gtk.NewFileChooserDialog(
			"Choose File...",
			window,
			gtk.FILE_CHOOSER_ACTION_OPEN,
			gtk.STOCK_OK,
			gtk.RESPONSE_ACCEPT)
		filter := gtk.NewFileFilter()
		filter.AddPattern("*.json")
		filechooserdialog.AddFilter(filter)
		filechooserdialog.Response(func() {
			iso8583.ReadSpecDefs(filechooserdialog.GetFilename());
			filechooserdialog.Destroy()
		})
		filechooserdialog.Run()

	})

	quit_mi := gtk.NewMenuItemWithMnemonic("_Quit")
	quit_mi.Connect("activate", func() {
		gtk.MainQuit()
	})
	file_menu.Add(quit_mi)

	menu_bar.Append(file_menu)

	vbox := gtk.NewVBox(false, 1)
	vbox.PackStart(menu_bar, false, false, 0)

	//add a vertical pane
	h_pane := gtk.NewHPaned()
	h_pane.SetPosition(40)
	vbox.Add(h_pane)

	frame1 := gtk.NewFrame("")
	framebox1 := gtk.NewVBox(false, 1)
	frame1.Add(framebox1)

	frame2 := gtk.NewFrame("")
	framebox2 := gtk.NewVBox(false, 1)
	frame2.Add(framebox2)
	h_pane.Pack1(frame1, false, false)
	h_pane.Pack2(frame2, false, false)

	window.Add(vbox)

	window.ShowAll()

	gtk.Main()

}
