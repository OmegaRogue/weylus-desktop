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

package input

import (
	"fmt"

	"github.com/holoplot/go-evdev"
)

// TODO extend and switch to https://github.com/holoplot/go-evdev, doesnt support setting properties yet, i like the api more than this

type UInputDevice struct {
}

func createStylus() {
	_, err := evdev.CreateDevice(
		"fake-device",
		evdev.InputID{
			BusType: 0x03,
			Vendor:  0x4711,
			Product: 0x0816,
			Version: 1,
		},
		map[evdev.EvType][]evdev.EvCode{
			evdev.EV_KEY: {
				evdev.BTN_LEFT,
				evdev.BTN_RIGHT,
				evdev.BTN_MIDDLE,
			},
			evdev.EV_REL: {
				evdev.REL_X,
				evdev.REL_Y,
				evdev.REL_WHEEL,
				evdev.REL_HWHEEL,
			},
		},
	)
	if err != nil {
		fmt.Printf("failed to create device: %s", err.Error())
		return
	}
}
