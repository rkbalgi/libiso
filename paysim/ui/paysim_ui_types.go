package ui

import (
	"encoding/hex"
	"errors"
	_ "fmt"

	_ "github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"net"
	"regexp"
	"strconv"
	"strings"
)

type SpecTab struct {
}

type PaysimUiContext struct {
	_window         *gtk.Window
	ipAddrEntry     *gtk.Entry
	portEntry       *gtk.Entry
	mliTypesCb      *gtk.ComboBox
	commsConfigVbox *gtk.VBox

	btnLd       *gtk.Button
	btnAssemble *gtk.Button
	btnSend     *gtk.Button
	btnBox      *gtk.HBox
}

func (pCtx *PaysimUiContext) Window() *gtk.Window {
	return pCtx._window
}

func (pCtx *PaysimUiContext) LoadButton() *gtk.Button {
	return pCtx.btnLd
}
func (pCtx *PaysimUiContext) AssembleButton() *gtk.Button {
	return pCtx.btnAssemble
}
func (pCtx *PaysimUiContext) SendButton() *gtk.Button {
	return pCtx.btnSend
}

/*func (ctx *PaysimUiContext) MliTypesComboBox() *gtk.ComboBox {
	return ctx.mli_types_cb
}

func (ctx *PaysimUiContext) IpAdrrWidget() *gtk.Entry {
	return ctx.ip_addr_entry
}

func (ctx *PaysimUiContext) PortEntryWidget() *gtk.Entry {
	return ctx.port_entry
}
*/
func (pCtx *PaysimUiContext) CommsConfigVBox() *gtk.VBox {
	return pCtx.commsConfigVbox
}

func NewUiContext() *PaysimUiContext {

	ctx := new(PaysimUiContext)
	ctx._window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)

	ctx._window.SetPosition(gtk.WIN_POS_CENTER)
	ctx._window.SetTitle("PaySim v1.00 - An ISO8583 Simulator")
	ctx._window.SetIconName("gtk-dialog-info")
	ctx._window.SetSizeRequest(600, 600)
	ctx._window.Connect("destroy", func(ctx *glib.CallbackContext) {
		gtk.MainQuit()
	}, "foo")

	//comms_config box setup
	ctx.ipAddrEntry = gtk.NewEntry()
	ctx.portEntry = gtk.NewEntry()
	ctx.mliTypesCb = gtk.NewComboBoxNewText()
	ctx.mliTypesCb.AppendText("2I")
	ctx.mliTypesCb.AppendText("2E")

	iter := gtk.TreeIter{}
	ctx.mliTypesCb.GetModel().GetIterFirst(&iter)
	ctx.mliTypesCb.SetActiveIter(&iter)

	ctx.constructCommsConfigVbox()
	ctx.createButtonsAndBox()

	return ctx
}

//GetCommsConfig returns the IP, port and the MLI to be used
//for the message
func (pCtx *PaysimUiContext) GetCommsConfig() (*net.TCPAddr, string, error) {

	ipStr := pCtx.ipAddrEntry.GetText()
	if len(strings.Trim(ipStr, " ")) == 0 {
		return nil, "", errors.New("invalid ip")
	}

	portStr := pCtx.portEntry.GetText()
	port := uint64(0)
	matched, err := regexp.Match("^[0-9]{4,6}$", []byte(portStr))
	if err != nil {
		return nil, "", err
	} else {
		if !matched {
			return nil, "", errors.New("invalid port")
		}

		_ip := net.ParseIP(ipStr)
		port, err = strconv.ParseUint(portStr, 10, 32)
		tcpAddr := &net.TCPAddr{IP: _ip, Port: int(port)}
		return tcpAddr, pCtx.mliTypesCb.GetActiveText(), nil
	}

}

func (pCtx *PaysimUiContext) constructCommsConfigVbox() {

	//default values
	pCtx.ipAddrEntry.SetText("127.0.0.1")
	pCtx.portEntry.SetText("5656")

	pCtx.commsConfigVbox = gtk.NewVBox(true, 5)
	//ip addr box
	tmpHbox := gtk.NewHBox(false, 5)
	tmpHbox.PackStart(gtk.NewLabel("Destination Ip   "), false, false, 5)
	tmpHbox.PackStart(pCtx.ipAddrEntry, false, false, 5)
	pCtx.commsConfigVbox.PackStart(tmpHbox, false, false, 5)

	//port box
	tmpHbox = gtk.NewHBox(false, 5)
	tmpHbox.PackStart(gtk.NewLabel("Destination Port "), false, false, 2)
	tmpHbox.PackStart(pCtx.portEntry, false, false, 5)
	pCtx.commsConfigVbox.PackStart(tmpHbox, false, false, 1)

	tmpHbox = gtk.NewHBox(false, 5)
	tmpHbox.PackStart(gtk.NewLabel("MLI Type "), false, false, 10)
	tmpHbox.PackStart(pCtx.mliTypesCb, false, false, 15)
	pCtx.commsConfigVbox.PackStart(tmpHbox, false, false, 1)

}

func (pCtx *PaysimUiContext) GetUsrTrace() ([]byte, error) {

	dialog := gtk.NewDialog()
	dialog.SetPosition(gtk.WIN_POS_CENTER)
	dialog.SetTitle("Input Trace")
	dialog.SetParent(pCtx.Window().GetTopLevelAsWindow())
	dialog.SetModal(true)
	dialog.SetSizeRequest(400, 200)

	data := make([]byte, 0)
	var err error

	textView := gtk.NewTextView()
	textView.SetEditable(true)
	textView.SetWrapMode(gtk.WRAP_CHAR)
	textView.GetBuffer().SetText("31313030702000000000000131353131313131313131313131313131343030343830303030303030303030303132323132333435360000000000000000")
	scrolledWindow := gtk.NewScrolledWindow(nil, nil)
	scrolledWindow.AddWithViewPort(textView)

	dialog.GetVBox().PackStart(scrolledWindow, true, true, 5)
	okBtn := dialog.AddButton("OK", gtk.RESPONSE_OK)

	okBtn.Connect("clicked", func() {
		textBuf := textView.GetBuffer()
		var sI, eI gtk.TextIter
		textBuf.GetStartIter(&sI)
		textBuf.GetEndIter(&eI)
		tmp := textBuf.GetText(&sI, &eI, false)
		data, err = hex.DecodeString(strings.Trim(tmp, " "))

		if err != nil {
			ShowErrorDialog(dialog, "invalid trace data")
		} else {

			dialog.Destroy()
			gtk.MainQuit()
		}
	})
	cancelBtn := dialog.AddButton("Cancel", gtk.RESPONSE_CANCEL)
	cancelBtn.Connect("clicked", func() {

		err = errors.New("cancel_op")
		data = nil
		dialog.Destroy()
		gtk.MainQuit()
	})

	dialog.ShowAll()
	gtk.Main()
	return data, err

}

func (pCtx *PaysimUiContext) ShowUsrTrace(traceData []byte) {

	dialog := gtk.NewDialog()
	dialog.SetTitle("Assembled Trace")
	dialog.SetPosition(gtk.WIN_POS_CENTER)
	dialog.SetModal(true)
	dialog.SetSizeRequest(400, 200)

	textView := gtk.NewTextView()
	textView.SetEditable(false)
	textView.SetWrapMode(gtk.WRAP_CHAR)
	textView.GetBuffer().SetText(hex.EncodeToString(traceData))
	ScrolledWindow := gtk.NewScrolledWindow(nil, nil)
	ScrolledWindow.AddWithViewPort(textView)

	dialog.GetVBox().PackStart(ScrolledWindow, true, true, 5)
	okBtn := dialog.AddButton("OK", gtk.RESPONSE_OK)

	okBtn.Connect("clicked", func() {
		dialog.Destroy()
	})

	dialog.ShowAll()

}
