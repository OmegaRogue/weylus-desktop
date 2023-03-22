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

// Package main
package main

import (
	stdlog "log"
	"os"

	"github.com/OmegaRogue/weylus-desktop/cmd"
	"github.com/OmegaRogue/weylus-desktop/logger/gliblogger"
	"github.com/OmegaRogue/weylus-desktop/logger/journald"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func main() {
	consoleWriter := zerolog.ConsoleWriter{
		Out:           os.Stderr,
		FieldsExclude: []string{journald.ThreadFieldName, gliblogger.GlibLevelFieldName},
	}
	multi := zerolog.MultiLevelWriter(consoleWriter, journald.NewBetterJournaldWriter())
	log.Logger = log.Output(multi).With().Caller().Logger().Hook(journald.ThreadHook{})
	stdlog.SetFlags(0)
	stdLogger := log.With().Str("component", "stdlog").Logger()
	stdlog.SetOutput(stdLogger)
	//glibLog := log.With().Str("component", "glib").Logger()
	//glib.LogSetWriter(gliblogger.LoggerHandler(glibLog))
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	cmd.Execute()
}
