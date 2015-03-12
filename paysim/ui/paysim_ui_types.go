package ui

import (
	"encoding/hex"
	"errors"
	_ "fmt"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"strings"
)

type SpecTab struct {
}

type PaysimUiContext struct {
	_window           *gtk.Window
	ip_addr_entry     *gtk.Entry
	port_entry        *gtk.Entry
	mli_types_cb      *gtk.ComboBox
	comms_config_vbox *gtk.VBox
}

func (ctx *PaysimUiContext) Window() *gtk.Window {
	return ctx._window
}

func (ctx *PaysimUiContext) MliTypesComboBox() *gtk.ComboBox {
	return ctx.mli_types_cb
}

func (ctx *PaysimUiContext) IpAdrrWidget() *gtk.Entry {
	return ctx.ip_addr_entry
}

func (ctx *PaysimUiContext) PortEntryWidget() *gtk.Entry {
	return ctx.port_entry
}

func (ctx *PaysimUiContext) CommsConfigVBox() *gtk.VBox {
	return ctx.comms_config_vbox
}

func NewUiContext() *PaysimUiContext {

	ctx := new(PaysimUiContext)
	ctx._window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)

	ctx._window.SetPosition(gtk.WIN_POS_CENTER)
	ctx._window.SetTitle("PaySim v1.00 - An ISO8583 Simulator")
	ctx._window.SetIconName("gtk-dialog-info")
	ctx._window.SetSizeRequest(gdk.ScreenWidth(), gdk.ScreenHeight()-40)
	ctx._window.Connect("destroy", func(ctx *glib.CallbackContext) {
		gtk.MainQuit()
	}, "foo")

	//comms_config box setup
	ctx.ip_addr_entry = gtk.NewEntry()
	ctx.port_entry = gtk.NewEntry()
	ctx.mli_types_cb = gtk.NewComboBoxNewText()
	ctx.mli_types_cb.AppendText("2I")
	ctx.mli_types_cb.AppendText("2E")

	iter := gtk.TreeIter{}
	ctx.mli_types_cb.GetModel().GetIterFirst(&iter)
	ctx.mli_types_cb.SetActiveIter(&iter)

	ctx.construct_comms_config_vbox()

	return ctx
}

func (ctx *PaysimUiContext) construct_comms_config_vbox() {

	ctx.comms_config_vbox = gtk.NewVBox(false, 5)
	//ip addr box
	tmp_hbox := gtk.NewHBox(false, 5)
	tmp_hbox.PackStart(gtk.NewLabel("Destination Ip   "), false, false, 10)
	tmp_hbox.PackStart(ctx.ip_addr_entry, false, false, 20)
	ctx.comms_config_vbox.PackStart(tmp_hbox, false, false, 1)

	//port box
	tmp_hbox = gtk.NewHBox(false, 5)
	tmp_hbox.PackStart(gtk.NewLabel("Destination Port "), false, false, 10)
	tmp_hbox.PackStart(ctx.port_entry, false, false, 20)
	ctx.comms_config_vbox.PackStart(tmp_hbox, false, false, 1)

	tmp_hbox = gtk.NewHBox(false, 5)
	tmp_hbox.PackStart(gtk.NewLabel("MLI Type "), false, false, 10)
	tmp_hbox.PackStart(ctx.mli_types_cb, false, false, 20)
	ctx.comms_config_vbox.PackStart(tmp_hbox, false, false, 1)

}

func (ctx *PaysimUiContext) GetUsrTrace() ([]byte, error) {

	dialog := gtk.NewDialog()
	dialog.SetPosition(gtk.WIN_POS_CENTER)
	dialog.SetTitle("Input Trace")
	dialog.SetParent(ctx.Window().GetTopLevelAsWindow())
	dialog.SetModal(true)
	dialog.SetSizeRequest(400, 200)

	data := make([]byte, 0)
	var err error

	text_view := gtk.NewTextView()
	text_view.SetEditable(true)
	text_view.SetWrapMode(gtk.WRAP_CHAR)
	text_view.GetBuffer().SetText("31313030702000000000000131353337313131313131313131313131343030343830303030303030303030303132323132333435360000000000000000")
	swin := gtk.NewScrolledWindow(nil, nil)
	swin.AddWithViewPort(text_view)

	dialog.GetVBox().PackStart(swin, true, true, 5)
	ok_btn := dialog.AddButton("OK", gtk.RESPONSE_OK)

	ok_btn.Connect("clicked", func() {
		text_buf := text_view.GetBuffer()
		var s_i, e_i gtk.TextIter
		text_buf.GetStartIter(&s_i)
		text_buf.GetEndIter(&e_i)
		tmp := text_buf.GetText(&s_i, &e_i, false)
		data, err = hex.DecodeString((strings.Trim(tmp, " ")))

		if err != nil {
			ShowErrorDialog(dialog, "invalid trace data")
		} else {

			dialog.Destroy()
			gtk.MainQuit()
		}
	})
	cancel_btn := dialog.AddButton("Cancel", gtk.RESPONSE_CANCEL)
	cancel_btn.Connect("clicked", func() {

		err = errors.New("cancel_op")
		data = nil
		dialog.Destroy()
		gtk.MainQuit()
	})

	dialog.ShowAll()
	gtk.Main()
	return data, err

}

func (ctx *PaysimUiContext) ShowUsrTrace(trace_data []byte) {

	dialog := gtk.NewDialog()
	dialog.SetTitle("Assembled Trace")
	dialog.SetModal(true)
	dialog.SetSizeRequest(400, 200)

	text_view := gtk.NewTextView()
	text_view.SetEditable(false)
	text_view.SetWrapMode(gtk.WRAP_CHAR)
	text_view.GetBuffer().SetText(hex.EncodeToString(trace_data))
	swin := gtk.NewScrolledWindow(nil, nil)
	swin.AddWithViewPort(text_view)

	dialog.GetVBox().PackStart(swin, true, true, 5)
	ok_btn := dialog.AddButton("OK", gtk.RESPONSE_OK)

	ok_btn.Connect("clicked", func() {
		dialog.Destroy()
	})

	dialog.ShowAll()

}
