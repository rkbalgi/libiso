package ui

import (
	"container/list"
	_ "fmt"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"go/paysim"
)

func ShowIsoResponseMsgDialog(tabDataList *list.List) {

	dialog := gtk.NewDialog()
	dialog.SetTitle("Message Response")

	tView := gtk.NewTreeView()
	tView.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Field Name", gtk.NewCellRendererText(), "text", 0))
	tView.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Field Value", gtk.NewCellRendererText(), "text", 1))
	tStore := gtk.NewListStore(glib.G_TYPE_STRING, glib.G_TYPE_STRING)

	treeIter := gtk.TreeIter{}
	for l := tabDataList.Front(); l != nil; l = l.Next() {
		tStore.Append(&treeIter)
		tuple := l.Value.(*paysim.Tuple)
		tStore.SetValue(&treeIter, 0, tuple.Nth(0))
		tStore.SetValue(&treeIter, 1, tuple.Nth(1))
	}

	tView.GetColumn(0).SetClickable(true)
	tView.GetColumn(1).SetClickable(true)
	tView.GetColumn(0).SetExpand(true)
	tView.GetColumn(1).SetExpand(true)

	tView.SetModel(tStore)

	scrolledWindow := gtk.NewScrolledWindow(nil, nil)
	scrolledWindow.AddWithViewPort(tView)

	dialog.GetVBox().Add(scrolledWindow)
	dialog.SetSizeRequest(400, 300)

	okBtn := dialog.AddButton("OK", gtk.RESPONSE_OK)
	okBtn.Connect("clicked", func() {

		dialog.Destroy()
		gtk.MainQuit()

	})

	dialog.SetPosition(gtk.WIN_POS_CENTER)

	dialog.ShowAll()
	gtk.Main()

}
