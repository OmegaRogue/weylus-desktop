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

package gliblogger

import (
	stdlog "log"
	"os"
	"strings"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/rs/zerolog"
)

const (
	CodeFileFieldName  = "file"
	CodeFuncFieldName  = "func"
	CodeLineFieldName  = "line"
	DomainFieldName    = "domain"
	GlibLevelFieldName = "glib_level"
)

func logLevelGlibZerolog(lvl glib.LogLevelFlags) zerolog.Level {
	switch lvl & 0b11111100 {
	case glib.LogLevelError:
		return zerolog.ErrorLevel
	case glib.LogLevelCritical:
		return zerolog.FatalLevel
	case glib.LogLevelWarning:
		return zerolog.WarnLevel
	case glib.LogLevelMessage:
		return zerolog.InfoLevel
	case glib.LogLevelDebug:
		return zerolog.DebugLevel
	}
	return zerolog.NoLevel
}

// LoggerHandler creates a new LogWriterFunc that LogUseLogger uses. For
// more information, see LogUseLogger's documentation.
//
//goland:noinspection SpellCheckingInspection
func LoggerHandler(l zerolog.Logger) glib.LogWriterFunc {

	// Treat Lshortfile and Llongfile the same, because we don't have
	// the full path in codeFile anyway.
	Lfile := stdlog.Flags()&(stdlog.Lshortfile|stdlog.Llongfile) != 0

	// Support $G_MESSAGES_DEBUG.
	debugDomains := make(map[string]struct{})
	for _, debugDomain := range strings.Fields(os.Getenv("G_MESSAGES_DEBUG")) {
		debugDomains[debugDomain] = struct{}{}
	}
	// Special case: G_MESSAGES_DEBUG=all.
	_, debugAll := debugDomains["all"]

	return func(lvl glib.LogLevelFlags, fields []glib.LogField) glib.LogWriterOutput {
		var codeFile, codeLine, codeFunc string
		domain := "GLib (no domain)"

		event := l.WithLevel(logLevelGlibZerolog(lvl))

		for _, field := range fields {
			if !Lfile {
				switch field.Key() {
				case "MESSAGE":
					event.Str(zerolog.MessageFieldName, field.Value())
				case "GLIB_DOMAIN":
					event.Str(DomainFieldName, field.Value())
					domain = field.Value()
				}
				// Skip setting code* if we don't have to.
				continue
			}

			switch field.Key() {
			case "MESSAGE":
				event.Str(zerolog.MessageFieldName, field.Value())
			case "CODE_FILE":
				codeFile = field.Value()
			case "CODE_LINE":

				codeLine = field.Value()
			case "CODE_FUNC":
				codeFunc = field.Value()
			case "GLIB_DOMAIN":
				event.Str(DomainFieldName, field.Value())
				domain = field.Value()
			}
		}

		if !debugAll && (lvl&glib.LogLevelDebug != 0) && domain != "" {
			if _, ok := debugDomains[domain]; !ok {
				return glib.LogWriterHandled
			}
		}

		// Minor issue: this works badly if consts are OR'd together.
		// Probably never.
		level := strings.TrimPrefix(lvl.String(), "Level")
		event.Str(GlibLevelFieldName, level)

		if !Lfile || (codeFile == "" && codeLine == "") {
			event.Send()
			return glib.LogWriterHandled
		}

		if codeFunc == "" {
			event.Str(CodeLineFieldName, codeLine).Str(CodeFileFieldName, codeFile).Send()
			return glib.LogWriterHandled
		}
		event.Str(CodeFuncFieldName, codeFunc).Send()
		return glib.LogWriterHandled
	}
}
