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
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <https://www.gnu.org/licenses/>.
 */

package main

import "C"
import (
	"unsafe"

	coreglib "github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
)

func (e *GstElement) Object() *glib.Object {
	return coreglib.AssumeOwnership(unsafe.Pointer(e.native))
}

func (e *GstElement) Link(elem GstElementer) error {
	return elementLink(e, elem)
}

func (e *GstElement) SetProperty(name string, value any) {
	elementSetProperty(e, name, value)
}

func (e *GstElement) Property(name string) any {
	return elementProperty(e, name)
}

func (e *GstElement) SetState(state GstState) (GstStateChangeReturn, error) {
	return elementSetState(e, state)
}

func (e *GstElement) LinkMany(elems ...GstElementer) error {
	return elementLinkMany(append([]GstElementer{e}, elems...)...)
}

func (e *GstElement) PropertyPaintable() gdk.Paintabler {
	return coreglib.NewValue(e.Property("paintable")).Object().Cast().(gdk.Paintabler)
}

func NewAppSource(name string) *GstElement {
	return NewGstElement("appsrc", name)
}
func NewVideoConvert(name string) *GstElement {
	return NewGstElement("videoconvert", name)
}
func NewGTK4PaintableSink(name string) *GstElement {
	return NewGstElement("gtk4paintablesink", name)
}
func NewVideoTestSource(name string) *GstElement {
	return NewGstElement("videotestsrc", name)
}
