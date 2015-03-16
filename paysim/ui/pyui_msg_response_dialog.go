package ui

import (
	"container/list"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"github.com/rkbalgi/go/paysim"
	"fmt"
)

func ShowIsoResponseMsgDialog(tab_data_list *list.List) {

	dialog := gtk.NewDialog()

	t_view := gtk.NewTreeView()
	t_view.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Field Name", gtk.NewCellRendererText(),"text",0))
	t_view.AppendColumn(gtk.NewTreeViewColumnWithAttributes("Field Value", gtk.NewCellRendererText(),"text",1))
	t_store := gtk.NewListStore(glib.G_TYPE_STRING, glib.G_TYPE_STRING)

	tree_iter := gtk.TreeIter{}
	for l := tab_data_list.Front(); l != nil; l = l.Next() {
		t_store.Append(&tree_iter)
		tuple := l.Value.(*paysim.Tuple)
		fmt.Println("???",tuple.String());
		t_store.SetValue(&tree_iter, 0, tuple.Nth(0))
		t_store.SetValue(&tree_iter, 1, tuple.Nth(1))
	}
	
	t_view.GetColumn(0).SetClickable(true);
	t_view.GetColumn(1).SetClickable(true)
	t_view.GetColumn(0).SetExpand(true);
	t_view.GetColumn(1).SetExpand(true)
	
	t_view.SetModel(t_store);

	swin := gtk.NewScrolledWindow(nil, nil)
	swin.AddWithViewPort(t_view)


	dialog.GetVBox().Add(swin)
	dialog.SetSizeRequest(400, 300)

	ok_btn := dialog.AddButton("OK", gtk.RESPONSE_OK)
	ok_btn.Connect("clicked", func() {

		dialog.Destroy()
		gtk.MainQuit()

	})

	dialog.SetPosition(gtk.WIN_POS_CENTER)

	dialog.ShowAll()
	gtk.Main()

}
