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

	pinBlockEntry := gtk.NewEntry()
	pinBlockEntry.SetEditable(false)
	keyEntry := gtk.NewEntry()
	keyEntry.SetText("90897656d3e4de568967f4deee3d56f2")
	clearPinEntry := gtk.NewEntry()
	clearPinEntry.SetText("1234")

	panEntry := gtk.NewEntry()
	panEntry.SetText("5400000000000006")

	okBtn := dialog.AddButton("Generate", gtk.RESPONSE_OK)
	cancelBtn := dialog.AddButton("Cancel", gtk.RESPONSE_CANCEL)

	pinAlgoCb := gtk.NewComboBoxText()
	pinAlgoCb.AppendText("IBM-3264")
	pinAlgoCb.AppendText("ISO-0")
	pinAlgoCb.AppendText("ISO-1")
	pinAlgoCb.AppendText("ISO-3")

	iter := gtk.TreeIter{}
	pinAlgoCb.GetModel().GetIterFirst(&iter)
	pinAlgoCb.SetActiveIter(&iter)

	table := gtk.NewVBox(false, 5)

	hboxTmp := gtk.NewHBox(false, 5)
	hboxTmp.PackStart(newLabel("PIN Key"), false, false, 5)
	hboxTmp.PackStart(keyEntry, true, true, 5)
	table.PackStart(hboxTmp, false, false, 5)

	hboxTmp = gtk.NewHBox(false, 5)
	hboxTmp.PackStart(newLabel("PAN"), false, false, 5)
	hboxTmp.PackStart(panEntry, false, false, 5)
	table.Add(hboxTmp)

	hboxTmp = gtk.NewHBox(false, 5)
	hboxTmp.PackStart(newLabel("PIN Algorithm"), false, false, 5)
	hboxTmp.PackStart(pinAlgoCb, false, false, 5)
	table.PackStart(hboxTmp, false, false, 5)

	hboxTmp = gtk.NewHBox(false, 5)
	hboxTmp.PackStart(newLabel("Clear PIN"), false, false, 5)
	hboxTmp.PackStart(clearPinEntry, false, false, 10)
	table.PackStart(hboxTmp, false, false, 5)

	hboxTmp = gtk.NewHBox(false, 5)
	hboxTmp.PackStart(newLabel("PIN Block"), false, false, 5)
	hboxTmp.PackStart(pinBlockEntry, false, false, 10)
	table.PackStart(hboxTmp, false, false, 5)

	dialog.GetVBox().Add(table)

	okBtn.Connect("clicked", func() {

		keyData, err := hex.DecodeString(keyEntry.GetText())
		if err != nil {
			ShowErrorDialog(dialog, "Invalid Key!")
			return
		}

		if len(keyData) == 16 {
		} else if len(keyData) == 8 {
		} else {
			ShowErrorDialog(dialog, "Invalid Key. Please enter single/double length DES key.")
			return
		}

		pan := panEntry.GetText()
		if len(pan) == 0 {
			ShowErrorDialog(dialog, "Invalid PAN.")
			return
		}

		//get only rightmost 12
		pan = pan[len(pan)-1-12 : len(pan)-1]

		clearPin := clearPinEntry.GetText()
		if len(clearPin) < 4 || len(clearPin) > 12 {
			ShowErrorDialog(dialog, "invalid clear pin (should be between 4 and 12 digits).")
			return
		}

		var pinBlock []byte
		//var err error;

		switch pinAlgoCb.GetActiveText() {
		case "IBM-3264":
			{
				pinBlockType := &pin.PinblockIbm3264{}
				pinBlock, err = pinBlockType.Encrypt(pan, clearPin, keyData)
				if err != nil {
					ShowErrorDialog(dialog, "System error - "+err.Error())
					return
				}
			}
		case "ISO-0":
			{
				pinBlockType := &pin.PinBlock_Iso0{}
				pinBlock, err = pinBlockType.Encrypt(pan, clearPin, keyData)
				if err != nil {
					ShowErrorDialog(dialog, "System error - "+err.Error())
					return
				}
			}
		case "ISO-1":
			{
				pinBlockType := &pin.PinblockIso1{}
				pinBlock, err = pinBlockType.Encrypt(pan, clearPin, keyData)
				if err != nil {
					ShowErrorDialog(dialog, "System error - "+err.Error())
					return
				}
			}
		case "ISO-3":
			{
				pinBlockType := &pin.PinblockIso3{}
				pinBlock, err = pinBlockType.Encrypt(pan, clearPin, keyData)
				if err != nil {
					ShowErrorDialog(dialog, "System error - "+err.Error())
					return
				}
			}
		default:
			{
				pinBlock = make([]byte, 8)
			}
		}

		pinBlockEntry.SetText(hex.EncodeToString(pinBlock))

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
