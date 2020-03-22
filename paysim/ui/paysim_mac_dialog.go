package ui

import (
	"encoding/hex"
	_ "fmt"
	"github.com/mattn/go-gtk/gtk"
	"go/crypto"
	"go/crypto/mac"
)

func ComputeMacDialog(widget gtk.IWidget, msg string) {

	dialog := gtk.NewDialog()
	dialog.SetParent(widget)
	dialog.SetPosition(gtk.WIN_POS_CENTER)
	dialog.SetModal(true)
	dialog.SetTitle("Compute Mac")

	macEntry := gtk.NewEntry()
	macEntry.SetEditable(false)
	keyEntry := gtk.NewEntry()
	keyEntry.SetText("90897656d3e4de56")

	okBtn := dialog.AddButton("Generate", gtk.RESPONSE_OK)
	cancelBtn := dialog.AddButton("Cancel", gtk.RESPONSE_CANCEL)

	macAlgoCb := gtk.NewComboBoxText()
	macAlgoCb.AppendText("X9.9")
	macAlgoCb.AppendText("X9.19")

	iter := gtk.TreeIter{}
	macAlgoCb.GetModel().GetIterFirst(&iter)
	macAlgoCb.SetActiveIter(&iter)

	paddingTypeCb := gtk.NewComboBoxText()
	paddingTypeCb.AppendText("9797-1")
	paddingTypeCb.AppendText("9792-2")

	iter = gtk.TreeIter{}
	paddingTypeCb.GetModel().GetIterFirst(&iter)
	paddingTypeCb.SetActiveIter(&iter)

	table := gtk.NewVBox(false, 5)

	hboxTmp := gtk.NewHBox(false, 5)
	hboxTmp.PackStart(newLabel("Mac Key"), false, false, 5)
	hboxTmp.PackStart(keyEntry, true, true, 5)
	table.PackStart(hboxTmp, false, false, 5)

	hboxTmp = gtk.NewHBox(false, 5)
	hboxTmp.PackStart(newLabel("MAC Data"), false, false, 5)
	macBufTv := gtk.NewTextView()
	macBufTv.SetWrapMode(gtk.WRAP_CHAR)
	macBufTv.GetBuffer().SetText("0000000000000000")
	scrolledWindow := gtk.NewScrolledWindow(nil, nil)
	scrolledWindow.AddWithViewPort(macBufTv)
	scrolledWindow.SetSizeRequest(180, 150)
	hboxTmp.Add(scrolledWindow)
	table.Add(hboxTmp)

	hboxTmp = gtk.NewHBox(false, 5)
	hboxTmp.PackStart(newLabel("MAC Algorithm"), false, false, 5)
	hboxTmp.PackStart(macAlgoCb, false, false, 5)
	table.PackStart(hboxTmp, false, false, 5)

	hboxTmp = gtk.NewHBox(false, 5)
	hboxTmp.PackStart(newLabel("Padding Type"), false, false, 5)
	hboxTmp.PackStart(paddingTypeCb, false, false, 5)
	table.PackStart(hboxTmp, false, false, 5)

	hboxTmp = gtk.NewHBox(false, 5)
	hboxTmp.PackStart(newLabel("MAC"), false, false, 5)
	hboxTmp.PackStart(macEntry, false, false, 10)
	table.PackStart(hboxTmp, false, false, 5)

	dialog.GetVBox().Add(table)

	okBtn.Connect("clicked", func() {

		keyData, err := hex.DecodeString(keyEntry.GetText())
		if err != nil {
			ShowErrorDialog(dialog, "Invalid Key!")
			return
		}

		var startIter, endIter gtk.TextIter
		macBufTv.GetBuffer().GetStartIter(&startIter)
		macBufTv.GetBuffer().GetEndIter(&endIter)
		macDataStr := macBufTv.GetBuffer().GetText(&startIter, &endIter, true)

		macData, err := hex.DecodeString(macDataStr)
		if err != nil {
			ShowErrorDialog(dialog, "Invalid Data (should be hex)")
			return
		}

		println("supplied key length -", len(keyData))
		if macAlgoCb.GetActiveText() == "X9.19" {
			if len(keyData) != 16 {
				ShowErrorDialog(dialog, "Invalid key. A double length DES key expected.")
				return
			}
		} else {
			//X9.9
			if len(keyData) != 8 {
				ShowErrorDialog(dialog, "Invalid key. A single length DES key expected.")
				return
			}
		}

		var paddedData []byte
		if paddingTypeCb.GetActiveText() == "9797-1" {
			paddedData = crypto.Iso9797M1Padding.Pad(macData)
		} else {
			paddedData = crypto.Iso9797M2Padding.Pad(macData)
		}

		var macVal []byte
		if len(keyData) == 8 {
			macVal, _ = mac.GenerateMacX99(paddedData, keyData)
		} else {
			macVal, _ = mac.GenerateMacX919(paddedData, keyData)
		}
		macEntry.SetText(hex.EncodeToString(macVal))

	})

	cancelBtn.Connect("clicked", func() {

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

func newLabel(labelTxt string) *gtk.Label {

	label := gtk.NewLabel(labelTxt)
	label.SetSizeRequest(100, 5)

	return label
}
