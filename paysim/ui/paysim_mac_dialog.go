package ui

import (
	"encoding/hex"
	_ "fmt"
	"github.com/mattn/go-gtk/gtk"
	"github.com/rkbalgi/go/crypto"
	"github.com/rkbalgi/go/crypto/mac"
)

func ComputeMacDialog(widget gtk.IWidget, msg string) {

	dialog := gtk.NewDialog()
	dialog.SetParent(widget)
	dialog.SetPosition(gtk.WIN_POS_CENTER)
	dialog.SetModal(true)
	dialog.SetTitle("Compute Mac")

	mac_entry := gtk.NewEntry()
	mac_entry.SetEditable(false)
	key_entry := gtk.NewEntry()
	key_entry.SetText("90897656d3e4de56")

	ok_btn := dialog.AddButton("Generate", gtk.RESPONSE_OK)
	cancel_btn := dialog.AddButton("Cancel", gtk.RESPONSE_CANCEL)

	mac_algo_cb := gtk.NewComboBoxText()
	mac_algo_cb.AppendText("X9.9")
	mac_algo_cb.AppendText("X9.19")

	iter := gtk.TreeIter{}
	mac_algo_cb.GetModel().GetIterFirst(&iter)
	mac_algo_cb.SetActiveIter(&iter)

	padding_type_cb := gtk.NewComboBoxText()
	padding_type_cb.AppendText("9797-1")
	padding_type_cb.AppendText("9792-2")

	iter = gtk.TreeIter{}
	padding_type_cb.GetModel().GetIterFirst(&iter)
	padding_type_cb.SetActiveIter(&iter)

	table := gtk.NewVBox(false, 5)

	hbox_tmp := gtk.NewHBox(false, 5)
	hbox_tmp.PackStart(new_label("Mac Key"), false, false, 5)
	hbox_tmp.PackStart(key_entry, true, true, 5)
	table.PackStart(hbox_tmp, false, false, 5)

	hbox_tmp = gtk.NewHBox(false, 5)
	hbox_tmp.PackStart(new_label("MAC Data"), false, false, 5)
	mac_buf_tv := gtk.NewTextView()
	mac_buf_tv.SetWrapMode(gtk.WRAP_CHAR)
	mac_buf_tv.GetBuffer().SetText("0000000000000000")
	swin := gtk.NewScrolledWindow(nil, nil)
	swin.AddWithViewPort(mac_buf_tv)
	swin.SetSizeRequest(180, 150)
	hbox_tmp.Add(swin)
	table.Add(hbox_tmp)

	hbox_tmp = gtk.NewHBox(false, 5)
	hbox_tmp.PackStart(new_label("MAC Algorithm"), false, false, 5)
	hbox_tmp.PackStart(mac_algo_cb, false, false, 5)
	table.PackStart(hbox_tmp, false, false, 5)

	hbox_tmp = gtk.NewHBox(false, 5)
	hbox_tmp.PackStart(new_label("Padding Type"), false, false, 5)
	hbox_tmp.PackStart(padding_type_cb, false, false, 5)
	table.PackStart(hbox_tmp, false, false, 5)

	hbox_tmp = gtk.NewHBox(false, 5)
	hbox_tmp.PackStart(new_label("MAC"), false, false, 5)
	hbox_tmp.PackStart(mac_entry, false, false, 10)
	table.PackStart(hbox_tmp, false, false, 5)

	dialog.GetVBox().Add(table)

	ok_btn.Connect("clicked", func() {

		key_data, err := hex.DecodeString(key_entry.GetText())
		if err != nil {
			ShowErrorDialog(dialog, "Invalid Key!")
			return
		}

		var start_iter, end_iter gtk.TextIter
		mac_buf_tv.GetBuffer().GetStartIter(&start_iter)
		mac_buf_tv.GetBuffer().GetEndIter(&end_iter)
		mac_data_str := mac_buf_tv.GetBuffer().GetText(&start_iter, &end_iter, true)

		mac_data, err := hex.DecodeString(mac_data_str)
		if err != nil {
			ShowErrorDialog(dialog, "Invalid Data (should be hex)")
			return
		}

		println("supplied key length -", len(key_data))
		if mac_algo_cb.GetActiveText() == "X9.19" {
			if len(key_data) != 16 {
				ShowErrorDialog(dialog, "Invalid key. A double length DES key expected.")
				return
			}
		} else {
			//X9.9
			if len(key_data) != 8 {
				ShowErrorDialog(dialog, "Invalid key. A single length DES key expected.")
				return
			}
		}

		var padded_data []byte
		if padding_type_cb.GetActiveText() == "9797-1" {
			padded_data = crypto.Iso9797M1Padding.Pad(mac_data)
		} else {
			padded_data = crypto.Iso9797M2Padding.Pad(mac_data)
		}

		var mac_val []byte
		if len(key_data) == 8 {
			mac_val = mac.GenerateMac_X99(padded_data, key_data)
		} else {
			mac_val = mac.GenerateMac_X919(padded_data, key_data)
		}
		mac_entry.SetText(hex.EncodeToString(mac_val))

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

func new_label(label_txt string) *gtk.Label {

	label := gtk.NewLabel(label_txt)
	label.SetSizeRequest(100, 5)

	return label
}
