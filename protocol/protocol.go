package protocol

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

const (
	CommandTryGetFrame       = "TryGetFrame"
	CommandGetCapturableList = "GetCapturableList"
)

type Config struct {
	UInputSupport bool   `json:"uinput_support"`
	CapturableID  uint   `json:"capturable_id"`
	CaptureCursor bool   `json:"capture_cursor"`
	MaxWidth      uint   `json:"max_width"`
	MaxHeight     uint   `json:"max_height"`
	ClientName    string `json:"client_name,omitempty"`
}

type MessageOutboundContent interface {
	PointerEvent | WheelEvent | KeyboardEvent | Config
}
type MessageOutbound interface {
	map[string]PointerEvent | map[string]WheelEvent | map[string]KeyboardEvent | map[string]Config | ~string
}

func WrapMessage[T MessageOutboundContent | ~string](a T) any {
	wrapper := make(map[string]T)
	switch any(a).(type) {
	case PointerEvent:
		wrapper["PointerEvent"] = a
	case WheelEvent:
		wrapper["WheelEvent"] = a
	case KeyboardEvent:
		wrapper["KeyboardEvent"] = a
	case Config:
		wrapper["Config"] = a
	case string:
		return a
	default:
		log.Fatal().Msg("invalid type")
	}
	return wrapper
}

type WeylusError struct {
	ErrorMessage string `json:"Error"`
}
type WeylusConfigError struct {
	ErrorMessage string `json:"ConfigError"`
}

func (e *WeylusError) Error() string {
	return fmt.Sprintf("WeylusError: %s", e.ErrorMessage)
}
func (e *WeylusConfigError) Error() string {
	return fmt.Sprintf("WeylusConfigError: %s", e.ErrorMessage)
}

type PointerEvent struct {
	EventType   PointerEventType `json:"event_type"`
	PointerId   int              `json:"pointer_id"`
	Timestamp   uint64           `json:"timestamp"`
	IsPrimary   bool             `json:"is_primary"`
	PointerType PointerType      `json:"pointer_type"`
	Button      ButtonFlags      `json:"button"`
	Buttons     ButtonFlags      `json:"buttons"`
	X           float64          `json:"x"`
	Y           float64          `json:"y"`
	MovementX   int64            `json:"movement_x"`
	MovementY   int64            `json:"movement_y"`
	Pressure    float64          `json:"pressure"`
	TiltX       int32            `json:"tilt_x"`
	TiltY       int32            `json:"tilt_y"`
	Twist       int32            `json:"twist"`
	Width       float64          `json:"width"`
	Height      float64          `json:"height"`
}

type WheelEvent struct {
	Dx        int32  `json:"dx"`
	Dy        int32  `json:"dy"`
	Timestamp uint64 `json:"timestamp"`
}

type KeyboardEvent struct {
	EventType KeyboardEventType `json:"event_type"`
	Code      string            `json:"code"`
	Key       string            `json:"key"`
	Location  KeyboardLocation  `json:"location"`
	Alt       bool              `json:"alt"`
	Ctrl      bool              `json:"ctrl"`
	Shift     bool              `json:"shift"`
	Meta      bool              `json:"meta"`
}
