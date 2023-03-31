/*
 * Copyright © 2023 omegarogue
 * SPDX-License-Identifier: AGPL-3.0-or-later
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

// #cgo pkg-config: gstreamer-1.0 glib-2.0 gstreamer-app-1.0 gstreamer-video-1.0 gtk4
// #include "go_gstreamer.h"
import "C"

import (
	"os"
	"unsafe"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func main() {
	app := gtk.NewApplication("com.github.diamondburned.gotk4-examples.gtk4.simple", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() { activate(app) })

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}

}

//func GetGtk4Gstreamer() gdk.Paintabler {
//	ptr := C.initialize()
//	obj := coreglib.AssumeOwnership(unsafe.Pointer(ptr))
//	return obj.Cast().(gdk.Paintabler)
//}

type GstElementer interface {
	Object() *glib.Object
	Native() *C.GstElement
	Link(elem GstElementer) error
	LinkMany(elems ...GstElementer) error
	SetProperty(name string, value any)
	Property(name string) any
	SetState(state GstState) (GstStateChangeReturn, error)
}

func elementLinkMany(elems ...GstElementer) error {
	var lastElem GstElementer
	for i, elem := range elems {
		if i == 0 {
			lastElem = elem
			continue
		}
		if err := lastElem.Link(elem); err != nil {
			return err
		}
		lastElem = elem
	}
	return nil
}

func elementSetProperty(e GstElementer, name string, value any) {
	e.Object().SetObjectProperty(name, value)
}
func elementProperty(e GstElementer, name string) any {
	return e.Object().ObjectProperty(name)
}

func elementSetState(e GstElementer, state GstState) (GstStateChangeReturn, error) {
	ret := GstStateChangeReturn(int(C.gstreamer_element_set_state(e.Native(), C.int(state))))
	if ret == GstStateChangeReturnFailure {
		return GstStateChangeReturnFailure, errors.New("State change resulted in failure")
	}
	return ret, nil
}

func elementLink(src, dest GstElementer) error {
	r := C.gst_element_link(src.Native(), dest.Native())
	if int(r) != 1 {
		return errors.New("Elements could not be linked.")
	}
	return nil
}

func activate(app *gtk.Application) {
	C.gstreamer_init()
	window := gtk.NewApplicationWindow(app)
	window.SetTitle("gotk4 Example")
	window.SetDefaultSize(400, 300)

	source := NewGstElement("videotestsrc", "source")
	convert := NewGstElement("videoconvert", "convert")
	sink := NewGstElement("gtk4paintablesink", "sink")

	picture := gtk.NewPicture()
	picture.SetPaintable(sink.PropertyPaintable())

	pipeline := NewGstPipeline("test-pipeline")

	pipeline.AddMany(source, convert, sink)

	if err := source.LinkMany(convert, sink); err != nil {
		C.gst_object_unref(C.gpointer(unsafe.Pointer(pipeline.native)))
		C.free(unsafe.Pointer(pipeline.native))
		log.Fatal().Err(err).Msg("Elements could not be linked.")
	}

	source.SetProperty("pattern", 0)

	if _, err := pipeline.SetState(GstStatePlaying); err != nil {
		C.gst_object_unref(C.gpointer(unsafe.Pointer(pipeline.native)))
		C.free(unsafe.Pointer(pipeline.native))
		log.Fatal().Err(err).Msg("Unable to set the pipeline to the playing state.")
	}

	window.SetChild(picture)

	window.Show()
}
