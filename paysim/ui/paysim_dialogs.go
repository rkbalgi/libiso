package ui
import(
		"github.com/mattn/go-gtk/gtk"
)
func ShowInfoDialog(widget gtk.IWidget, msg string) {

	dialog := gtk.NewMessageDialog(
		widget.GetTopLevelAsWindow(),
		gtk.DIALOG_MODAL,
		gtk.MESSAGE_INFO,
		gtk.BUTTONS_OK,
		msg)
	dialog.SetTitle("Paysim Info Message")
	dialog.Response(func() {
		dialog.Destroy()
	})
	dialog.Run()
}


func ShowErrorDialog(widget gtk.IWidget, msg string) {

	dialog := gtk.NewMessageDialog(
		widget.GetTopLevelAsWindow(),
		gtk.DIALOG_MODAL,
		gtk.MESSAGE_ERROR,
		gtk.BUTTONS_OK,
		msg)
	dialog.SetTitle("Paysim Error Message")
	dialog.Response(func() {
		dialog.Destroy()
	})
	dialog.Run()
}