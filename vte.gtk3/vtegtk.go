// Package vte is a cgo binding for Vte. Supports version 2.91 (0.40) and later.
//
// This package provides the Vte terminal wrapped as a gotk3 widget, and using
// this library ressources.
//
// https://developer.gnome.org/vte/0.40/VteTerminal.html
// https://developer.gnome.org/vte/unstable/VteTerminal.html
// https://golang.hotexamples.com/examples/c/-/toVTerminal/golang-tovterminal-function-examples.html
// https://github.com/golang/go/issues/58625
//

package vte

// #include <vte/vte.h>
// #cgo pkg-config: vte-2.91
import "C"

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"github.com/mrhdias/vte"

	"errors"
	"runtime"
	"unsafe"
)

// Terminal is a representation of Vte's VteTerminal.
type Terminal struct {
	gtk.Widget
	vte.Terminal
}

// NewTerminal creates a new terminal widget.
func NewTerminal() *Terminal {
	c := vte.NewTerminal()
	if c == nil {
		return nil
	}

	// obj := &glib.Object{glib.ToGObject(unsafe.Pointer(c.Native()))}
	obj := &glib.Object{
		GObject: glib.ToGObject(unsafe.Pointer(c.Native())),
	}

	obj.RefSink()
	runtime.SetFinalizer(obj, (*glib.Object).Unref)

	return wrapTerminal(obj, c)
}

// NewTerminalWindow creates a new terminal widget packed in a dedicated window.
func NewTerminalWindow() (*Terminal, *gtk.Window, error) {
	window, e := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if e != nil {
		return nil, nil, e
	}

	terminal := NewTerminal()
	if terminal == nil {
		return nil, nil, errors.New("create terminal failed")
	}

	// Packing.
	swin, _ := gtk.ScrolledWindowNew(nil, nil)
	swin.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
	swin.Add(terminal)
	window.Add(swin)
	window.ShowAll()

	return terminal, window, nil
}

func (v *Terminal) termNative() *C.VteTerminal {
	return (*C.VteTerminal)(unsafe.Pointer(v.Terminal.Native()))
}

func wrapTerminal(obj *glib.Object, term *vte.Terminal) *Terminal {
	// return &Terminal{gtk.Widget{glib.InitiallyUnowned{obj}}, *term}
	return &Terminal{gtk.Widget{
		InitiallyUnowned: glib.InitiallyUnowned{
			Object: obj,
		},
	}, *term}
}

// SetBgColor sets the background color for text which does not have a specific
// background color assigned. Only has effect when no background image is set
// and when the terminal is not transparent.
// The gdk RGBA is used as input.
func (v *Terminal) SetBgColor(color *gdk.RGBA) {
	// C.vte_terminal_set_color_background(v.termNative(), (*C.GdkRGBA)(unsafe.Pointer(color.Native())))

	// var addr uintptr = color.Native()
	// p := *(*unsafe.Pointer)(unsafe.Pointer(&addr))
	// C.vte_terminal_set_color_background(v.termNative(), (*C.GdkRGBA)(p))
	C.vte_terminal_set_color_background(v.termNative(), (*C.GdkRGBA)(uintptrToUnsafePointer(color.Native())))
}

// SetFgColor sets the foreground color used to draw normal text.
// The gdk RGBA is used as input.
func (v *Terminal) SetFgColor(color *gdk.RGBA) {
	// C.vte_terminal_set_color_foreground(v.termNative(), (*C.GdkRGBA)(unsafe.Pointer(color.Native())))

	// var addr uintptr = color.Native()
	// p := *(*unsafe.Pointer)(unsafe.Pointer(&addr))
	// C.vte_terminal_set_color_foreground(v.termNative(), (*C.GdkRGBA)(p))
	C.vte_terminal_set_color_foreground(v.termNative(), (*C.GdkRGBA)(uintptrToUnsafePointer(color.Native())))
}

// SetFont sets the font used for rendering all text displayed by the terminal,
// overriding any fonts set using widget.ModifyFont().
// The terminal will immediately attempt to load the desired font, retrieve its
// metrics, and attempt to resize itself to keep the same number of rows and
// columns. The font scale is applied to the specified font.
// The pango FontDescription is used as input.

func (v *Terminal) SetFont(fontDesc *pango.FontDescription) {
	// C.vte_terminal_set_font(v.termNative(), (*C.PangoFontDescription)(unsafe.Pointer(font.Native())))

	// var addr uintptr = fontDesc.Native()
	// p := *(*unsafe.Pointer)(unsafe.Pointer(&addr))
	// C.vte_terminal_set_font(v.termNative(), (*C.PangoFontDescription)(p))
	C.vte_terminal_set_font(v.termNative(), (*C.PangoFontDescription)(uintptrToUnsafePointer(fontDesc.Native())))
}

func uintptrToUnsafePointer(addr uintptr) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Pointer(&addr))
}
