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

package utils

import (
	"bufio"
	"io"
	"reflect"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"nhooyr.io/websocket"
)

func Remove[V any](collection []V, index int) []V {
	return lo.Filter(collection, func(_ V, i int) bool {
		return i != index
	})
}

func NewBufPipe() *bufio.ReadWriter {
	pr, pw := io.Pipe()
	return bufio.NewReadWriter(bufio.NewReader(pr), bufio.NewWriter(pw))
}

type Msg struct {
	Type websocket.MessageType
	Data []byte
}

func IsStringUnderlying(v any) bool {
	b := reflect.ValueOf(v)
	return b.Kind() == reflect.String
}
func GetUnderlyingString(v any) (string, error) {
	if b := reflect.ValueOf(v); b.Kind() == reflect.String {
		return b.String(), nil
	} else {
		return "", errors.Errorf("%v of invalid type %s", v, b.Type().String())
	}
}
