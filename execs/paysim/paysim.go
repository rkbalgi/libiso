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
	"unsafe"
)

var right_box, spec_tree_vbox *gtk.VBox
var status_bar *gtk.Statusbar
var status_bar_context_id uint
var active_spec_frame *gtk.Frame
var swin_spec_tree *gtk.ScrolledWindow
var ui_ctx *pyui.PaysimUiContext
var bmp_iter gtk.TreeIter
var h_box *gtk.HBox
var swin_console *gtk.ScrolledWindow

//the tree holding all the specs
var spec_tree_view *gtk.TreeView

//the currently loaded iso_msg
var req_iso_msg, resp_iso_msg *iso8583.Iso8583Message
var comms_config_bx *gtk.VBox

func main() {

	gtk.Init(nil)
	gtk.SettingsGetDefault().SetStringProperty("gtk-font-name", "Dejavu Sans 8", "")
	//a vbox that will hold the spec tree
	//and its contents
	pyui.Init()
	spec_tree_vbox = gtk.NewVBox(false, 1)

	ui_ctx = pyui.NewUiContext()

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
			ui_ctx.Window(),
			gtk.FILE_CHOOSER_ACTION_OPEN,
			gtk.STOCK_OK,
			gtk.RESPONSE_ACCEPT)
		filter := gtk.NewFileFilter()
		filter.AddPattern("*.json")
		filechooserdialog.AddFilter(filter)
		filechooserdialog.Response(func() {
			iso8583.ReadSpecDefs(filechooserdialog.GetFilename())
			make_and_populate_spec_tree()
			filechooserdialog.Destroy() 
		})
		filechooserdialog.Run()

	})

	quit_mi := gtk.NewMenuItemWithMnemonic("E_xit")
	quit_mi.Connect("activate", func() {
		gtk.MainQuit()
	})
	file_menu.Add(quit_mi)

	cascade_utils_mi := gtk.NewMenuItemWithMnemonic("_Utils")
	utils_menu := gtk.NewMenu()
	cascade_utils_mi.SetSubmenu(utils_menu)

	mac_mi := gtk.NewMenuItemWithMnemonic("_MAC")
	utils_menu.Add(mac_mi)
	mac_mi.Connect("activate", func() {
		pyui.ComputeMacDialog(ui_ctx.Window(), "")
	})

	pin_mi := gtk.NewMenuItemWithMnemonic("_PIN")
	utils_menu.Add(pin_mi)
	pin_mi.Connect("activate", func() {})

	menu_bar.Append(cascade_utils_mi)

	cascade_about_mi := gtk.NewMenuItemWithMnemonic("_About")

	cascade_about_mi.Connect("activate", func() {
		dialog := gtk.NewAboutDialog()
		dialog.SetName("About Paysim")
		dialog.SetProgramName("PaySim v1.0 build 03052015	")
		dialog.SetAuthors([]string{"Raghavendra Balgi (rkbalgi@gmail.com)"})

		dialog.SetLicense("This application is available under the same terms and conditions as the Go, the BSD style license, and the LGPL (Lesser GNU Public License).")
		dialog.SetWrapLicense(true)
		dialog.Run()
		dialog.Destroy()
	})

	menu_bar.Append(cascade_about_mi)

	vbox := gtk.NewVBox(false, 1)
	vbox.PackStart(menu_bar, false, false, 0)

	//add a vertical pane
	h_pane := gtk.NewHPaned()
	h_pane.SetPosition(40)
	vbox.Add(h_pane)

	frame1 := gtk.NewFrame("")
	frame1.SetSizeRequest(100, 600)
	frame1.Add(spec_tree_vbox)

	frame2 := gtk.NewFrame("")
	right_box = gtk.NewVBox(false, 1)
	frame2.Add(right_box)
	h_pane.Pack1(frame1, false, false)
	h_pane.Pack2(frame2, false, false)

	ui_ctx.Window().Add(vbox)

	status_bar = gtk.NewStatusbar()
	status_bar_context_id = status_bar.GetContextId("pay-sim")
	status_bar.Push(status_bar_context_id, "OK")
	vbox.PackStart(status_bar, false, false, 1)

	right_box.PackStart(ui_ctx.CommsConfigVBox(), false, false, 12)

	reset_console()
	active_spec_frame = gtk.NewFrame("                 ")
	active_spec_frame.SetName("spec_frame")
	active_spec_frame.SetSizeRequest(400, 400)

	right_box.PackStart(active_spec_frame, false, false, 5)
	right_box.PackStart(swin_console, false, false, 5)
	right_box.ShowAll()

	pylog.Log("Paysim v1.00 starting...\n#********************************#")

	iso8583.ReadDemoSpecDefs()
	make_and_populate_spec_tree()

	ui_ctx.Window().ShowAll()
	gtk.Main()

}

func make_and_populate_spec_tree() {
	//remove the old tree if there is already one
	//and then create a new one with the specs loaded in the
	//tree
	if spec_tree_view != nil {
		spec_tree_vbox.Remove(swin_spec_tree)
	}

	tree_store := gtk.NewTreeStore(glib.G_TYPE_STRING)
	spec_tree_view = gtk.NewTreeView()
	spec_tree_view.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Specs", gtk.NewCellRendererText(), "text", 0))
	spec_tree_view.SetModel(tree_store.ToTreeModel())

	spec_tree_view.SetSizeRequest(250, 800)

	var ti_1, ti_2 gtk.TreeIter
	tree_store.Append(&ti_1, nil)
	tree_store.Set(&ti_1, "Specifications")
	n := 0
	for _, spec := range iso8583.GetSpecs() {
		tree_store.Append(&ti_2, &ti_1)
		tree_store.Set(&ti_2, spec.Name())
		n++
	}

	spec_tree_view.SetHeadersVisible(false)

	spec_tree_view.Connect("button-release-event", func(ctx *glib.CallbackContext) {

		event := (*gdk.EventButton)(unsafe.Pointer(ctx.Args(0)))
		if event.Button == 1 {

			var path *gtk.TreePath
			var col *gtk.TreeViewColumn
			spec_tree_view.GetCursor(&path, &col)

			if path.String() == "0" {
				return
			}

			val := new(glib.GValue)
			var t_iter gtk.TreeIter
			if tree_store.GetIter(&t_iter, path) {
				tree_store.GetValue(&t_iter, 0, val)
				show_spec_layout(val.GetString())
			}

		}

	})

	status_bar.Push(status_bar_context_id, fmt.Sprintf("%d specs loaded. OK.", n))

	swin_spec_tree = gtk.NewScrolledWindow(nil, nil)
	swin_spec_tree.SetSizeRequest(200, 800)
	swin_spec_tree.AddWithViewPort(spec_tree_view)

	spec_tree_vbox.PackStart(swin_spec_tree, true, true, 2)

	spec_tree_vbox.ShowAll()

}

func show_spec_layout(spec_name string) {

	//this 'req_iso_msg' will act as the model
	//that backs the TreeView
	req_iso_msg = iso8583.NewIso8583Message(spec_name)

	//remove the active frame
	if active_spec_frame != nil {
		right_box.Remove(active_spec_frame)
		right_box.Remove(h_box)
	}
	//right_box.Remove(swin_console)

	spec_msg_tree := pyui.NewPaysimSpecMsgTree()

	spec_lyt_tree := spec_msg_tree.View()
	spec_lyt_store := spec_msg_tree.Store()

	var i1 gtk.TreeIter 

	spec_lyt := iso8583.GetSpecLayout(spec_name)
	for _, field_def := range spec_lyt {
		spec_lyt_store.Append(&i1, nil)
		spec_lyt_store.SetValue(&i1, 0, false)
		spec_lyt_store.SetValue(&i1, 1, field_def.BitPosition)
		spec_lyt_store.SetValue(&i1, 2, field_def.Name)
		spec_lyt_store.SetValue(&i1, 3, "") 
		spec_lyt_store.SetValue(&i1, 4, field_def.Id)
		if field_def.Name == "Bitmap" || field_def.Name == "Message Type" {
			spec_lyt_store.SetValue(&i1, 0, true)
		} else {
			spec_lyt_store.SetValue(&i1, 0, false)
		}

		//bitmap will be non-editable
		if field_def.Name == "Bitmap" {
			//lets store the iter
			bmp_iter = i1
			spec_lyt_store.SetValue(&i1, 3, "00000000000000000000000000000000")
			spec_lyt_store.SetValue(&i1, 5, false)
		} else {
			spec_lyt_store.SetValue(&i1, 5, true)
		}

	}

	spec_msg_tree.FieldValueRenderer().Connect("edited",
		func(ctx *glib.CallbackContext) {

			field_name := get_current_field_name(spec_lyt_tree, spec_lyt_store)
			if field_name == "Bitmap" {
				//do not allow edits on Bitmap field
				return
			}

			set_current_field_value(spec_lyt_tree, spec_lyt_store, ctx.Args(1).ToString())

		})

	spec_msg_tree.FieldToggleRenderer().Connect("toggled",
		func(ctx *glib.CallbackContext) {

			var path *gtk.TreePath
			var col *gtk.TreeViewColumn
			spec_lyt_tree.GetCursor(&path, &col)
			log.Println(col.GetTitle())
			if col.GetTitle() != "Is Selected" {
				return
			}

			var i1 gtk.TreeIter
			if spec_lyt_store.GetIter(&i1, path) {

				field_name_val := glib.GValue{}
				spec_lyt_store.GetValue(&i1, 2, &field_name_val)
				fmt.Println(field_name_val.GetString())

				if field_name_val.GetString() == "Bitmap" || field_name_val.GetString() == "Message Type" {
					return
				}

				val := glib.GValue{}
				bit_pos_val := glib.GValue{}
				spec_lyt_store.GetValue(&i1, 0, &val)
				spec_lyt_store.GetValue(&i1, 1, &bit_pos_val)

				spec_lyt_store.SetValue(&i1, 0, !val.GetBool())

				if !val.GetBool() {
					req_iso_msg.Bitmap().SetOn(bit_pos_val.GetInt())
				} else {
					req_iso_msg.Bitmap().SetOff(bit_pos_val.GetInt())
				}

				recompute_bitmap(spec_lyt_store)
			}

		})

	ui_ctx.LoadButton().Connect("clicked", func() {

		trace_data, err := ui_ctx.GetUsrTrace()
		log.Println("user trace: ", hex.EncodeToString(trace_data))
		if err == nil {
			msg_buf := bytes.NewBuffer(trace_data)
			err = req_iso_msg.Parse(msg_buf)
			if err != nil {
				pyui.ShowErrorDialog(ui_ctx.Window(), "Trace Parse Error")
			}

			populate_model(spec_lyt_store)
		}

	})

	ui_ctx.AssembleButton().Connect("clicked", func() {

		trace_data := req_iso_msg.Bytes()
		ui_ctx.ShowUsrTrace(trace_data)
	})
	ui_ctx.SendButton().Connect("clicked", func() {

		tcp_addr, mli_type_str, err := ui_ctx.GetCommsConfig()
		if err != nil {
			pyui.ShowErrorDialog(ui_ctx.Window(), err.Error())
			return
		}

		//send this msg, get the response and then display
		//the response
		resp_iso_msg, err = pynet.SendIsoMsg(tcp_addr.String(), mli_type_str, req_iso_msg)
		if err != nil {
			pyui.ShowErrorDialog(ui_ctx.Window(), err.Error())
		} else {
			//hoo-hah! response received
			//display it as a dialog to the user
			pyui.ShowIsoResponseMsgDialog(resp_iso_msg.TabularFormat())

		}

	})

	align := gtk.NewAlignment(0.5, 0.5, 0.0, 0.0)
	align.Add(ui_ctx.ButtonBox())
	spec_lyt_tree.ShowAll()
	active_spec_frame = gtk.NewFrame("       [" + spec_name + "]        ")
	active_spec_frame.SetName("spec_frame")
	active_spec_frame.SetSizeRequest(300, 300)

	//add a scrolled winow containing
	//the spec_lty_tree
	swin := gtk.NewScrolledWindow(nil, nil)
	swin.AddWithViewPort(spec_lyt_tree)
	hbox_tmp := gtk.NewHBox(false, 10)
	hbox_tmp.PackStart(swin, true, true, 10)

	active_spec_frame.Add(hbox_tmp)

	right_box.SetSizeRequest(400, 200)
	right_box.PackStart(active_spec_frame, false, false, 20)
	right_box.PackStart(align, false, false, 0)
	//reset_console()
	right_box.ReorderChild(swin_console, -1) //.PackStart(swin_console, false, false, 0)

	right_box.ShowAll()

}

func reset_console() {

	swin_console = gtk.NewScrolledWindow(nil, nil)
	swin_console.AddWithViewPort(pyui.PaysimConsole.TextView())
	swin_console.SetSizeRequest(400, 150)

}

func get_current_field_name(iso_msg_tree *gtk.TreeView,
	iso_msg_tree_store *gtk.TreeStore) string {

	var path *gtk.TreePath
	var col *gtk.TreeViewColumn
	iso_msg_tree.GetCursor(&path, &col)
	var i1 gtk.TreeIter

	if iso_msg_tree_store.GetIter(&i1, path) {
		field_name_val := glib.GValue{}
		iso_msg_tree_store.GetValue(&i1, 2, &field_name_val)
		return field_name_val.GetString()
	}

	return ""

}

func set_current_field_value(iso_msg_tree *gtk.TreeView,
	iso_msg_tree_store *gtk.TreeStore, val string) {

	var path *gtk.TreePath
	var col *gtk.TreeViewColumn
	iso_msg_tree.GetCursor(&path, &col)
	var i1 gtk.TreeIter

	if iso_msg_tree_store.GetIter(&i1, path) {

		f_name := glib.GValue{}
		iso_msg_tree_store.GetValue(&i1, 2, &f_name)
		log.Println("setting field to value ", f_name.GetString(), "to ", val)
		fld_id := glib.GValue{}
		iso_msg_tree_store.GetValue(&i1, 4, &fld_id)
		req_iso_msg.SetFieldData(fld_id.GetInt(), val)
		//data might have been truncated or padded based on the
		//field definition
		new_val := req_iso_msg.GetFieldDataById(fld_id.GetInt()).String()
		iso_msg_tree_store.SetValue(&i1, 3, new_val)

	}

}

//recompute_bitmap will show a new value of the bitmap
//on the UI
func recompute_bitmap(iso_tree_store *gtk.TreeStore) {

	val := glib.GValue{}
	bmp_data := req_iso_msg.Bitmap().Bytes()
	iso_tree_store.GetValue(&bmp_iter, 3, &val)
	log.Println("Bitmap's current value", val.GetString())

	iso_tree_store.SetValue(&bmp_iter, 3, hex.EncodeToString(bmp_data))

}

//populates the store based on the data in
//req_iso_msg
func populate_model(iso_tree_store *gtk.TreeStore) {

	iter := gtk.TreeIter{}

	if iso_tree_store.GetIterFirst(&iter) {
		g_val := glib.GValue{}
		iso_tree_store.GetValue(&iter, 4, &g_val)
		get_and_set_field_val(iso_tree_store, g_val.GetInt(), &iter)
		for iso_tree_store.IterNext(&iter) {
			g_val := glib.GValue{}
			iso_tree_store.GetValue(&iter, 4, &g_val)
			get_and_set_field_val(iso_tree_store, g_val.GetInt(), &iter)

		}
	}

}

func get_and_set_field_val(iso_tree_store *gtk.TreeStore, id int, iter *gtk.TreeIter) {
	f_name_val := glib.GValue{}
	iso_tree_store.GetValue(iter, pyui.COL_FIELD_NAME, &f_name_val)
	val := req_iso_msg.GetFieldDataById(id)
	if len(val.Bytes()) > 0 {
		log.Printf("setting field [%s] value [%s]\n", f_name_val.GetString(), val.String())
		iso_tree_store.SetValue(iter, pyui.COL_FIELD_VAL, val.String())

		if val.Def() != nil {

			if val.Def().BitPosition() > 0 {
				if req_iso_msg.IsSelected(val.Def().BitPosition()) {
					iso_tree_store.SetValue(iter, pyui.COL_FIELD_SELECTION, true)
				} else {
					iso_tree_store.SetValue(iter, pyui.COL_FIELD_SELECTION, false)
				}

			}
		} else {
			iso_tree_store.SetValue(iter, pyui.COL_FIELD_SELECTION, true)
		}

	}
}
