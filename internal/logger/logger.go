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

package logger

import (
	stdlog "log"
	"os"

	"github.com/OmegaRogue/weylus-desktop/logger/gliblogger"
	"github.com/OmegaRogue/weylus-desktop/logger/journald"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func SetupLogger() {
	var multi zerolog.LevelWriter
	journaldWriter := journald.NewBetterJournaldWriter()

	if os.Getenv("WEYLUS_LOG_JSON") == "true" {
		multi = zerolog.MultiLevelWriter(os.Stderr, journaldWriter)
	} else {
		consoleWriter := zerolog.ConsoleWriter{
			Out:           os.Stderr,
			FieldsExclude: []string{journald.ThreadFieldName, gliblogger.GlibLevelFieldName},
		}
		multi = zerolog.MultiLevelWriter(consoleWriter, journaldWriter)
	}
	log.Logger = log.Output(multi).With().Caller().Logger().Hook(journald.ThreadHook{})
	stdLogger := log.With().Str("component", "stdlog").Logger()
	stdlog.SetOutput(stdLogger)
	if os.Getenv("G_DEBUG") != "" {
		glibLog := log.With().Str("component", "glib").Logger()
		glib.LogSetWriter(gliblogger.LoggerHandler(&glibLog))
	}
	switch os.Getenv("WEYLUS_LOG_LEVEL") {
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "WARN":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "TRACE":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}
