package protocol

type ButtonFlags byte

const (
	// ButtonNone is a ButtonFlags of type None.
	ButtonNone ButtonFlags = 0
	// ButtonPrimary is a ButtonFlags of type Primary. Usually the left button
	ButtonPrimary ButtonFlags = 1 << (iota - 1)
	// ButtonSecondary is a ButtonFlags of type Secondary. Usually the right button
	ButtonSecondary
	// ButtonAuxiliary is a ButtonFlags of type Auxiliary. Usually the wheel button or the middle button (if present)
	ButtonAuxiliary
	// ButtonFourth is a ButtonFlags of type Fourth. Typically the Browser Back button
	ButtonFourth
	// ButtonFifth is a ButtonFlags of type Fifth. Typically the Browser Forward button
	ButtonFifth
	// ButtonEraser is a ButtonFlags of type Eraser.
	ButtonEraser
)
