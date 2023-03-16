//go:generate go-enum --marshal --names --values
package protocol

// KeyboardEventType indicates what kind of event is represented.
// ENUM(down,up,repeat)
type KeyboardEventType string
