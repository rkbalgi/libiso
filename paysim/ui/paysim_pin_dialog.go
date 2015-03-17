package ui

import (
	"encoding/hex"
	_ "fmt"
	"github.com/mattn/go-gtk/gtk"
	"github.com/rkbalgi/go/crypto/pin"
)

func ComputePinBlockDialog(widget gtk.IWidget, msg string) {

	dialog := gtk.NewDialog()
	dialog.SetParent(widget)
	dialog.SetPosition(gtk.WIN_POS_CENTER)
	dialog.SetModal(true)
	dialog.SetTitle("Compute Pin Block (ISO DF52)")

	pin_block_entry := gtk.NewEntry()
	pin_block_entry.SetEditable(false)
	key_entry := gtk.NewEntry()
	key_entry.SetText("90897656d3e4de568967f4deee3d56f2")
	clear_pin_entry := gtk.NewEntry()
	clear_pin_entry.SetText("1234")

	pan_entry := gtk.NewEntry()
	pan_entry.SetText("5400000000000006")

	ok_btn := dialog.AddButton("Generate", gtk.BUTTONS_OK)
	cancel_btn := dialog.AddButton("Cancel", gtk.BUTTONS_CANCEL)

	pin_algo_cb := gtk.NewComboBoxText()
	pin_algo_cb.AppendText("IBM-3264")
	pin_algo_cb.AppendText("ISO-0")
	pin_algo_cb.AppendText("ISO-1")
	pin_algo_cb.AppendText("ISO-3")

	iter := gtk.TreeIter{}
	pin_algo_cb.GetModel().GetIterFirst(&iter)
	pin_algo_cb.SetActiveIter(&iter)

	table := gtk.NewVBox(false, 5)

	hbox_tmp := gtk.NewHBox(false, 5)
	hbox_tmp.PackStart(new_label("PIN Key"), false, false, 5)
	hbox_tmp.PackStart(key_entry, true, true, 5)
	table.PackStart(hbox_tmp, false, false, 5)

	hbox_tmp = gtk.NewHBox(false, 5)
	hbox_tmp.PackStart(new_label("PAN"), false, false, 5)
	hbox_tmp.PackStart(pan_entry, false, false, 5)
	table.Add(hbox_tmp)

	hbox_tmp = gtk.NewHBox(false, 5)
	hbox_tmp.PackStart(new_label("PIN Algorithm"), false, false, 5)
	hbox_tmp.PackStart(pin_algo_cb, false, false, 5)
	table.PackStart(hbox_tmp, false, false, 5)

	hbox_tmp = gtk.NewHBox(false, 5)
	hbox_tmp.PackStart(new_label("Clear PIN"), false, false, 5)
	hbox_tmp.PackStart(clear_pin_entry, false, false, 10)
	table.PackStart(hbox_tmp, false, false, 5)

	hbox_tmp = gtk.NewHBox(false, 5)
	hbox_tmp.PackStart(new_label("PIN Block"), false, false, 5)
	hbox_tmp.PackStart(pin_block_entry, false, false, 10)
	table.PackStart(hbox_tmp, false, false, 5)

	dialog.GetVBox().Add(table)

	ok_btn.Connect("clicked", func() {

		key_data, err := hex.DecodeString(key_entry.GetText())
		if err != nil {
			ShowErrorDialog(dialog, "Invalid Key!")
			return
		}

		if len(key_data) == 16 {
		} else if len(key_data) == 8 {
		} else {
			ShowErrorDialog(dialog, "Invalid Key. Please enter single/double length DES key.")
			return
		}

		pan := pan_entry.GetText()
		if len(pan) == 0 {
			ShowErrorDialog(dialog, "Invalid PAN.")
			return
		}
		
		//get only rightmost 12
		pan=pan[len(pan)-1-12:len(pan)-1]

		clear_pin := clear_pin_entry.GetText()
		if len(clear_pin) < 4 || len(clear_pin) > 12 {
			ShowErrorDialog(dialog, "invalid clear pin (should be between 4 and 12 digits).")
			return
		}

		var pin_block []byte
		//var err error;

		switch pin_algo_cb.GetActiveText() {
		case "IBM-3264":
			{
				pin_block_type := &pin.PinBlock_Ibm3264{}
				pin_block = pin_block_type.Encrypt(pan, clear_pin, key_data)
			}
		case "ISO-0":
			{
				pin_block_type := &pin.PinBlock_Iso0{};
				pin_block = pin_block_type.Encrypt(pan, clear_pin, key_data)
			}
		case "ISO-1":
			{
				pin_block_type := &pin.PinBlock_Iso1{};
				pin_block = pin_block_type.Encrypt(pan, clear_pin, key_data)
			}	
		case "ISO-3":
			{
				pin_block_type := &pin.PinBlock_Iso3{};
				pin_block = pin_block_type.Encrypt(pan, clear_pin, key_data)
			}							
		default:
			{
				pin_block = make([]byte, 8)
			}
		}

		if err != nil {
			ShowErrorDialog(dialog, "error handling pin block ["+err.Error()+"]")
			return
		}

		pin_block_entry.SetText(hex.EncodeToString(pin_block))

	})

	cancel_btn.Connect("clicked", func() {

		dialog.Destroy()
		gtk.MainQuit()

	})

	dialog.SetResizable(false)
	//table.ShowAll()
	dialog.SetSizeRequest(400, 250)
	dialog.GetVBox().ShowAll()
	dialog.ShowAll()
	gtk.Main()

}
