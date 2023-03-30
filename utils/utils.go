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

package utils

import (
	"bufio"
	"io"

	"github.com/samber/lo"
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
