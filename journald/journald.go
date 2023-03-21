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

package journald

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"

	"github.com/coreos/go-systemd/v22/journal"
	"github.com/coreos/go-systemd/v22/sdjournal"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"weylus-surface/utils/pthread"
)

const (
	ThreadFieldName = "thread"
)

type betterJournaldWriter struct {
}

func (b betterJournaldWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	var output map[string]interface{}
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

func (b betterJournaldWriter) WriteJSONLevel(level zerolog.Level, o map[string]interface{}, p []byte) (int, error) {
	var message string
	prio := zerologLeveltoJournaldPriority(level)
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

	err := journal.Send(message, prio, nil)
	if err != nil {
		return 0, errors.Wrap(err, "send journal message")
	}
	args["JSON"] = string(p)
	err = journal.Send(message, prio, args)

	return len(p), nil
}

func zerologLeveltoJournaldPriority(level zerolog.Level) journal.Priority {
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

func (b betterJournaldWriter) Write(p []byte) (int, error) {
	var output map[string]interface{}
	if err := json.Unmarshal(p, &output); err != nil {
		return 0, errors.Wrap(err, "unmarshal intermediate log message")
	}
	lvl, err := zerolog.ParseLevel(output["level"].(string))
	if err != nil {
		return 0, errors.Wrap(err, "parse level")
	}
	return b.WriteJSONLevel(lvl, output, p)
}

func NewBetterJournaldWriter() io.Writer {
	return betterJournaldWriter{}
}
