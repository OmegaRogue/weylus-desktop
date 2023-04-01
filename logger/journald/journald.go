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

// Package journald
package journald

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/OmegaRogue/weylus-desktop/utils/pthread"
	"github.com/coreos/go-systemd/v22/journal"
	"github.com/coreos/go-systemd/v22/sdjournal"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	ThreadFieldName = "thread"
)

type betterJournaldWriter struct {
}

func (b betterJournaldWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	var output map[string]any
	if err := json.Unmarshal(p, &output); err != nil {
		return 0, errors.Wrap(err, "unmarshal intermediate log message")
	}
	return b.WriteJSONLevel(level, output, p)
}

type ThreadHook struct {
}

func (t ThreadHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	e.Uint64("thread", pthread.Self().ID())
}

// GetGID gets the current goroutine ID (copied from https://blog.sgmansfield.com/2015/12/goroutine-ids/)
func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func (b betterJournaldWriter) WriteJSONLevel(level zerolog.Level, o map[string]any, p []byte) (int, error) {
	var message string
	prio := zerologLevelToJournaldPriority(level)
	args := make(map[string]string)

	for key, value := range o {
		jKey := strings.ToUpper(key)
		switch key {
		case zerolog.LevelFieldName, zerolog.TimestampFieldName:
			continue
		case zerolog.MessageFieldName:
			message, _ = value.(string)
			continue
		case zerolog.CallerFieldName:

			call := strings.Split(value.(string), ":")
			args[sdjournal.SD_JOURNAL_FIELD_CODE_FILE] = call[0]
			args[sdjournal.SD_JOURNAL_FIELD_CODE_LINE] = call[1]
			continue
		case ThreadFieldName:
			args["TID"] = fmt.Sprint(uint64(value.(float64)))
			continue
		case zerolog.ErrorStackFieldName:
			if stackTrace, ok := value.([]any); ok {
				if frame, ok := stackTrace[0].(map[string]any); ok {
					args[sdjournal.SD_JOURNAL_FIELD_CODE_FUNC] = frame["func"].(string)
					args[sdjournal.SD_JOURNAL_FIELD_CODE_LINE] = frame["line"].(string)
					args[sdjournal.SD_JOURNAL_FIELD_CODE_FILE] = frame["source"].(string)
				}
			}
		}
		switch v := value.(type) {
		case string:
			args[jKey] = v
		case json.Number:
			args[jKey] = fmt.Sprint(value)
		default:
			b, err := zerolog.InterfaceMarshalFunc(value)
			if err != nil {
				args[jKey] = fmt.Sprintf("[error: %v]", err)
			} else {
				args[jKey] = string(b)
			}
		}
	}

	args["JSON"] = string(p)
	if err := journal.Send(message, prio, args); err != nil {
		return 0, errors.Wrap(err, "send journal message")
	}
	return len(p), nil
}

func zerologLevelToJournaldPriority(level zerolog.Level) journal.Priority {
	switch level {
	case zerolog.TraceLevel:
		return journal.PriDebug
	case zerolog.DebugLevel:
		return journal.PriDebug
	case zerolog.InfoLevel:
		return journal.PriInfo
	case zerolog.WarnLevel:
		return journal.PriWarning
	case zerolog.ErrorLevel:
		return journal.PriErr
	case zerolog.FatalLevel:
		return journal.PriCrit
	case zerolog.PanicLevel:
		return journal.PriEmerg
	case zerolog.NoLevel:
		return journal.PriNotice
	}
	return journal.PriNotice
}

func (b betterJournaldWriter) Write(p []byte) (n int, err error) {
	var output map[string]any
	if err := json.Unmarshal(p, &output); err != nil {
		return 0, errors.Wrap(err, "unmarshal intermediate log message")
	}
	level, ok := output["level"].(string)
	var lvl zerolog.Level
	if ok {
		lvl, err = zerolog.ParseLevel(level)
		if err != nil {
			return 0, errors.Wrap(err, "parse level")
		}
	} else {
		lvl = zerolog.NoLevel
	}

	return b.WriteJSONLevel(lvl, output, p)
}

func NewBetterJournaldWriter() zerolog.LevelWriter {
	return betterJournaldWriter{}
}
