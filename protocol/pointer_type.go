//go:generate go-enum --marshal --names --values

package protocol

// PointerType indicates what type the pointer triggering this event is.
/*
ENUM(
unknown=""
mouse
pen
touch
)
*/
type PointerType string
