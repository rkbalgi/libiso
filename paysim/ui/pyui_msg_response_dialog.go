package ui

import (
	_ "github.com/mattn/go-gtk/gdk"
	_ "github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"github.com/rkbalgi/go/iso8583"
)

func ShowIsoResponseMsgDialog(resp_iso_msg *iso8583.Iso8583Message) {

	dialog := gtk.NewDialog()
	resp_data_view := gtk.NewTextView()
	//TODO:: temp solution
	
	resp_data_view.GetBuffer().SetText(resp_iso_msg.Dump())

	swin := gtk.NewScrolledWindow(nil, nil)
	swin.AddWithViewPort(resp_data_view)
   
    //vbox:=gtk.NewVBox(true,5);
    //vbox.PackStart(swin,true,true,5);

	dialog.GetVBox().Add(swin);
	dialog.SetSizeRequest(400,400);

	ok_btn := dialog.AddButton("OK", gtk.RESPONSE_OK)
	ok_btn.Connect("clicked", func() {

		dialog.Destroy()
		gtk.MainQuit()

	})

	dialog.SetPosition(gtk.WIN_POS_CENTER)

	dialog.ShowAll()
	gtk.Main()

}
