package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/mattn/go-gtk/gdk"
	_ "github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"github.com/rkbalgi/go/iso8583"
	pylog "github.com/rkbalgi/go/paysim/log"
	pynet "github.com/rkbalgi/go/paysim/net"
	pyui "github.com/rkbalgi/go/paysim/ui"
	"log"
	"net"
	"unsafe"
)

var rightBox, specTreeVbox *gtk.VBox
var statusBar *gtk.Statusbar
var statusBarContextId uint
var activeSpecFrame *gtk.Frame
var swinSpecTree *gtk.ScrolledWindow
var uiCtx *pyui.PaysimUiContext
var bmpIter gtk.TreeIter
var hBox *gtk.HBox
var swinConsole *gtk.ScrolledWindow

//the tree holding all the specs
var specTreeView *gtk.TreeView

//the currently loaded iso_msg
var reqIsoMsg, respIsoMsg *iso8583.Iso8583Message
var commsConfigBx *gtk.VBox

func main() {

	gtk.Init(nil)
	gtk.SettingsGetDefault().SetStringProperty("gtk-font-name", "Dejavu Sans 8", "")
	//a vbox that will hold the spec tree
	//and its contents
	pyui.Init()
	specTreeVbox = gtk.NewVBox(false, 1)

	uiCtx = pyui.NewUiContext()

	menu_bar := gtk.NewMenuBar()
	cascade_file_mi := gtk.NewMenuItemWithMnemonic("_File")
	file_menu := gtk.NewMenu()
	cascade_file_mi.SetSubmenu(file_menu)
	menu_bar.Append(cascade_file_mi)

	openSpecMi := gtk.NewMenuItemWithMnemonic("_Open Specs Def")
	file_menu.Add(openSpecMi)
	openSpecMi.Connect("activate", func() {

		filechooserdialog := gtk.NewFileChooserDialog(
			"Choose File...",
			uiCtx.Window(),
			gtk.FILE_CHOOSER_ACTION_OPEN,
			gtk.STOCK_OK,
			gtk.RESPONSE_ACCEPT)
		filter := gtk.NewFileFilter()
		filter.AddPattern("*.json")
		filechooserdialog.AddFilter(filter)
		filechooserdialog.Response(func() {
			iso8583.ReadSpecDefs(filechooserdialog.GetFilename())
			makeAndPopulateSpecTree()
			filechooserdialog.Destroy()
		})
		filechooserdialog.Run()

	})

	quitMi := gtk.NewMenuItemWithMnemonic("E_xit")
	quitMi.Connect("activate", func() {
		gtk.MainQuit()
	})
	file_menu.Add(quitMi)

	cascadeUtilsMi := gtk.NewMenuItemWithMnemonic("_Utils")
	utilsMenu := gtk.NewMenu()
	cascadeUtilsMi.SetSubmenu(utilsMenu)

	macMi := gtk.NewMenuItemWithMnemonic("_MAC")
	utilsMenu.Add(macMi)
	macMi.Connect("activate", func() {
		pyui.ComputeMacDialog(uiCtx.Window(), "")
	})

	pinMi := gtk.NewMenuItemWithMnemonic("_PIN")
	utilsMenu.Add(pinMi)
	pinMi.Connect("activate", func() {
		pyui.ComputePinBlockDialog(uiCtx.Window(), "")
	})

	menu_bar.Append(cascadeUtilsMi)

	cascadeAboutMi := gtk.NewMenuItemWithMnemonic("_About")

	cascadeAboutMi.Connect("activate", func() {
		dialog := gtk.NewAboutDialog()
		dialog.SetName("About Paysim")
		dialog.SetProgramName("PaySim v1.0 build 03052015	")
		dialog.SetAuthors([]string{"Raghavendra Balgi (rkbalgi@gmail.com)"})

		dialog.SetLicense("This application is available under the same terms and conditions as the Go, the BSD style license, and the LGPL (Lesser GNU Public License).")
		dialog.SetWrapLicense(true)
		dialog.Run()
		dialog.Destroy()
	})

	menu_bar.Append(cascadeAboutMi)

	vbox := gtk.NewVBox(false, 1)
	vbox.PackStart(menu_bar, false, false, 0)

	//add a vertical pane
	hPane := gtk.NewHPaned()
	hPane.SetPosition(40)
	vbox.Add(hPane)

	frame1 := gtk.NewFrame("")
	frame1.SetSizeRequest(100, 600)
	frame1.Add(specTreeVbox)

	frame2 := gtk.NewFrame("")
	rightBox = gtk.NewVBox(false, 1)
	frame2.Add(rightBox)
	hPane.Pack1(frame1, false, false)
	hPane.Pack2(frame2, false, false)

	uiCtx.Window().Add(vbox)

	statusBar = gtk.NewStatusbar()
	statusBarContextId = statusBar.GetContextId("pay-sim")
	statusBar.Push(statusBarContextId, "OK")
	vbox.PackStart(statusBar, false, false, 1)

	rightBox.PackStart(uiCtx.CommsConfigVBox(), false, false, 12)

	resetConsole()
	activeSpecFrame = gtk.NewFrame("                 ")
	activeSpecFrame.SetName("spec_frame")
	activeSpecFrame.SetSizeRequest(400, 400)

	rightBox.PackStart(activeSpecFrame, false, false, 5)
	rightBox.PackStart(swinConsole, false, false, 5)
	rightBox.ShowAll()

	pylog.Log("Paysim v1.00 starting...\n#********************************#")

	iso8583.ReadDemoSpecDefs()
	makeAndPopulateSpecTree()

	uiCtx.Window().ShowAll()
	gtk.Main()

}

func makeAndPopulateSpecTree() {
	//remove the old tree if there is already one
	//and then create a new one with the specs loaded in the
	//tree
	if specTreeView != nil {
		specTreeVbox.Remove(swinSpecTree)
	}

	treeStore := gtk.NewTreeStore(glib.G_TYPE_STRING)
	specTreeView = gtk.NewTreeView()
	specTreeView.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Specs", gtk.NewCellRendererText(), "text", 0))
	specTreeView.SetModel(treeStore.ToTreeModel())

	specTreeView.SetSizeRequest(250, 800)

	var ti1, ti2 gtk.TreeIter
	treeStore.Append(&ti1, nil)
	treeStore.Set(&ti1, "Specifications")
	n := 0
	for _, spec := range iso8583.GetSpecs() {
		treeStore.Append(&ti2, &ti1)
		treeStore.Set(&ti2, spec.Name())
		n++
	}

	specTreeView.SetHeadersVisible(false)

	specTreeView.Connect("button-release-event", func(ctx *glib.CallbackContext) {

		event := (*gdk.EventButton)(unsafe.Pointer(ctx.Args(0)))
		if event.Button == 1 {

			var path *gtk.TreePath
			var col *gtk.TreeViewColumn
			specTreeView.GetCursor(&path, &col)

			if path.String() == "0" {
				return
			}

			val := new(glib.GValue)
			var tIter gtk.TreeIter
			if treeStore.GetIter(&tIter, path) {
				treeStore.GetValue(&tIter, 0, val)
				showSpecLayout(val.GetString())

			}

		}

	})

	statusBar.Push(statusBarContextId, fmt.Sprintf("%d specs loaded. OK.", n))

	swinSpecTree = gtk.NewScrolledWindow(nil, nil)
	swinSpecTree.SetSizeRequest(200, 800)
	swinSpecTree.AddWithViewPort(specTreeView)

	specTreeVbox.PackStart(swinSpecTree, true, true, 2)

	specTreeVbox.ShowAll()

}

func showSpecLayout(specName string) {

	//this 'req_iso_msg' will act as the model
	//that backs the TreeView
	reqIsoMsg = iso8583.NewIso8583Message(specName)

	//remove the active frame
	if activeSpecFrame != nil {
		rightBox.Remove(activeSpecFrame)
		rightBox.Remove(hBox)
	}
	//right_box.Remove(swin_console)

	specMsgTree := pyui.NewPaysimSpecMsgTree()

	specLytTree := specMsgTree.View()
	specLytStore := specMsgTree.Store()

	var i1 gtk.TreeIter

	specLyt := iso8583.GetSpecLayout(specName)
	for _, fieldDef := range specLyt {
		specLytStore.Append(&i1, nil)
		specLytStore.SetValue(&i1, 0, false)
		specLytStore.SetValue(&i1, 1, fieldDef.BitPosition)
		specLytStore.SetValue(&i1, 2, fieldDef.Name)
		specLytStore.SetValue(&i1, 3, "")
		specLytStore.SetValue(&i1, 4, fieldDef.Id)
		if fieldDef.Name == "Bitmap" || fieldDef.Name == "Message Type" {
			specLytStore.SetValue(&i1, 0, true)
		} else {
			specLytStore.SetValue(&i1, 0, false)
		}

		//bitmap will be non-editable
		if fieldDef.Name == "Bitmap" {
			//lets store the iter
			bmpIter = i1
			specLytStore.SetValue(&i1, 3, "00000000000000000000000000000000")
			specLytStore.SetValue(&i1, 5, false)
		} else {
			specLytStore.SetValue(&i1, 5, true)
		}

	}

	specMsgTree.FieldValueRenderer().Connect("edited",
		func(ctx *glib.CallbackContext) {

			fieldName := getCurrentFieldName(specLytTree, specLytStore)
			if fieldName == "Bitmap" {
				//do not allow edits on Bitmap field
				return
			}

			setCurrentFieldValue(specLytTree, specLytStore, ctx.Args(1).ToString())

		})

	specMsgTree.FieldToggleRenderer().Connect("toggled",
		func(ctx *glib.CallbackContext) {

			var path *gtk.TreePath
			var col *gtk.TreeViewColumn
			specLytTree.GetCursor(&path, &col)
			log.Println(col.GetTitle())
			if col.GetTitle() != "Is Selected" {
				return
			}

			var i1 gtk.TreeIter
			if specLytStore.GetIter(&i1, path) {

				fieldNameVal := glib.GValue{}
				specLytStore.GetValue(&i1, 2, &fieldNameVal)
				fmt.Println(fieldNameVal.GetString())

				if fieldNameVal.GetString() == "Bitmap" || fieldNameVal.GetString() == "Message Type" {
					return
				}

				val := glib.GValue{}
				bitPosVal := glib.GValue{}
				specLytStore.GetValue(&i1, 0, &val)
				specLytStore.GetValue(&i1, 1, &bitPosVal)

				specLytStore.SetValue(&i1, 0, !val.GetBool())

				if !val.GetBool() {
					reqIsoMsg.Bitmap().SetOn(bitPosVal.GetInt())
				} else {
					reqIsoMsg.Bitmap().SetOff(bitPosVal.GetInt())
				}

				recomputeBitmap(specLytStore)
			}

		})

	uiCtx.LoadButton().Connect("clicked", func() {

		trace_data, err := uiCtx.GetUsrTrace()
		log.Println("user trace: ", hex.EncodeToString(trace_data))
		if err == nil {
			msg_buf := bytes.NewBuffer(trace_data)
			err = reqIsoMsg.Parse(msg_buf)
			if err != nil {
				pyui.ShowErrorDialog(uiCtx.Window(), "Trace Parse Error")
			}

			populateModel(specLytStore)
		}

	})

	uiCtx.AssembleButton().Connect("clicked", func() {

		traceData := reqIsoMsg.Bytes()
		uiCtx.ShowUsrTrace(traceData)
	})
	uiCtx.SendButton().Connect("clicked", func() {

		tcpAddr, mliTypeStr, err := uiCtx.GetCommsConfig()
		if err != nil {
			pyui.ShowErrorDialog(uiCtx.Window(), err.Error())
			return
		}

		//send this msg, get the response and then display
		//the response
		respIsoMsg, err = pynet.SendIsoMsg(tcpAddr.String(), mliTypeStr, reqIsoMsg)
		if err != nil {

			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() {
				pyui.ShowErrorDialog(uiCtx.Window(), "Message Timed Out.")
			} else {
				pyui.ShowErrorDialog(uiCtx.Window(), err.Error())
			}
		} else {
			//hoo-hah! response received
			//display it as a dialog to the user
			pyui.ShowIsoResponseMsgDialog(respIsoMsg.TabularFormat())
		}

	})

	align := gtk.NewAlignment(0.5, 0.5, 0.0, 0.0)
	align.Add(uiCtx.ButtonBox())
	specLytTree.ShowAll()
	activeSpecFrame = gtk.NewFrame("       [" + specName + "]        ")
	activeSpecFrame.SetName("spec_frame")
	activeSpecFrame.SetSizeRequest(300, 300)

	//add a scrolled winow containing
	//the spec_lty_tree
	swin := gtk.NewScrolledWindow(nil, nil)
	swin.AddWithViewPort(specLytTree)
	hboxTmp := gtk.NewHBox(false, 10)
	hboxTmp.PackStart(swin, true, true, 10)

	activeSpecFrame.Add(hboxTmp)

	rightBox.SetSizeRequest(400, 200)
	rightBox.PackStart(activeSpecFrame, false, false, 20)
	rightBox.PackStart(align, false, false, 0)
	//reset_console()
	rightBox.ReorderChild(swinConsole, -1) //.PackStart(swin_console, false, false, 0)

	rightBox.ShowAll()

}

func resetConsole() {

	swinConsole = gtk.NewScrolledWindow(nil, nil)
	swinConsole.AddWithViewPort(pyui.PaysimConsole.TextView())
	swinConsole.SetSizeRequest(400, 150)

}

func getCurrentFieldName(isoMsgTree *gtk.TreeView, isoMsgTreeStore *gtk.TreeStore) string {

	var path *gtk.TreePath
	var col *gtk.TreeViewColumn
	isoMsgTree.GetCursor(&path, &col)
	var i1 gtk.TreeIter

	if isoMsgTreeStore.GetIter(&i1, path) {
		fieldNameVal := glib.GValue{}
		isoMsgTreeStore.GetValue(&i1, 2, &fieldNameVal)
		return fieldNameVal.GetString()
	}

	return ""

}

func setCurrentFieldValue(isoMsgTree *gtk.TreeView,
	isoMsgTreeStore *gtk.TreeStore, val string) {

	var path *gtk.TreePath
	var col *gtk.TreeViewColumn
	isoMsgTree.GetCursor(&path, &col)
	var i1 gtk.TreeIter

	if isoMsgTreeStore.GetIter(&i1, path) {

		fName := glib.GValue{}
		isoMsgTreeStore.GetValue(&i1, 2, &fName)
		log.Println("setting field to value ", fName.GetString(), "to ", val)
		fldId := glib.GValue{}
		isoMsgTreeStore.GetValue(&i1, 4, &fldId)
		reqIsoMsg.SetFieldData(fldId.GetInt(), val)
		//data might have been truncated or padded based on the
		//field definition
		newVal := reqIsoMsg.GetFieldDataById(fldId.GetInt()).String()
		isoMsgTreeStore.SetValue(&i1, 3, newVal)

	}

}

//recompute_bitmap will show a new value of the bitmap
//on the UI
func recomputeBitmap(isoTreeStore *gtk.TreeStore) {

	val := glib.GValue{}
	bmpData := reqIsoMsg.Bitmap().Bytes()
	isoTreeStore.GetValue(&bmpIter, 3, &val)
	log.Println("Bitmap's current value", val.GetString())

	isoTreeStore.SetValue(&bmpIter, 3, hex.EncodeToString(bmpData))

}

//populates the store based on the data in
//req_iso_msg
func populateModel(isoTreeStore *gtk.TreeStore) {

	iter := gtk.TreeIter{}

	if isoTreeStore.GetIterFirst(&iter) {
		gVal := glib.GValue{}
		isoTreeStore.GetValue(&iter, 4, &gVal)
		getAndSetFieldVal(isoTreeStore, gVal.GetInt(), &iter)
		for isoTreeStore.IterNext(&iter) {
			gVal := glib.GValue{}
			isoTreeStore.GetValue(&iter, 4, &gVal)
			getAndSetFieldVal(isoTreeStore, gVal.GetInt(), &iter)

		}
	}

}

func getAndSetFieldVal(isoTreeStore *gtk.TreeStore, id int, iter *gtk.TreeIter) {
	fNameVal := glib.GValue{}
	isoTreeStore.GetValue(iter, pyui.COL_FIELD_NAME, &fNameVal)
	val := reqIsoMsg.GetFieldDataById(id)
	if len(val.Bytes()) > 0 {
		log.Printf("setting field [%s] value [%s]\n", fNameVal.GetString(), val.String())
		isoTreeStore.SetValue(iter, pyui.COL_FIELD_VAL, val.String())

		if val.Def() != nil {

			if val.Def().BitPosition() > 0 {
				if reqIsoMsg.IsSelected(val.Def().BitPosition()) {
					isoTreeStore.SetValue(iter, pyui.COL_FIELD_SELECTION, true)
				} else {
					isoTreeStore.SetValue(iter, pyui.COL_FIELD_SELECTION, false)
				}

			}
		} else {
			isoTreeStore.SetValue(iter, pyui.COL_FIELD_SELECTION, true)
		}

	}
}
