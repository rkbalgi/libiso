package ui

//This file contains all of the necessary types and functions
//to support the tree on the right hand side of the application that
//controls the spec message layout

import (
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

//columns in the SpecMsgTree
const (
	COL_FIELD_SELECTION     int = 0
	COL_BIT_POS                 = 1
	COL_FIELD_NAME              = 2
	COL_FIELD_VAL               = 3
	COL_FIELD_ID                = 4
	COLD_FIELD_EDIT_ENABLED     = 5
)

type PaysimSpecMsgTree struct {
	_store *gtk.TreeStore
	_view  *gtk.TreeView

	_val_rend    *gtk.CellRendererText
	_toggle_rend *gtk.CellRendererToggle
}

func NewPaysimSpecMsgTree() *PaysimSpecMsgTree {

	self := new(PaysimSpecMsgTree)

	field_name_renderer := gtk.NewCellRendererText()
	self._val_rend = gtk.NewCellRendererText()
	self._toggle_rend = gtk.NewCellRendererToggle()

	self._toggle_rend.SetActivatable(true)
	//  0             1             2         3            4        5
	//is selected, bit position,field name, field value, field id,editable
	self._store = gtk.NewTreeStore(glib.G_TYPE_BOOL, glib.G_TYPE_INT, glib.G_TYPE_STRING, glib.G_TYPE_STRING, glib.G_TYPE_INT, glib.G_TYPE_BOOL)
	self._view = gtk.NewTreeView()

	self._view.AppendColumn(gtk.NewTreeViewColumnWithAttribute("Is Selected", self._toggle_rend))
	self._view.AppendColumn(gtk.NewTreeViewColumnWithAttribute("Bit Position", field_name_renderer))
	self._view.AppendColumn(gtk.NewTreeViewColumnWithAttribute("Field Name", field_name_renderer))
	self._view.AppendColumn(gtk.NewTreeViewColumnWithAttribute("Data", self._val_rend))

	self._view.GetColumn(0).AddAttribute(self._toggle_rend, "active", 0)
	self._view.GetColumn(1).AddAttribute(field_name_renderer, "text", 1)
	self._view.GetColumn(2).AddAttribute(field_name_renderer, "text", 2)
	self._view.GetColumn(3).AddAttribute(self._val_rend, "text", 3)
	self._view.GetColumn(3).AddAttribute(self._val_rend, "editable", 5)

	for i := 0; i < 4; i++ {
		self._view.GetColumn(i).SetProperty("resizable", glib.ValueFromNative(true))
	}

    self._view.ModifyFontEasy("Consolas 10");
    
	self._view.SetModel(self._store)

	return self

}

func (self *PaysimSpecMsgTree) Store() *gtk.TreeStore {
	return self._store
}

func (self *PaysimSpecMsgTree) View() *gtk.TreeView {
	return self._view
}

func (self *PaysimSpecMsgTree) FieldValueRenderer() *gtk.CellRendererText {
	return self._val_rend
}

func (self *PaysimSpecMsgTree) FieldToggleRenderer() *gtk.CellRendererToggle {
	return self._toggle_rend
}
