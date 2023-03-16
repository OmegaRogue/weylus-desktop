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

// Package protocol
package protocol

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	CommandTryGetFrame       = "TryGetFrame"
	CommandGetCapturableList = "GetCapturableList"
)

type Config struct {
	UInputSupport bool   `json:"uinput_support"`
	CaptureCursor bool   `json:"capture_cursor"`
	CapturableID  uint   `json:"capturable_id"`
	MaxWidth      uint   `json:"max_width"`
	MaxHeight     uint   `json:"max_height"`
	ClientName    string `json:"client_name,omitempty"`
}

type MessageOutboundContent interface {
	PointerEvent | WheelEvent | KeyboardEvent | Config
}
type MessageOutbound interface {
	map[string]PointerEvent | map[string]WheelEvent | map[string]KeyboardEvent | map[string]Config | ~string
}

func WrapMessage[T MessageOutboundContent | ~string](a T) any {
	wrapper := make(map[string]T)
	switch any(a).(type) {
	case PointerEvent:
		wrapper["PointerEvent"] = a
	case WheelEvent:
		wrapper["WheelEvent"] = a
	case KeyboardEvent:
		wrapper["KeyboardEvent"] = a
	case Config:
		wrapper["Config"] = a
	case string:
		return a
	default:
		log.Fatal().Msg("invalid type")
	}
	return wrapper
}

type CapturableList struct {
	CapturableList []string `json:"CapturableList"`
}

type MessageInbound interface {
	CapturableList | ~string | WeylusError | WeylusConfigError
}

func ParseMessage(data []byte) (any, error) {
	dataString := string(data)
	switch {
	case strings.Contains(dataString, "CapturableList"):
		var l CapturableList
		err := json.Unmarshal(data, &l)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal CapturableList")
		}
		return l, nil
	case strings.Contains(dataString, "ConfigError"):
		var err2 WeylusConfigError
		err := json.Unmarshal(data, &err2)
		if err != nil {
			return "", errors.Errorf("failed unmarshaling error: %s", dataString)
		}
		return &err2, nil
	case strings.Contains(dataString, "Error"):
		var err2 WeylusError
		err := json.Unmarshal(data, &err2)
		if err != nil {
			return "", errors.Errorf("failed unmarshaling error: %s", dataString)
		}
		return &err2, nil
	}
	var foo any
	err := json.Unmarshal(data, &foo)
	if err != nil {
		return nil, errors.Errorf("")
	}
	return foo, nil
}

type WeylusError struct {
	ErrorMessage string `json:"Error"`
}
type WeylusConfigError struct {
	ErrorMessage string `json:"ConfigError"`
}

func (e *WeylusError) Error() string {
	return fmt.Sprintf("WeylusError: %s", e.ErrorMessage)
}
func (e *WeylusConfigError) Error() string {
	return fmt.Sprintf("WeylusConfigError: %s", e.ErrorMessage)
}

type PointerEvent struct {
	EventType   PointerEventType `json:"event_type"`
	PointerType PointerType      `json:"pointer_type"`
	X           float64          `json:"x"`
	Y           float64          `json:"y"`
	Pressure    float64          `json:"pressure"`
	Width       float64          `json:"width"`
	Height      float64          `json:"height"`
	PointerID   int              `json:"pointer_id"`
	Timestamp   uint64           `json:"timestamp"`
	MovementX   int64            `json:"movement_x"`
	MovementY   int64            `json:"movement_y"`
	TiltX       int32            `json:"tilt_x"`
	TiltY       int32            `json:"tilt_y"`
	Twist       int32            `json:"twist"`
	Button      ButtonFlags      `json:"button"`
	Buttons     ButtonFlags      `json:"buttons"`
	IsPrimary   bool             `json:"is_primary"`
}

type WheelEvent struct {
	Dx        int32  `json:"dx"`
	Dy        int32  `json:"dy"`
	Timestamp uint64 `json:"timestamp"`
}

type KeyboardEvent struct {
	EventType KeyboardEventType `json:"event_type"`
	Code      string            `json:"code"`
	Key       string            `json:"key"`
	Location  KeyboardLocation  `json:"location"`
	Alt       bool              `json:"alt"`
	Ctrl      bool              `json:"ctrl"`
	Shift     bool              `json:"shift"`
	Meta      bool              `json:"meta"`
}
