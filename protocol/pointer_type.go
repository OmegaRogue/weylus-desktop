//go:generate go-enum --marshal --names --values

package protocol

/*
ENUM(
unknown=""
mouse
pen
touch
)
*/
type PointerType string
