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
	"io"
	"os"
	"unsafe"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type GstElement struct {
	native *C.GstElement
}

func (e *GstElement) Native() *C.GstElement {
	return e.native
}
func NewGstElement(factoryName, name string) *GstElement {
	elem := new(GstElement)
	var _factoryName *C.char
	var _name *C.char

	_factoryName = (*C.char)(unsafe.Pointer(C.CString(factoryName)))
	_name = (*C.char)(unsafe.Pointer(C.CString(name)))
	elem.native = C.gst_element_factory_make(_factoryName, _name)
	return elem
}

type GstPipeline struct {
	native *C.GstElement
}

func (p *GstPipeline) Native() *C.GstElement {
	return p.native
}

func NewGstPipeline(name string) *GstPipeline {
	pipe := new(GstPipeline)

	_name := (*C.char)(unsafe.Pointer(C.CString(name)))
	ptr := C.gstreamer_pipeline_new(_name)
	pipe.native = ptr
	return pipe
}

func (p *GstPipeline) Add(elem GstElementer) {
	C.gstreamer_bin_add(p.native, elem.Native())
}

type GstBuffer struct {
	native *C.GstBuffer
}

func NewGstBuffer(size int) *GstBuffer {
	buffer := new(GstBuffer)
	_size := C.size_t(size)
	buffer.native = C.gstreamer_new_buffer(_size)

	return buffer
}

func (b *GstBuffer) Fill(data []byte, offset int) int {
	_offset := C.size_t(offset)
	_size := C.size_t(len(data))
	return int(C.gstreamer_buffer_fill(b.native, _offset, C.CBytes(data), _size))
}

func (e *GstElement) AppSrcPushBuffer(b *GstBuffer) int {
	return int(C.gst_app_src_push_buffer(C.gstreamer_app_src_cast(e.native), b.native))
}

type AppSrcWriter struct {
	elem *GstElement
}

func NewAppSrcWriter(appsrc *GstElement) *AppSrcWriter {
	w := new(AppSrcWriter)
	w.elem = appsrc

	return w
}

func (a *AppSrcWriter) Close() error {
	C.gst_app_src_end_of_stream(C.gstreamer_app_src_cast(a.elem.native))
	return nil
}

func (a *AppSrcWriter) Write(p []byte) (n int, err error) {
	buf := NewGstBuffer(len(p))
	n = buf.Fill(p, 0)

	a.elem.AppSrcPushBuffer(buf)
	C.free(unsafe.Pointer(buf.native))
	return
}

var _ io.WriteCloser = &AppSrcWriter{}

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

	source := NewGstElement("v4l2src", "source")
	capsfilter := NewGstElement("capsfilter", "filter")

	//source.SetProperty("device", "/dev/video0")

	C.gstreamer_set_caps_example(capsfilter.native)

	convert := NewVideoConvert("convert")
	sink := NewGTK4PaintableSink("sink")

	picture := gtk.NewPicture()
	picture.SetPaintable(sink.PropertyPaintable())

	pipeline := NewGstPipeline("test-pipeline")

	pipeline.AddMany(source, capsfilter, convert, sink)

	if err := source.LinkMany(capsfilter, convert, sink); err != nil {
		C.gst_object_unref(C.gpointer(unsafe.Pointer(pipeline.native)))
		C.free(unsafe.Pointer(pipeline.native))
		log.Fatal().Err(err).Msg("Elements could not be linked.")
	}

	//source.SetProperty("pattern", 0)

	if _, err := pipeline.SetState(GstStatePlaying); err != nil {
		C.gst_object_unref(C.gpointer(unsafe.Pointer(pipeline.native)))
		C.free(unsafe.Pointer(pipeline.native))
		log.Fatal().Err(err).Msg("Unable to set the pipeline to the playing state.")
	}

	window.SetChild(picture)

	window.Show()
}
