//go:generate go-enum --marshal --names --values
package protocol

/*
ENUM(
down=pointerdown
up=pointerup
cancel=pointercancel
move=pointermove
)
*/
type PointerEventType string
