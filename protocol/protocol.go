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

// Package protocol
package protocol

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/OmegaRogue/weylus-desktop/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
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
	PointerEvent | WheelEvent | KeyboardEvent | Config | ~string
}
type MessageOutbound interface {
	map[WeylusCommand]PointerEvent | map[WeylusCommand]WheelEvent | map[WeylusCommand]KeyboardEvent | map[WeylusCommand]Config | ~string
}

func WrapMessage[T MessageOutboundContent](a T) any {
	wrapper := make(map[WeylusCommand]T)
	switch any(a).(type) {
	case PointerEvent:
		wrapper[WeylusCommandPointerEvent] = a
	case WheelEvent:
		wrapper[WeylusCommandWheelEvent] = a
	case KeyboardEvent:
		wrapper[WeylusCommandKeyboardEvent] = a
	case Config:
		wrapper[WeylusCommandConfig] = a
	case string, WeylusCommand:
		return a
	default:
		str, err := utils.GetUnderlyingString(a)
		if err != nil {
			log.Panic().Err(err).Msg("what the fuck did you do? (I'm genuinely curious) This should never happen, immediately report this as an issue")
		}
		return str
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
	case strings.Contains(dataString, string(WeylusResponseCapturableList)):
		var l CapturableList
		err := json.Unmarshal(data, &l)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal CapturableList")
		}
		return l, nil
	case strings.Contains(dataString, string(WeylusResponseConfigError)):
		var err2 WeylusConfigError
		err := json.Unmarshal(data, &err2)
		if err != nil {
			return nil, errors.Wrapf(err, "failed unmarshaling error: %s", dataString)
		}
		return &err2, nil
	case strings.Contains(dataString, string(WeylusResponseError)):
		var err2 WeylusError
		err := json.Unmarshal(data, &err2)
		if err != nil {
			return nil, errors.Wrapf(err, "failed unmarshaling error: %s", dataString)
		}
		return &err2, nil
	}
	var foo any
	err := json.Unmarshal(data, &foo)
	if err != nil {
		return nil, errors.Wrapf(err, "failed unmarshaling data: %s", dataString)
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
