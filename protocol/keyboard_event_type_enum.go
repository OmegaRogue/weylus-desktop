// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package protocol

import (
	"fmt"
	"strings"
)

const (
	// KeyboardEventTypeDown is a KeyboardEventType of type down.
	KeyboardEventTypeDown KeyboardEventType = "down"
	// KeyboardEventTypeUp is a KeyboardEventType of type up.
	KeyboardEventTypeUp KeyboardEventType = "up"
	// KeyboardEventTypeRepeat is a KeyboardEventType of type repeat.
	KeyboardEventTypeRepeat KeyboardEventType = "repeat"
)

var ErrInvalidKeyboardEventType = fmt.Errorf("not a valid KeyboardEventType, try [%s]", strings.Join(_KeyboardEventTypeNames, ", "))

var _KeyboardEventTypeNames = []string{
	string(KeyboardEventTypeDown),
	string(KeyboardEventTypeUp),
	string(KeyboardEventTypeRepeat),
}

// KeyboardEventTypeNames returns a list of possible string values of KeyboardEventType.
func KeyboardEventTypeNames() []string {
	tmp := make([]string, len(_KeyboardEventTypeNames))
	copy(tmp, _KeyboardEventTypeNames)
	return tmp
}

// KeyboardEventTypeValues returns a list of the values for KeyboardEventType
func KeyboardEventTypeValues() []KeyboardEventType {
	return []KeyboardEventType{
		KeyboardEventTypeDown,
		KeyboardEventTypeUp,
		KeyboardEventTypeRepeat,
	}
}

// String implements the Stringer interface.
func (x KeyboardEventType) String() string {
	return string(x)
}

// String implements the Stringer interface.
func (x KeyboardEventType) IsValid() bool {
	_, err := ParseKeyboardEventType(string(x))
	return err == nil
}

var _KeyboardEventTypeValue = map[string]KeyboardEventType{
	"down":   KeyboardEventTypeDown,
	"up":     KeyboardEventTypeUp,
	"repeat": KeyboardEventTypeRepeat,
}

// ParseKeyboardEventType attempts to convert a string to a KeyboardEventType.
func ParseKeyboardEventType(name string) (KeyboardEventType, error) {
	if x, ok := _KeyboardEventTypeValue[name]; ok {
		return x, nil
	}
	return KeyboardEventType(""), fmt.Errorf("%s is %w", name, ErrInvalidKeyboardEventType)
}

// MarshalText implements the text marshaller method.
func (x KeyboardEventType) MarshalText() ([]byte, error) {
	return []byte(string(x)), nil
}

// UnmarshalText implements the text unmarshaller method.
func (x *KeyboardEventType) UnmarshalText(text []byte) error {
	tmp, err := ParseKeyboardEventType(string(text))
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}
