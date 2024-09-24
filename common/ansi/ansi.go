package ansi

import "strconv"

// Moves cursor to home position (0, 0)
const (
	CurToHome  = "\x1b[H"
	CurUp      = "\x1b[1A"
	CurDown    = "\x1b[1B"
	CurRight   = "\x1b[1C"
	CurLeft    = "\x1b[1D"
	CurSave    = "\x1b 7"
	CurRestore = "\x1b 8"
	CurPos     = "\x1b[6n"
)

func CurUpN(n int) string {
	return "\x1b[" + strconv.Itoa(n) + "A"
}

func CurDownN(n int) string {
	return "\x1b[" + strconv.Itoa(n) + "B"
}

func CurRightN(n int) string {
	return "\x1b[" + strconv.Itoa(n) + "C"
}
func CurLeftN(n int) string {
	return "\x1b[" + strconv.Itoa(n) + "D"
}

// CUR_{direction}_LINE moves cursor to the beg of line N lines up/down
const (
	CurDownLine = "\x1b[1E"
	CurUpLine   = "\x1b[1F"
)

func CurDownLineN(n int) string {
	return "\x1b[" + strconv.Itoa(n) + "E"
}
func CurUpLineN(n int) string {
	return "\x1b[" + strconv.Itoa(n) + "F"
}

func CurToColN(n int) string {
	return "\x1b[" + strconv.Itoa(n) + "G"
}
func CurToPos(x, y int) string {
	//format!("\x1b[{};{}H", y, x)
	return "\x1b[" + strconv.Itoa(y) + ";" + strconv.Itoa(x) + "G"
}

// Figure out difference between "H" and "f"
// Move cursor to line Y, col X
// "\x1b[{Y};{X}H"
// "\x1b[{Y};{x}f"

// TODO: Figure out difference between "\x1b[J" and "\x1b[2J"
const (
	ClearScreenRight = "\x1b[0J"
	ClearScreenLeft  = "\x1b[1J"
	ClearScreen      = "\x1b[2J"
)

// TODO: Figure out difference between "\x1b[K" and "\x1b[2K"
const (
	ClearLineRight = "\x1b[0K"
	ClearLineLeft  = "\x1b[1K"
	ClearLine      = "\x1b[2K"
)

const (
	ResetAll = "\x1b[0m"

	SetBold          = "\x1b[1m"
	SetDim           = "\x1b[2m"
	SetItalic        = "\x1b[3m"
	SetUnderline     = "\x1b[4m"
	SetBlinking      = "\x1b[5m"
	SetInverse       = "\x1b[7m"
	SetInvisible     = "\x1b[8m"
	SetStrikethrough = "\x1b[9m"

	ResetBold          = "\x1b[22m"
	ResetDim           = "\x1b[22m"
	ResetItalic        = "\x1b[23m"
	ResetUnderline     = "\x1b[24m"
	ResetBlinking      = "\x1b[25m"
	ResetInverse       = "\x1b[27m"
	ResetInvisible     = "\x1b[28m"
	ResetStrikethrough = "\x1b[29m"
)

// FORE = foreground color; BACK = background color
const (
	ForeBlack   = "\x1b[30m"
	ForeRed     = "\x1b[31m"
	ForeGreen   = "\x1b[32m"
	ForeYellow  = "\x1b[33m"
	ForeBlue    = "\x1b[34m"
	ForeMagenta = "\x1b[35m"
	ForeCyan    = "\x1b[36m"
	ForeWhite   = "\x1b[37m"
	ForeDefault = "\x1b[39m"
)

const (
	BackBlack   = "\x1b[40m"
	BackRed     = "\x1b[41m"
	BackGreen   = "\x1b[42m"
	BackYellow  = "\x1b[43m"
	BackBlue    = "\x1b[44m"
	BackMagenta = "\x1b[45m"
	BackCyan    = "\x1b[46m"
	BackWhite   = "\x1b[47m"
	BackDefault = "\x1b[49m"
)

func SetForeColorRgb(r, g, b int) string {
	//format!("\x1b[38;2;{};{};{}m", r, g, b)
	return "\x1b[38;2;" +
		strconv.Itoa(r) + ";" +
		strconv.Itoa(g) + ";" +
		strconv.Itoa(b) + "m"
}

func SetBackColorRgb(r, g, b int) string {
	//format!("\x1b[48;2;{};{};{}m", r, g, b)
	return "\x1b[48;2;" +
		strconv.Itoa(r) + ";" +
		strconv.Itoa(g) + ";" +
		strconv.Itoa(b) + "m"
}

/*
#[repr(u8)]
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
enum GraphicsMode {
    Reset(u8) = 0,
    Bold = 1,
    Dim = 2,
    Italic = 3,
    Underline = 4,
    Blinking = 5,
    Inverse = 7,
    Invisible = 8,
    Strikethrough = 9,
}
*/

/*
impl GraphicsMode {
    func as_u8(self) -> u8 {
        use self::GraphicsMode::*
        match self {
            Bold => 1,
            Dim => 2,
            Italic => 3,
            Underline => 4,
            Blinking => 5,
            Inverse => 7,
            Invisible => 8,
            Strikethrough => 9,
            Reset(code) => code,
        }
    }

    const func reset_code(self) -> u8 {
        use self::GraphicsMode::*
        match self {
            Bold | Dim => 22,
            Italic => 23,
            Underline => 24,
            Blinking => 25,
            Inverse => 27,
            Invisible => 28,
            Strikethrough => 29,
            Reset(code) => code,
        }
    }

    const func reset(self) -> Self {
        Self::Reset(self.reset_code())
    }

    const func const_str(self) -> &'static str {
        use self::GraphicsMode::*
        match self {
            Bold => RESET_BOLD,
            Dim => RESET_DIM,
            Italic => RESET_ITALIC,
            Underline => RESET_UNDERLINE,
            Blinking => RESET_BLINKING,
            Inverse => RESET_INVERSE,
            Invisible => RESET_INVISIBLE,
            Strikethrough => RESET_STRIKETHROUGH,
            Reset(code) if code == Bold.reset_code() => RESET_BOLD,
            Reset(code) if code == Dim.reset_code() => RESET_DIM,
            Reset(code) if code == Italic.reset_code() => RESET_ITALIC,
            Reset(code) if code == Underline.reset_code() => RESET_UNDERLINE,
            Reset(code) if code == Blinking.reset_code() => RESET_BLINKING,
            Reset(code) if code == Inverse.reset_code() => RESET_INVERSE,
            Reset(code) if code == Invisible.reset_code() => RESET_INVISIBLE,
            Reset(code) if code == Strikethrough.reset_code() => RESET_STRIKETHROUGH,
            Reset(_) => panic!("unknown reset code"),
        }
    }
}

impl fmt::Display for GraphicsMode {
    func fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        //write!(f, "{}", self.const_str())
        //write!(f, "\x1b[{}m", *self as u8)
        write!(f, "\x1b[{}m", self.as_u8())
    }
}
*/

/*
#[repr(u8)]
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
enum ForeColor {
    RGB(u8, u8, u8),
    Black = 30,
    Red = 31,
    Green = 32,
    Yellow = 33,
    Blue = 34,
    Magenta = 35,
    Cyan = 36,
    White = 37,
    Default = 39,
}

impl ForeColor {
    func as_back(self) -> BackColor {
        use self::ForeColor::*
        match self {
            Black => BackColor::Black,
            Red => BackColor::Red,
            Green => BackColor::Green,
            Yellow => BackColor::Yellow,
            Blue => BackColor::Blue,
            Magenta => BackColor::Magenta,
            Cyan => BackColor::Cyan,
            White => BackColor::White,
            Default => BackColor::Default,
            RGB(r, g, b) => BackColor::RGB(r, g, b),
        }
    }

    func as_u8(self) -> u8 {
        use self::ForeColor::*
        match self {
            Black => 30,
            Red => 31,
            Green => 32,
            Yellow => 33,
            Blue => 34,
            Magenta => 35,
            Cyan => 36,
            White => 37,
            Default => 39,
            RGB(_, _, _) => 39,
        }
    }

    const func const_str(self) -> &'static str {
        use self::ForeColor::*
        match self {
            Black => FORE_BLACK,
            Red => FORE_RED,
            Green => FORE_GREEN,
            Yellow => FORE_YELLOW,
            Blue => FORE_BLUE,
            Magenta => FORE_MAGENTA,
            Cyan => FORE_CYAN,
            White => FORE_WHITE,
            Default => FORE_DEFAULT,
        }
    }
}

impl fmt::Display for ForeColor {
    func fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        //write!(f, "{}", self.const_str())
        if let ForeColor::RGB(r, b, g) = self {
            write!(f, "\x1b[38;2;{};{};{}m", r, g, b)
        } else {
            write!(f, "\x1b[{}m", self.as_u8())
        }
    }
}

#[repr(u8)]
#[derive(Clone, Copy, Debug, PartialEq, Eq)]
enum BackColor {
    RGB(u8, u8, u8),
    Black = 40,
    Red = 41,
    Green = 42,
    Yellow = 43,
    Blue = 44,
    Magenta = 45,
    Cyan = 46,
    White = 47,
    Default = 49,
}

impl BackColor {
    func as_fore(self) -> ForeColor {
        use self::BackColor::*
        match self {
            Black => ForeColor::Black,
            Red => ForeColor::Red,
            Green => ForeColor::Green,
            Yellow => ForeColor::Yellow,
            Blue => ForeColor::Blue,
            Magenta => ForeColor::Magenta,
            Cyan => ForeColor::Cyan,
            White => ForeColor::White,
            Default => ForeColor::Default,
            RGB(r, g, b) => ForeColor::RGB(r, g, b),
        }
    }

    func as_u8(self) -> u8 {
        use self::BackColor::*
        match self {
            Black => 40,
            Red => 41,
            Green => 42,
            Yellow => 43,
            Blue => 44,
            Magenta => 45,
            Cyan => 46,
            White => 47,
            Default => 49,
            RGB(_, _, _) => 49,
        }
    }

    const func const_str(self) -> &'static str {
        use self::BackColor::*
        match self {
            Black => BACK_BLACK,
            Red => BACK_RED,
            Green => BACK_GREEN,
            Yellow => BACK_YELLOW,
            Blue => BACK_BLUE,
            Magenta => BACK_MAGENTA,
            Cyan => BACK_CYAN,
            White => BACK_WHITE,
            Default => BACK_DEFAULT,
        }
    }
}

impl fmt::Display for BackColor {
    func fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        //write!(f, "{}", self.const_str())
        if let BackColor::RGB(r, b, g) = self {
            write!(f, "\x1b[48;2;{};{};{}m", r, g, b)
        } else {
            write!(f, "\x1b[{}m", self.as_u8())
        }
    }
}
*/

/*
func SET_GRAPHICS(graphics: GraphicsMode) string {
    format!("\x1b[{}m", graphics.as_u8())
}

func WRAP_GRAPHICS(text, graphics: GraphicsMode) string {
    format!(
        "\x1b[{}m{}{}",
        graphics.as_u8(),
        text,
        graphics.reset_code()
    )
}

func SET_FORE_COLOR_ID(color: ForeColor) string {
    format!("\x1b[{}m", color.as_u8())
}

func WRAP_FORE_COLOR_ID(text, color: ForeColor) string {
    format!("\x1b[{}m{}{}", color.as_u8(), text, FORE_DEFAULT)
}

func SET_BACK_COLOR_ID(color: BackColor) string {
    format!("\x1b[{}m", color.as_u8())
}

func WRAP_BACK_COLOR_ID(text, color: BackColor) string {
    format!("\x1b[{}m{}{}", color.as_u8(), text, BACK_DEFAULT)
}
*/

/*
func WRAP_FORE_COLOR_RGB(text, (r, g, b): (u8, u8, u8)) string {
    format!("\x1b[38;2;{};{};{}m{}{}", r, g, b, text, FORE_DEFAULT)
}
*/

/*
func WRAP_BACK_COLOR_RGB(text, (r, g, b): (u8, u8, u8)) string {
    format!("\x1b[48;2;{};{};{}m{}{}", r, g, b, text, BACK_DEFAULT)
}
*/

/*
#[derive(Clone, Debug, PartialEq, Eq)]
enum TextNode {
    Text(String),
    Fore(ForeColor),
    Back(BackColor),
    Graphics(GraphicsMode),
}

impl fmt::Display for TextNode {
    func fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            TextNode::Text(text) => write!(f, "{}", text),
            TextNode::Fore(color) => write!(f, "{}", color),
            TextNode::Back(color) => write!(f, "{}", color),
            TextNode::Graphics(graphics) => write!(f, "{}", graphics),
        }
    }
}

#[derive(Default)]
struct Text {
    nodes: Vec<TextNode>,
}

impl Text {
    func new() -> Self {
        Default::default()
    }

    func text(&mut self, text: impl ToString) -> &mut Self {
        self.nodes.push(TextNode::Text(text.to_string()))
        self
    }

    func graphics(&mut self, graphics: GraphicsMode) -> &mut Self {
        self.nodes.push(TextNode::Graphics(graphics))
        self
    }

    func wrap_graphics(&mut self, text: impl ToString, graphics: GraphicsMode) -> &mut Self {
        self.nodes.extend_from_slice(&[
            TextNode::Graphics(graphics),
            TextNode::Text(text.to_string()),
            TextNode::Graphics(graphics.reset()),
        ])
        self
    }

    func reset_graphics(&mut self, graphics: GraphicsMode) -> &mut Self {
        self.nodes.push(TextNode::Graphics(graphics.reset()))
        self
    }

    func fore(&mut self, color: ForeColor) -> &mut Self {
        self.nodes.push(TextNode::Fore(color))
        self
    }

    func wrap_fore(&mut self, text: impl ToString, color: ForeColor) -> &mut Self {
        self.nodes.extend_from_slice(&[
            TextNode::Fore(color),
            TextNode::Text(text.to_string()),
            TextNode::Fore(ForeColor::Default),
        ])
        self
    }

    func reset_fore(&mut self) -> &mut Self {
        self.nodes.push(TextNode::Fore(ForeColor::Default))
        self
    }

    func back(&mut self, color: BackColor) -> &mut Self {
        self.nodes.push(TextNode::Back(color))
        self
    }

    func wrap_back(&mut self, text: impl ToString, color: BackColor) -> &mut Self {
        self.nodes.extend_from_slice(&[
            TextNode::Back(color),
            TextNode::Text(text.to_string()),
            TextNode::Back(BackColor::Default),
        ])
        self
    }

    func reset_back(&mut self) -> &mut Self {
        self.nodes.push(TextNode::Back(BackColor::Default))
        self
    }

    func reset(&mut self) {
        self.nodes.clear()
    }
}

impl Add for Text {
    type Output = Text

    func add(mut self, rhs: Text) -> Self::Output {
        self.nodes.extend(rhs.nodes)
        self
    }
}

impl Add<GraphicsMode> for Text {
    type Output = Text

    func add(mut self, rhs: GraphicsMode) -> Self::Output {
        self.graphics(rhs)
        self
    }
}

impl Add<ForeColor> for Text {
    type Output = Text

    func add(mut self, rhs: ForeColor) -> Self::Output {
        self.fore(rhs)
        self
    }
}

impl Add<BackColor> for Text {
    type Output = Text

    func add(mut self, rhs: BackColor) -> Self::Output {
        self.back(rhs)
        self
    }
}

impl fmt::Display for Text {
    func fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        self.nodes.iter().try_for_each(|node| write!(f, "{}", node))
    }
}

#[derive(Default)]
struct TextBuilder {
    text: Text,
}

impl TextBuilder {
    func new() -> Self {
        Default::default()
    }

    func text(mut self, text: impl ToString) -> Self {
        self.text.nodes.push(TextNode::Text(text.to_string()))
        self
    }

    func graphics(mut self, graphics: GraphicsMode) -> Self {
        self.text.nodes.push(TextNode::Graphics(graphics))
        self
    }

    func wrap_graphics(mut self, text: impl ToString, graphics: GraphicsMode) -> Self {
        self.text.nodes.extend_from_slice(&[
            TextNode::Graphics(graphics),
            TextNode::Text(text.to_string()),
            TextNode::Graphics(graphics.reset()),
        ])
        self
    }

    func reset_graphics(mut self, graphics: GraphicsMode) -> Self {
        self.text.nodes.push(TextNode::Graphics(graphics.reset()))
        self
    }

    func fore(mut self, color: ForeColor) -> Self {
        self.text.nodes.push(TextNode::Fore(color))
        self
    }

    func wrap_fore(mut self, text: impl ToString, color: ForeColor) -> Self {
        self.text.nodes.extend_from_slice(&[
            TextNode::Fore(color),
            TextNode::Text(text.to_string()),
            TextNode::Fore(ForeColor::Default),
        ])
        self
    }

    func reset_fore(mut self) -> Self {
        self.text.nodes.push(TextNode::Fore(ForeColor::Default))
        self
    }

    func back(mut self, color: BackColor) -> Self {
        self.text.nodes.push(TextNode::Back(color))
        self
    }

    func wrap_back(mut self, text: impl ToString, color: BackColor) -> Self {
        self.text.nodes.extend_from_slice(&[
            TextNode::Back(color),
            TextNode::Text(text.to_string()),
            TextNode::Back(BackColor::Default),
        ])
        self
    }

    func reset_back(mut self) -> Self {
        self.text.nodes.push(TextNode::Back(BackColor::Default))
        self
    }

    func build(self) -> Text {
        self.text
    }
}
*/

const (
	CurInvisible  = "\x1b[?25l"
	CurVisible    = "\x1b[?25h"
	ScreenRestore = "\x1b[?47l"
	ScreenSave    = "\x1b[?47h"
)
