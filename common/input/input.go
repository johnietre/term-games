package input

type Input uint32

const (
	Unknown               Input = 0x00
	CtrlC                 Input = 0x03
	Backspace             Input = 0x08
	NewLine, LineFeed     Input = 0x0A, 0x0A
	Enter, CarriageReturn Input = 0x0D, 0x0D
	Escape                Input = 0x1B
	Del                   Input = 0x7F
	// TODO
	CtrlC2     Input = 0x00_1B_5B_03
	ArrowUp    Input = 0x00_1B_5B_41
	ArrowDown  Input = 0x00_1B_5B_42
	ArrowRight Input = 0x00_1B_5B_43
	ArrowLeft  Input = 0x00_1B_5B_44
	Del2       Input = 0x1B_5B_33_7E
)

func FromBytes(b []byte) Input {
	switch b[0] {
	case 0x1B:
		switch b[1] {
		case 0x5B:
			switch b[2] {
			case 0x03:
				return CtrlC2
			case 0x33:
				switch b[3] {
				case 0x7E:
					return Del2
				}
			case 0x41:
				return ArrowUp
			case 0x42:
				return ArrowDown
			case 0x43:
				return ArrowRight
			case 0x44:
				return ArrowLeft
			}
		}
	}
	return Input(b[0])
}

func (i Input) Byte() byte {
	return byte(rune(i))
}

func (i Input) Rune() rune {
	return rune(i)
}

func (i Input) IsAsciiAlpha() bool {
	return (i >= Input('A') && i <= Input('Z')) ||
		(i >= Input('a') && i <= Input('z'))
}

func (i Input) IsAsciiUpper() bool {
	return i >= Input('A') && i <= Input('Z')
}

func (i Input) IsAsciiLower() bool {
	return i >= Input('a') && i <= Input('z')
}

func (i Input) IsAsciiDigit() bool {
	return i >= Input('0') && i <= Input('9')
}

func (i Input) ToAsiiUpper() Input {
	if i.IsAsciiLower() {
		i -= Input('a' - 'A')
	}
	return i
}
