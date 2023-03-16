//go:generate go-enum --marshal --names --values
package protocol

// PointerEventType indicates what kind of event is represented.
/*
ENUM(
down=pointerdown
up=pointerup
cancel=pointercancel
move=pointermove
)
*/
type PointerEventType string
