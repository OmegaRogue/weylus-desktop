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
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package server

import (
	"context"
	_ "embed"
	"html/template"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/OmegaRogue/weylus-desktop/utils"
	"github.com/OmegaRogue/weylus-desktop/web"
	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type data struct {
	AccessCode           string
	WebsocketPort        uint16
	LogLevel             int
	UInputEnabled        bool
	CaptureCursorEnabled bool
}

type WeylusServer struct {
	websiteAddr     string
	websocketAddr   string
	msgs            chan utils.Msg
	websiteServer   *http.Server
	websocketServer *http.Server
}

func newWeylusWebsiteServer(ctx context.Context, logger zerolog.Logger, addr string, websocketPort uint16) *http.Server {
	mux := http.NewServeMux()
	c := middleware(logger)
	h := c.Then(http.HandlerFunc(handleWebsite(websocketPort)))
	mux.Handle("/", h)
	mux.Handle("/style.css", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		if _, err := w.Write([]byte(web.StyleCSS)); err != nil {
			hlog.FromRequest(r).Err(err).Msg("error on write style.css")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	mux.Handle("/access_code.html", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		if _, err := w.Write([]byte(web.AccessHTML)); err != nil {
			hlog.FromRequest(r).Err(err).Msg("error on write access_code.html")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	mux.Handle("/lib.js", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/javascript")
		if _, err := w.Write([]byte(web.LibJS)); err != nil {
			hlog.FromRequest(r).Err(err).Msg("error on write lib.js")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	return &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 10 * time.Second,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
		Handler: mux,
	}
}
func newWeylusWebsocketServer(ctx context.Context, logger zerolog.Logger, addr string) *http.Server {
	mux := http.NewServeMux()
	c := middleware(logger)
	h := c.Then(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		c, err := websocket.Accept(writer, request, nil)
		if err != nil {
			hlog.FromRequest(request).Err(err).Msg("error on accept websocket")
		}

		defer func(c *websocket.Conn, code websocket.StatusCode, reason string) {
			err := c.Close(code, reason)
			if err != nil {
				hlog.FromRequest(request).Err(err).Msg("error on close websocket")
			}
		}(c, websocket.StatusInternalError, "the sky is falling")
		ctx, cancel := context.WithTimeout(request.Context(), time.Second*100)
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				writer.WriteHeader(http.StatusInternalServerError)
				if err := c.Close(websocket.StatusInternalError, "the sky is falling"); err != nil {
					hlog.FromRequest(request).Err(err).Msg("error on close websocket")
				}
				return
			default:
				var v interface{}
				err = wsjson.Read(ctx, c, &v)
				if err != nil {
					log.Fatal().Err(err).Msg("read")
				}
				log.Info().Msgf("received: %v", v)
			}
		}
	}))
	mux.Handle("/", h)
	return &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 10 * time.Second,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
		Handler: mux,
	}
}

func NewWeylusServer(ctx context.Context, hostname string, websitePort, websocketPort uint16) *WeylusServer {
	s := new(WeylusServer)
	s.msgs = make(chan utils.Msg)
	logger := log.With().Str("component", "server").Logger()
	ctx = logger.WithContext(ctx)
	s.websiteAddr = net.JoinHostPort(hostname, strconv.FormatUint(uint64(websitePort), 10))
	s.websocketAddr = net.JoinHostPort(hostname, strconv.FormatUint(uint64(websocketPort), 10))
	s.websiteServer = newWeylusWebsiteServer(ctx, logger, s.websiteAddr, websocketPort)
	s.websocketServer = newWeylusWebsocketServer(ctx, logger, s.websocketAddr)
	return s
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

func (s *WeylusServer) RunWebsite() {
	if err := s.websiteServer.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("website failed")
	}
}
func (s *WeylusServer) RunWebsocket() {
	if err := s.websocketServer.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("websocket failed")
	}
}

func handleWebsite(websocketPort uint16) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "text/html")
		d := getBaseConfig()
		d.WebsocketPort = websocketPort

		//authed := false
		//if accessCode := viper.GetString("access-code"); accessCode != "" {
		//	if code := r.URL.Query().Get("access_code"); code != "" {
		//		d.AccessCode = code
		//		authed = true
		//		hlog.FromRequest(r).Debug().Msg("web client authenticated")
		//	}
		//} else {
		//	authed = true
		//}
		//
		//if !authed {
		//	w.Header().Add("Content-Type", "text/html")
		//	if _, err := w.Write([]byte(web.AccessHTML)); err != nil {
		//		hlog.FromRequest(r).Err(err).Msg("error on write access_code.html")
		//		w.WriteHeader(http.StatusInternalServerError)
		//	}
		//	return
		//}

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
}
