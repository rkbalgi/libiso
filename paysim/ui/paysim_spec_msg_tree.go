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

	valRend    *gtk.CellRendererText
	toggleRend *gtk.CellRendererToggle
}

func NewPaysimSpecMsgTree() *PaysimSpecMsgTree {

	self := new(PaysimSpecMsgTree)

	fieldNameRenderer := gtk.NewCellRendererText()
	self.valRend = gtk.NewCellRendererText()
	self.toggleRend = gtk.NewCellRendererToggle()

	self.toggleRend.SetActivatable(true)
	//  0             1             2         3            4        5
	//is selected, bit position,field name, field value, field id,editable
	self._store = gtk.NewTreeStore(glib.G_TYPE_BOOL, glib.G_TYPE_INT, glib.G_TYPE_STRING, glib.G_TYPE_STRING, glib.G_TYPE_INT, glib.G_TYPE_BOOL)
	self._view = gtk.NewTreeView()

	self._view.AppendColumn(gtk.NewTreeViewColumnWithAttribute("Is Selected", self.toggleRend))
	self._view.AppendColumn(gtk.NewTreeViewColumnWithAttribute("Bit Position", fieldNameRenderer))
	self._view.AppendColumn(gtk.NewTreeViewColumnWithAttribute("Field Name", fieldNameRenderer))
	self._view.AppendColumn(gtk.NewTreeViewColumnWithAttribute("Data", self.valRend))

	self._view.GetColumn(0).AddAttribute(self.toggleRend, "active", 0)
	self._view.GetColumn(1).AddAttribute(fieldNameRenderer, "text", 1)
	self._view.GetColumn(2).AddAttribute(fieldNameRenderer, "text", 2)
	self._view.GetColumn(3).AddAttribute(self.valRend, "text", 3)
	self._view.GetColumn(3).AddAttribute(self.valRend, "editable", 5)

	for i := 0; i < 4; i++ {
		self._view.GetColumn(i).SetProperty("resizable", glib.ValueFromNative(true))
	}

	self._view.ModifyFontEasy("Dejavu Sans 8")

	self._view.SetModel(self._store)

	return self

}

func (pTree *PaysimSpecMsgTree) Store() *gtk.TreeStore {
	return pTree._store
}

func (pTree *PaysimSpecMsgTree) View() *gtk.TreeView {
	return pTree._view
}

func (pTree *PaysimSpecMsgTree) FieldValueRenderer() *gtk.CellRendererText {
	return pTree.valRend
}

func (pTree *PaysimSpecMsgTree) FieldToggleRenderer() *gtk.CellRendererToggle {
	return pTree.toggleRend
}
