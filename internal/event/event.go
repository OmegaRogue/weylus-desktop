/*
 * Copyright © 2023 omegarogue
 * SPDX-License-Identifier: GPL-3.0-or-later
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

// Package event
package event

import (
	"math"
	"strings"
	"time"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"weylus-surface/protocol"
)

type ControllerManager struct {
	Stylus *gtk.GestureStylus
	Click  *gtk.GestureClick
	Drag   *gtk.GestureDrag
	Motion *gtk.EventControllerMotion
	Scroll *gtk.EventControllerScroll
	Key    *gtk.EventControllerKey

	StylusState protocol.PointerEvent
	TouchState  protocol.PointerEvent
	MouseState  protocol.PointerEvent
	KeyState    protocol.KeyboardEvent
	ScrollState protocol.WheelEvent

	callbacks []func(m *ControllerManager)
}

func (m *ControllerManager) AddCallback(f func(m *ControllerManager)) {
	m.callbacks = append(m.callbacks, f)
}

func (m *ControllerManager) runCallbacks() {
	for _, callback := range m.callbacks {
		callback(m)
	}
}

func NewControllerManager() *ControllerManager {
	m := new(ControllerManager)

	m.Stylus = gtk.NewGestureStylus()
	m.Stylus.SetButton(0)
	m.Stylus.SetExclusive(false)
	m.Click = gtk.NewGestureClick()
	m.Click.SetButton(0)
	m.Click.SetExclusive(false)
	m.Drag = gtk.NewGestureDrag()
	m.Drag.SetButton(0)
	m.Drag.SetExclusive(false)
	m.Drag.SetTouchOnly(true)
	m.Motion = gtk.NewEventControllerMotion()
	m.Key = gtk.NewEventControllerKey()
	m.Scroll = gtk.NewEventControllerScroll(gtk.EventControllerScrollBothAxes)

	m.Stylus.ConnectUp(m.StylusUpEventHandler)
	m.Stylus.ConnectDown(m.StylusDownEventHandler)
	m.Stylus.ConnectProximity(m.StylusProximityEventHandler)
	m.Stylus.ConnectMotion(m.StylusMotionEventHandler)

	m.Key.ConnectKeyPressed(m.KeyDownHandler)
	m.Key.ConnectKeyReleased(m.KeyReleasedHandler)
	m.Key.ConnectModifiers(m.KeyModHandler)

	m.Scroll.ConnectScroll(m.ScrollHandler)

	m.Motion.ConnectMotion(m.MotionHandler)

	m.Click.ConnectPressed(m.PressedHandler)
	m.Click.ConnectReleased(m.ReleasedHandler)
	m.Click.ConnectUnpairedRelease(m.UnpairedReleaseHandler)

	m.Drag.ConnectDragBegin(m.DragBeginHandler)
	m.Drag.ConnectDragEnd(m.DragEndHandler)
	m.Drag.ConnectDragUpdate(m.DragUpdateHandler)

	return m
}

func (m *ControllerManager) ScrollHandler(dx, dy float64) (ok bool) {
	ok = false
	m.ScrollState.Timestamp = uint64(time.Now().UnixMilli())
	defer m.runCallbacks()
	m.ScrollState.Dx = int32(math.Round(10 * dx))
	m.ScrollState.Dy = int32(math.Round(10 * dy))
	return
}

func (m *ControllerManager) KeyDownHandler(keyVal, keycode uint, state gdk.ModifierType) (ok bool) {
	ok = false
	defer m.runCallbacks()
	m.KeyState.Alt = state.Has(gdk.AltMask)
	m.KeyState.Shift = state.Has(gdk.ShiftMask)
	m.KeyState.Ctrl = state.Has(gdk.ControlMask)
	m.KeyState.Meta = state.Has(gdk.MetaMask)
	m.KeyState.EventType = protocol.KeyboardEventTypeDown
	m.KeyState.Code = protocol.CodeValue[keycode]
	m.KeyState.Location = protocol.KeyboardLocationStandard
	var ok2 bool
	m.KeyState.Key, ok2 = protocol.KeyValue[keyVal]
	if !ok2 {
		m.KeyState.Key = string(rune(gdk.KeyvalToUnicode(keyVal)))
	}
	m.KeyState.EventType = protocol.KeyboardEventTypeDown

	return
}
func (m *ControllerManager) KeyReleasedHandler(keyVal, keycode uint, state gdk.ModifierType) {
	defer m.runCallbacks()
	m.KeyState.Alt = state.Has(gdk.AltMask)
	m.KeyState.Shift = state.Has(gdk.ShiftMask)
	m.KeyState.Ctrl = state.Has(gdk.ControlMask)
	m.KeyState.Meta = state.Has(gdk.MetaMask)
	m.KeyState.EventType = protocol.KeyboardEventTypeDown
	m.KeyState.Code = protocol.CodeValue[keycode]
	m.KeyState.Location = protocol.KeyboardLocationStandard
	var ok bool
	m.KeyState.Key, ok = protocol.KeyValue[keyVal]
	if !ok {
		m.KeyState.Key = string(rune(gdk.KeyvalToUnicode(keyVal)))
	}
	m.KeyState.EventType = protocol.KeyboardEventTypeUp
}
func (m *ControllerManager) KeyModHandler(keyVal gdk.ModifierType) (ok bool) {
	ok = false
	m.KeyState.Alt = keyVal.Has(gdk.AltMask)
	m.KeyState.Shift = keyVal.Has(gdk.ShiftMask)
	m.KeyState.Ctrl = keyVal.Has(gdk.ControlMask)
	m.KeyState.Meta = keyVal.Has(gdk.MetaMask)
	return
}

func (m *ControllerManager) ConnectControllers(overlay *gtk.Overlay) {
	overlay.AddController(m.Drag)
	overlay.AddController(m.Click)
	overlay.AddController(m.Stylus)
	overlay.AddController(m.Motion)
}

func (m *ControllerManager) stylusEventHandler(x, y float64) {
	m.StylusState.X = x
	m.StylusState.Y = y
	xTilt, _ := m.Stylus.Axis(gdk.AxisXtilt)
	yTilt, _ := m.Stylus.Axis(gdk.AxisYtilt)
	m.StylusState.TiltX = int32(90 * xTilt)
	m.StylusState.TiltY = int32(90 * yTilt)
	m.StylusState.Pressure, _ = m.Stylus.Axis(gdk.AxisPressure)
	m.StylusState.Timestamp = uint64(time.Now().UnixMilli())
}

func (m *ControllerManager) StylusUpEventHandler(x, y float64) {
	m.stylusEventHandler(x, y)
	defer m.runCallbacks()
	m.StylusState.TiltX = 0
	m.StylusState.TiltY = 0
	m.StylusState.Pressure = 0
	m.StylusState.MovementX = 0
	m.StylusState.MovementY = 0
	m.StylusState.Buttons &= ^(protocol.ButtonPrimary | protocol.ButtonEraser)
	m.StylusState.Button = protocol.ButtonNone
}

func (m *ControllerManager) StylusDownEventHandler(x, y float64) {
	m.stylusEventHandler(x, y)
	defer m.runCallbacks()
	if tool := m.Stylus.DeviceTool(); tool != nil {
		switch tool.ToolType() {
		case gdk.DeviceToolTypePen:
			m.StylusState.Button = protocol.ButtonPrimary
			m.StylusState.Buttons |= protocol.ButtonPrimary
			m.StylusState.Buttons &= ^protocol.ButtonEraser
		case gdk.DeviceToolTypeEraser:
			m.StylusState.Button = protocol.ButtonEraser
			m.StylusState.Buttons &= ^protocol.ButtonPrimary
			m.StylusState.Buttons |= protocol.ButtonEraser
		}
	}
}

func (m *ControllerManager) StylusProximityEventHandler(x, y float64) {
	m.stylusEventHandler(x, y)
	defer m.runCallbacks()
	m.StylusState.Buttons &= ^(protocol.ButtonPrimary | protocol.ButtonEraser)
	m.StylusState.Button = protocol.ButtonNone
}

func (m *ControllerManager) StylusMotionEventHandler(x, y float64) {
	m.stylusEventHandler(x, y)
	defer m.runCallbacks()
}

func (m *ControllerManager) PressedHandler(_ int, x, y float64) {
	defer m.runCallbacks()
	name := m.Click.CurrentEventDevice().ObjectProperty("name").(string)
	source := m.Click.CurrentEventDevice().ObjectProperty("source").(gdk.InputSource)
	button := m.Click.CurrentButton()
	var btn protocol.ButtonFlags
	switch button {
	case gdk.BUTTON_PRIMARY:
		btn = protocol.ButtonPrimary
	case gdk.BUTTON_MIDDLE:
		btn = protocol.ButtonAuxiliary
	case gdk.BUTTON_SECONDARY:
		btn = protocol.ButtonSecondary
	case 8:
		btn = protocol.ButtonFourth
	case 9:
		btn = protocol.ButtonFifth
	}
	switch {
	case strings.Contains(name, "Stylus"):
		m.StylusState.Timestamp = uint64(time.Now().UnixMilli())
		m.StylusState.X = x
		m.StylusState.Y = y
		m.StylusState.Button = btn
		m.StylusState.Buttons |= btn
	case source == gdk.SourceTouchscreen:
		m.TouchState.Timestamp = uint64(time.Now().UnixMilli())
		m.TouchState.X = x
		m.TouchState.Y = y
		m.TouchState.Button = btn
		m.TouchState.Buttons |= btn
	default:
		m.MouseState.Timestamp = uint64(time.Now().UnixMilli())
		m.MouseState.X = x
		m.MouseState.Y = y
		m.MouseState.Button = btn
		m.MouseState.Buttons |= btn
	}
}

func (m *ControllerManager) ReleasedHandler(_ int, x, y float64) {
	defer m.runCallbacks()
	name := m.Click.CurrentEventDevice().ObjectProperty("name").(string)
	source := m.Click.CurrentEventDevice().ObjectProperty("source").(gdk.InputSource)
	button := m.Click.CurrentButton()
	var btn protocol.ButtonFlags
	switch button {
	case gdk.BUTTON_PRIMARY:
		btn = protocol.ButtonPrimary
	case gdk.BUTTON_MIDDLE:
		btn = protocol.ButtonAuxiliary
	case gdk.BUTTON_SECONDARY:
		btn = protocol.ButtonSecondary
	case 8:
		btn = protocol.ButtonFourth
	case 9:
		btn = protocol.ButtonFifth
	}
	switch {
	case strings.Contains(name, "Stylus"):
		m.StylusState.Timestamp = uint64(time.Now().UnixMilli())
		m.StylusState.X = x
		m.StylusState.Y = y
		m.StylusState.Button = protocol.ButtonNone
		m.StylusState.Buttons &= ^btn
	case source == gdk.SourceTouchscreen:
		m.TouchState.Timestamp = uint64(time.Now().UnixMilli())
		m.TouchState.X = x
		m.TouchState.Y = y
		m.TouchState.Button = protocol.ButtonNone
		m.TouchState.Buttons &= ^btn
	default:
		m.MouseState.Timestamp = uint64(time.Now().UnixMilli())
		m.MouseState.X = x
		m.MouseState.Y = y
		m.MouseState.Button = protocol.ButtonNone
		m.MouseState.Buttons &= ^btn
	}
}
func (m *ControllerManager) MotionHandler(x, y float64) {
	defer m.runCallbacks()
	var name string
	var source gdk.InputSource
	if dev := m.Motion.CurrentEventDevice(); dev != nil {
		if n := dev.ObjectProperty("name"); n != nil {
			name = n.(string)
		}
		if s := dev.ObjectProperty("source"); s != nil {
			source = s.(gdk.InputSource)
		}
	}

	switch {
	case strings.Contains(name, "Stylus"):
		m.StylusState.Timestamp = uint64(time.Now().UnixMilli())
		m.StylusState.X = x
		m.StylusState.Y = y
	case source == gdk.SourceTouchscreen:
		m.TouchState.Timestamp = uint64(time.Now().UnixMilli())
		m.TouchState.X = x
		m.TouchState.Y = y
	default:
		m.MouseState.Timestamp = uint64(time.Now().UnixMilli())
		m.MouseState.X = x
		m.MouseState.Y = y
	}
}

func (m *ControllerManager) UnpairedReleaseHandler(x, y float64, button uint, _ *gdk.EventSequence) {
	defer m.runCallbacks()
	var name string
	var source gdk.InputSource
	if dev := m.Click.CurrentEventDevice(); dev != nil {
		if n := dev.ObjectProperty("name"); n != nil {
			name = n.(string)
		}
		if s := dev.ObjectProperty("source"); s != nil {
			source = s.(gdk.InputSource)
		}
	}
	var btn protocol.ButtonFlags
	switch button {
	case gdk.BUTTON_PRIMARY:
		btn = protocol.ButtonPrimary
	case gdk.BUTTON_MIDDLE:
		btn = protocol.ButtonAuxiliary
	case gdk.BUTTON_SECONDARY:
		btn = protocol.ButtonSecondary
	case 8:
		btn = protocol.ButtonFourth
	case 9:
		btn = protocol.ButtonFifth
	}
	switch {
	case strings.Contains(name, "Stylus"):
		m.StylusState.Timestamp = uint64(time.Now().UnixMilli())
		m.StylusState.Button = protocol.ButtonNone
		m.StylusState.Buttons &= ^btn
		m.StylusState.X = x
		m.StylusState.Y = y
	case source == 3:
		m.TouchState.Timestamp = uint64(time.Now().UnixMilli())
		m.TouchState.Button = protocol.ButtonNone
		m.TouchState.Buttons &= ^btn
		m.TouchState.X = x
		m.TouchState.Y = y
	default:
		m.MouseState.Timestamp = uint64(time.Now().UnixMilli())
		m.MouseState.Button = protocol.ButtonNone
		m.MouseState.Buttons &= ^btn
		m.MouseState.X = x
		m.MouseState.Y = y
	}
}

func (m *ControllerManager) DragBeginHandler(x, y float64) {
	defer m.runCallbacks()
	m.TouchState.X = x
	m.TouchState.Y = y
	m.TouchState.Timestamp = uint64(time.Now().UnixMilli())
}

func (m *ControllerManager) DragEndHandler(x, y float64) {
	defer m.runCallbacks()
	m.TouchState.X = x
	m.TouchState.Y = y
	m.TouchState.Timestamp = uint64(time.Now().UnixMilli())
}

func (m *ControllerManager) DragUpdateHandler(x, y float64) {
	defer m.runCallbacks()
	m.TouchState.X = x
	m.TouchState.Y = y
	m.TouchState.Timestamp = uint64(time.Now().UnixMilli())
}
