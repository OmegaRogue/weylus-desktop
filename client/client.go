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
package client

import (
	"bufio"
	"context"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/OmegaRogue/weylus-desktop/protocol"
	"github.com/OmegaRogue/weylus-desktop/utils"
	"github.com/kr/pretty"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var WebsocketNotStartedError = errors.New("Websocket not initialized")

type Callback func(msg Msg)
type WeylusClient struct {
	ws                    *websocket.Conn
	msgs                  chan Msg
	callbacks             map[protocol.WeylusResponse][]Callback
	callbackMutex         sync.Mutex
	ctx                   context.Context
	cancel                context.CancelFunc
	Framerate             uint
	frameTimer            *time.Ticker
	BufPipe               *bufio.ReadWriter
	receivedVideoResponse bool
	requestedFirstFrame   bool
}

func (w *WeylusClient) AddCallback(event protocol.WeylusResponse, callback Callback) int {
	w.callbackMutex.Lock()
	defer w.callbackMutex.Unlock()
	n := len(w.callbacks[event])
	w.callbacks[event] = append(w.callbacks[event], callback)
	return n
}
func (w *WeylusClient) RemoveCallback(event protocol.WeylusResponse, i int) int {
	w.callbackMutex.Lock()
	defer w.callbackMutex.Unlock()
	n := len(w.callbacks[event])

	w.callbacks[event] = utils.Remove(w.callbacks[event], i)
	return n
}

func (w *WeylusClient) AddCallbackNext(event protocol.WeylusResponse, callback Callback) {
	var i int
	i = w.AddCallback(event, func(msg Msg) {
		defer w.RemoveCallback(event, i)
		callback(msg)
	})
}

func NewWeylusClient(ctx context.Context, fps uint) *WeylusClient {
	w := new(WeylusClient)
	w.msgs = make(chan Msg)
	w.callbacks = make(map[protocol.WeylusResponse][]Callback)
	ctx = log.With().Str("component", "client").Logger().WithContext(ctx)
	w.ctx, w.cancel = context.WithCancel(ctx)
	w.Framerate = fps
	log.Ctx(ctx).Info().Dur("frame_time", time.Second/time.Duration(w.Framerate)).Uint("fps", fps).Msg("video times")
	w.frameTimer = time.NewTicker(time.Second / time.Duration(w.Framerate))

	return w
}

func commandWithReceive[T protocol.MessageInbound, V protocol.MessageOutboundContent](w *WeylusClient, command V) (a T, err error) {
	if w.ws == nil {
		return a, errors.Wrap(WebsocketNotStartedError, "commandWithReceive failed")
	}
	var wg sync.WaitGroup
	wg.Add(1)
	w.AddCallbackNext(protocol.ResponseFromOutboundContent(command), func(msg Msg) {
		var r any
		r, err = protocol.ParseMessage(msg.Data)
		if b, ok := r.(T); ok {
			a = b
		} else {
			err = errors.Errorf("wrong type returned by ParseMessage: %v\n %v", reflect.TypeOf(r), pretty.Sprint(b))
		}
		wg.Done()
	})
	if err := wsjson.Write(w.ctx, w.ws, protocol.WrapMessage(command)); err != nil {
		return a, errors.Wrap(err, string(protocol.CommandFromOutboundContent(command)))
	}
	wg.Wait()
	if err != nil {
		return a, errors.Wrap(err, "parsing message")
	}
	return a, nil
}

func (w *WeylusClient) GetCapturableList() (protocol.CapturableList, error) {
	return commandWithReceive[protocol.CapturableList](w, protocol.WeylusCommandGetCapturableList)
}
func (w *WeylusClient) Config(config protocol.Config) (protocol.WeylusResponse, error) {
	resp, err := commandWithReceive[string](w, config)
	if err != nil {
		return protocol.WeylusResponseError, errors.Wrap(err, "error on receive")
	}
	return protocol.ParseWeylusResponse(resp)
}

func (w *WeylusClient) StartVideo() error {
	var wg sync.WaitGroup
	w.requestedFirstFrame = true
	wg.Add(1)
	w.AddCallbackNext(protocol.WeylusResponseNewVideo, func(msg Msg) {
		w.receivedVideoResponse = true
		log.Ctx(w.ctx).Info().Msg("video")
		wg.Done()
	})
	err := w.TryGetFrame()
	if err != nil {
		return errors.Wrap(err, "start video")
	}
	return nil
}

func (w *WeylusClient) TryGetFrame() error {
	if w.ws == nil {
		return errors.Wrap(WebsocketNotStartedError, "TryGetFrame failed")
	}
	if err := wsjson.Write(w.ctx, w.ws, protocol.WeylusCommandTryGetFrame); err != nil {
		return errors.Wrap(err, string(protocol.WeylusCommandTryGetFrame))
	}
	return nil
}
func (w *WeylusClient) SendPointerEvent(e protocol.PointerEvent) error {
	if w.ws == nil {
		return errors.Wrap(WebsocketNotStartedError, "SendPointerEvent failed")
	}
	if err := wsjson.Write(w.ctx, w.ws, e); err != nil {
		return errors.Wrap(err, string(protocol.WeylusCommandPointerEvent))
	}
	return nil
}
func (w *WeylusClient) SendWheelEvent(e protocol.WheelEvent) error {
	if w.ws == nil {
		return errors.Wrap(WebsocketNotStartedError, "SendWheelEvent failed")
	}
	if err := wsjson.Write(w.ctx, w.ws, e); err != nil {
		return errors.Wrap(err, string(protocol.WeylusCommandWheelEvent))
	}
	return nil
}
func (w *WeylusClient) SendKeyboardEvent(e protocol.KeyboardEvent) error {
	if w.ws == nil {
		return errors.Wrap(WebsocketNotStartedError, "SendKeyboardEvent failed")
	}
	if err := wsjson.Write(w.ctx, w.ws, e); err != nil {
		return errors.Wrap(err, string(protocol.WeylusCommandKeyboardEvent))
	}
	return nil
}

func (w *WeylusClient) Dial(address string) error {
	c, _, err := websocket.Dial(w.ctx, address, nil)
	if err != nil {
		return errors.Wrap(err, "dial weylusClient")
	}
	c.SetReadLimit(32769 * 16)
	w.ws = c
	return nil
}

type Msg struct {
	Type websocket.MessageType
	Data []byte
}

func (w *WeylusClient) Listen() {
	if w.ws == nil {
		log.Fatal().Msg("Listen failed")
	}
	for {
		select {
		case <-w.ctx.Done():
			log.Ctx(w.ctx).Err(errors.Wrap(w.ctx.Err(), "closed context")).Msg("closed context")
			return
		default:
			t, d, err := w.ws.Read(w.ctx)
			if err != nil {
				_ = w.Close()
				log.Ctx(w.ctx).Fatal().Err(err).Msg("error on listen")
				return
			}
			w.msgs <- Msg{
				Type: t,
				Data: d,
			}
		}
	}
}

func (w *WeylusClient) Close() error {
	w.cancel()
	close(w.msgs)
	if err := w.ws.Close(websocket.StatusNormalClosure, "closing"); err != nil {
		return errors.Wrap(err, "close websocket and channel")
	}
	return nil
}

func (w *WeylusClient) Run() {
	// avoid racecondition
	time.Sleep(time.Second)
	for {
		select {
		case <-w.ctx.Done():
			log.Ctx(w.ctx).Err(errors.Wrap(w.ctx.Err(), "closed context")).Msg("closed context")
			return
		case msg := <-w.msgs:
			switch msg.Type {
			case websocket.MessageText:
				log.Ctx(w.ctx).Info().RawJSON("data", msg.Data).Msg("received data")
				for response, callbacks := range w.callbacks {
					if strings.Contains(string(msg.Data), string(response)) {
						for _, callback := range callbacks {
							callback(msg)
						}
					}
				}
			case websocket.MessageBinary:
				if _, err := w.BufPipe.Write(msg.Data); err != nil {
					log.Ctx(w.ctx).Err(err).Msg("error on write data")
				}
			}
		}
	}
}

func (w *WeylusClient) RunVideo() {
	w.receivedVideoResponse = false
	w.requestedFirstFrame = false
	time.Sleep(time.Second * 2)
	for {
		select {
		case <-w.ctx.Done():
			w.frameTimer.Stop()
			log.Ctx(w.ctx).Err(errors.Wrap(w.ctx.Err(), "closed context")).Msg("closed context")
			return
		case <-w.frameTimer.C:
			if !w.requestedFirstFrame {
				if err := w.StartVideo(); err != nil {
					log.Ctx(w.ctx).Err(err).Msg("send TryGetFrame for first frame, dropped frame")
				}
				log.Ctx(w.ctx).Trace().Msg("tick")
				continue
			} else if !w.receivedVideoResponse {
				continue
			}
			if err := w.TryGetFrame(); err != nil {
				log.Ctx(w.ctx).Err(err).Msg("send TryGetFrame, dropped frame")
			}
			log.Ctx(w.ctx).Trace().Msg("tick")
		}
	}
}
