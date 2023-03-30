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

package server

import (
	"context"
	_ "embed"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/OmegaRogue/weylus-desktop/web"
	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"nhooyr.io/websocket"
)

type data struct {
	AccessCode           string
	WebsocketPort        uint16
	LogLevel             int
	UInputEnabled        bool
	CaptureCursorEnabled bool
}

func middleware(logger zerolog.Logger) alice.Chain {

	c := alice.New()
	c = c.Append(hlog.NewHandler(logger))
	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))

	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.UserAgentHandler("user_agent"))
	c = c.Append(hlog.RefererHandler("referer"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))
	return c
}

func WeylusWeb(logger zerolog.Logger) {
	c := middleware(logger)
	h := c.Then(http.HandlerFunc(HandleWebsite))
	s := http.NewServeMux()

	s.Handle("/", h)
	s.Handle("/style.css", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		if _, err := w.Write([]byte(web.StyleCSS)); err != nil {
			hlog.FromRequest(r).Err(err).Msg("error on write style.css")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	s.Handle("/access_code.html", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		if _, err := w.Write([]byte(web.AccessHTML)); err != nil {
			hlog.FromRequest(r).Err(err).Msg("error on write access_code.html")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	s.Handle("/lib.js", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/javascript")
		if _, err := w.Write([]byte(web.LibJS)); err != nil {
			hlog.FromRequest(r).Err(err).Msg("error on write lib.js")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	if err := http.ListenAndServe(net.JoinHostPort(viper.GetString("hostname"), strconv.FormatUint(uint64(viper.GetUint16("web-port")), 10)), s); err != nil {
		log.Fatal().Err(err).Msg("Startup failed")
	}
}

func WeylusWebsocket(logger zerolog.Logger) {
	s := http.NewServeMux()
	s.Handle("/", http.HandlerFunc(HandleWebsocket))
	if err := http.ListenAndServe(":9001", s); err != nil {
		log.Fatal().Err(err).Msg("Startup failed")
	}
}

func HandleWebsite(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/html")
	d := getBaseConfig()
	d.WebsocketPort = viper.GetUint16("websocket-port")
	d.LogLevel = int(zerolog.GlobalLevel())

	authed := false
	if accessCode := viper.GetString("access-code"); accessCode != "" {
		if code := r.URL.Query().Get("access_code"); code != "" {
			d.AccessCode = code
			authed = true
			hlog.FromRequest(r).Debug().Msg("web client authenticated")
		}
	} else {
		authed = true
	}

	if !authed {
		w.Header().Add("Content-Type", "text/html")
		if _, err := w.Write([]byte(web.AccessHTML)); err != nil {
			hlog.FromRequest(r).Err(err).Msg("error on write access_code.html")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	tmpl, err := template.New("IndexHTML").Parse(web.IndexHTML)
	if err != nil {
		hlog.FromRequest(r).Err(err).Msg("error on parse template")
		w.WriteHeader(http.StatusInternalServerError)
	}
	if err := tmpl.Execute(w, d); err != nil {
		hlog.FromRequest(r).Err(err).Msg("error on execute template")
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		hlog.FromRequest(r).Err(err).Msg("error on accept websocket")
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*100)
	defer cancel()
	for {
		_, reader, _ := c.Reader(ctx)
		io.Copy(os.Stdout, reader)
	}
	c.Close(websocket.StatusNormalClosure, "")
}
