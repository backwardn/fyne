package glfw

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
	publicWidget "fyne.io/fyne/widget"
)

var _ fyne.Widget = (*MenuBarItem)(nil)
var _ desktop.Hoverable = (*MenuBarItem)(nil)

// MenuBarItem is a widget for displaying an item for a fyne.Menu in a MenuBar.
type MenuBarItem struct {
	widget.Base
	Child  *publicWidget.Menu
	Menu   *fyne.Menu
	Parent *MenuBar

	hovered bool
}

// CreateRenderer returns a new renderer for the menu bar item.
// Implements: fyne.Widget
func (i *MenuBarItem) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(i.Menu.Label, theme.TextColor())
	objects := []fyne.CanvasObject{text}
	if i.Child == nil {
		child := publicWidget.NewMenu(i.Menu)
		child.Hide()
		child.DismissAction = i.Parent.deactivate
		i.Child = child
	}
	objects = append(objects, i.Child)

	return &menuBarItemRenderer{
		widget.NewBaseRenderer(objects),
		i,
		text,
	}
}

// Hide hides the menu bar item.
// Implements: fyne.Widget
func (i *MenuBarItem) Hide() {
	widget.HideWidget(&i.Base, i)
}

// MinSize returns the minimal size of the menu bar item.
// Implements: fyne.Widget
func (i *MenuBarItem) MinSize() fyne.Size {
	return widget.MinSizeOf(i)
}

// MouseIn changes the item to be hovered and shows the menu if the bar is active.
// The menu that was displayed before will be hidden.
// Implements: desktop.Hoverable
func (i *MenuBarItem) MouseIn(_ *desktop.MouseEvent) {
	if i.Parent.active {
		i.hovered = true
		i.activateChild()
		i.Refresh()
	} else {
		i.hovered = true
		i.Refresh()
	}
}

// MouseMoved does nothing.
// Implements: desktop.Hoverable
func (i *MenuBarItem) MouseMoved(_ *desktop.MouseEvent) {
}

// MouseOut changes the item to not be hovered but has no effect on the visibility of the menu.
// Implements: desktop.Hoverable
func (i *MenuBarItem) MouseOut() {
	i.hovered = false
	i.Refresh()
}

// Refresh triggers a redraw of the menu bar item.
// Implements: fyne.Widget
func (i *MenuBarItem) Refresh() {
	widget.RefreshWidget(i)
}

// Resize changes the size of the menu bar item.
// Implements: fyne.Widget
func (i *MenuBarItem) Resize(size fyne.Size) {
	widget.ResizeWidget(&i.Base, i, size)
	i.updateChildPosition()
}

// Show makes the menu bar item visible.
// Implements: fyne.Widget
func (i *MenuBarItem) Show() {
	widget.ShowWidget(&i.Base, i)
}

// Tapped toggles the activation state of the menu bar.
// It shows the item’s menu if the bar is activated and hides it if the bar is deactivated.
// Implements: fyne.Tappable
func (i *MenuBarItem) Tapped(*fyne.PointEvent) {
	if i.Parent.active {
		i.Parent.deactivate()
	} else {
		i.Parent.activate()
		i.activateChild()
	}
	i.Refresh()
}

func (i *MenuBarItem) activateChild() {
	if i.Child != nil {
		i.Child.DeactivateChild()
	}
	if i.Parent.activeChild == i.Child {
		return
	}

	if i.Parent.activeChild != nil {
		i.Parent.activeChild.Hide()
	}
	i.Parent.activeChild = i.Child
	if i.Child != nil {
		if i.Child.Size().IsZero() {
			i.Child.Resize(i.Child.MinSize())
			i.updateChildPosition()
		}
		i.Child.Show()
	}
}

func (i *MenuBarItem) updateChildPosition() {
	i.Child.Move(fyne.NewPos(0, i.Size().Height))
}

type menuBarItemRenderer struct {
	widget.BaseRenderer
	i    *MenuBarItem
	text *canvas.Text
}

func (r *menuBarItemRenderer) BackgroundColor() color.Color {
	if r.i.hovered || (r.i.Child != nil && r.i.Child.Visible()) {
		return theme.HoverColor()
	}

	return color.Transparent
}

func (r *menuBarItemRenderer) Layout(_ fyne.Size) {
	padding := r.padding()

	r.text.TextSize = theme.TextSize()
	r.text.Color = theme.TextColor()
	r.text.Resize(r.text.MinSize())
	r.text.Move(fyne.NewPos(padding.Width/2, padding.Height/2))
}

func (r *menuBarItemRenderer) MinSize() fyne.Size {
	return r.text.MinSize().Add(r.padding())
}

func (r *menuBarItemRenderer) Refresh() {
	canvas.Refresh(r.i)
}

func (r *menuBarItemRenderer) padding() fyne.Size {
	return fyne.NewSize(theme.Padding()*4, theme.Padding()*2)
}
