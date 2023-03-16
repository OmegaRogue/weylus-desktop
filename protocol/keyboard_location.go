//go:generate go-enum --marshal --names --values
package protocol

// KeyboardLocation identifies which part of the keyboard the key event originates from.
/*
 ENUM(
 standard // The key described by the event is not identified as being located in a particular area of the keyboard.
 left // The key is on the left side of the keyboard.
 right // The key is located on the right side of the keyboard.
 numpad // The key is located on the numeric keypad.
)
*/
type KeyboardLocation int
