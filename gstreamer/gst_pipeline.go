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
	"github.com/diamondburned/gotk4/pkg/glib/v2"
)

func (p *GstPipeline) Object() *glib.Object {
	return coreglib.AssumeOwnership(unsafe.Pointer(p.native))
}

func (p *GstPipeline) Link(elem GstElementer) error {
	return elementLink(p, elem)
}

func (p *GstPipeline) SetProperty(name string, value any) {
	elementSetProperty(p, name, value)
}

func (p *GstPipeline) Property(name string) any {
	return elementProperty(p, name)
}

func (p *GstPipeline) SetState(state GstState) (GstStateChangeReturn, error) {
	return elementSetState(p, state)
}

func (p *GstPipeline) LinkMany(elems ...GstElementer) error {
	return elementLinkMany(append([]GstElementer{p}, elems...)...)
}

func (p *GstPipeline) AddMany(elems ...GstElementer) {
	for _, elem := range elems {
		p.Add(elem)
	}
}
