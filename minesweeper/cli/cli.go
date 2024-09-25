package cli

// TODO: Text to make sure works
// TODO: Custom sizes/num mines
// TODO: viewable area
// TODO: time
// TODO: scores (times)

import (
	"math/rand"
	"time"

	"github.com/johnietre/term-games/common"
	"github.com/johnietre/term-games/common/ansi"
	inputpkg "github.com/johnietre/term-games/common/input"
	"github.com/spf13/cobra"
)

func MakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "minesweeper",
		Short:                 "Play Minesweeper",
		Run:                   run,
		DisableFlagsInUseLine: true,
	}
	return cmd
}

func Run() {
	if err := MakeCmd().Execute(); err != nil {
		common.Fatal("error running: ", err)
	}
}

func run(_ *cobra.Command, _ []string) {
	defer common.GlobalDeferrer.Run()

	restore, _, err := common.MakeTermRaw()
	if err != nil {
		common.Fatal("error making term raw: ", err)
	}
	common.GlobalDefer(restore)

	game := NewGame()
	game.Run()
}

type Game struct {
	board *Board
	// Represents the viewable area
	// bottomRightView is exclusive
	topLeftView, bottomRightView Point
	// Top left is 0, 0
	pos Point

	tb *common.TermBuffer
}

func NewGame() *Game {
	return &Game{
		board:           NewBoard(9, 9),
		topLeftView:     Point{X: 0, Y: 0},
		bottomRightView: Point{X: 9, Y: 9},
		pos:             Point{X: 0, Y: 0},
		tb:              common.NewTermBuffer(nil),
	}
}

func (game *Game) Run() {
	game.runMenu()
	game.runBoard()
}

func (game *Game) runMenu() {
	game.tb.WriteString("Minesweeper!\n\r")
	game.tb.WriteString("1) Beginner (9x9/10)\n\r")
	game.tb.WriteString("2) Intermediate (16x16/40)\n\r")
	game.tb.WriteString("3) Advanced (30x16/99)\n\r")
	game.tb.WriteString("0) Exit\n\r")
	game.tb.Flush()
	for {
		buf, n, err := common.ReadStdinBytes()
		if err != nil {
			common.Fatal("error reading from stdin: ", err)
		} else if n == 0 {
			continue
		}
		switch input := inputpkg.FromByteArray(buf); input {
		case inputpkg.Input('1'):
			game.board = NewBoard(9, 9)
			game.RandomizeBoard(10)
		case inputpkg.Input('2'):
			game.board = NewBoard(16, 16)
			game.RandomizeBoard(40)
		case inputpkg.Input('3'):
			game.board = NewBoard(30, 16)
			game.RandomizeBoard(99)
		case inputpkg.Input('0'), inputpkg.CtrlC, inputpkg.Escape:
			common.Exit(0)
		default:
			continue
		}
		break
	}
	for i := 0; i < 5; i++ {
		game.tb.WriteString(ansi.CurUp + ansi.ClearLine)
	}
	game.tb.WriteByte('\r')
	game.tb.Flush()
}

func (game *Game) runBoard() {
	minesLeft := game.board.MinesLeft()
	game.tb.Print("*: ", minesLeft, "\n\r")
	game.DisplayBoard()
	game.tb.Flush()

	for {
		buf, n, err := common.ReadStdinBytes()
		if err != nil {
			common.Fatal("error reading from stdin: ", err)
		} else if n == 0 {
			continue
		}
		input := inputpkg.FromBytes(buf[:]).ToAsiiUpper()
		switch input {
		case inputpkg.ArrowUp, inputpkg.Input('W'), inputpkg.Input('K'):
			if game.pos.Y != 0 {
				game.pos.Y--
				game.tb.WriteString(ansi.CurUp)
			}
		case inputpkg.ArrowDown, inputpkg.Input('S'), inputpkg.Input('J'):
			if game.pos.Y != game.board.BottomRight().Y-1 {
				game.pos.Y++
				game.tb.WriteString(ansi.CurDown)
			}
		case inputpkg.ArrowRight, inputpkg.Input('L'), inputpkg.Input('D'):
			if game.pos.X != game.board.BottomRight().X-1 {
				game.pos.X++
				game.tb.WriteString(ansi.CurRight)
			}
		case inputpkg.ArrowLeft, inputpkg.Input('A'), inputpkg.Input('H'):
			if game.pos.X != 0 {
				game.pos.X--
				game.tb.WriteString(ansi.CurLeft)
			}
		case inputpkg.Input('F'):
			sq := game.board.Get(game.pos)
			if !sq.IsShowing() {
				if sq.IsFlagged() {
					game.board.Set(game.pos, sq.WithHidden())
				} else {
					game.board.Set(game.pos, sq.WithFlagged())
				}
				game.DisplayCurrSquare()
			}
		case inputpkg.Input('?'):
			sq := game.board.Get(game.pos)
			if !sq.IsShowing() {
				if sq.IsPossible() {
					game.board.Set(game.pos, sq.WithHidden())
				} else {
					game.board.Set(game.pos, sq.WithPossible())
				}
				game.DisplayCurrSquare()
			}
		case inputpkg.Input(' '):
			sq := game.board.Get(game.pos)
			if !sq.IsShowing() {
				if sq.IsHidden() {
					game.board.Set(game.pos, sq.WithFlagged())
				} else if sq.IsFlagged() {
					game.board.Set(game.pos, sq.WithPossible())
				} else {
					game.board.Set(game.pos, sq.WithHidden())
				}
				game.DisplayCurrSquare()
			}
		case inputpkg.Enter:
			changed := game.board.TryUncover(game.pos)
			if len(changed) == 0 {
				game.DisplayFailedBoard()
				game.MoveTo(game.board.BottomRight())
				game.tb.WriteByte('\r')
				game.tb.Flush()
				return
			}
			old := game.pos
			for _, pos := range changed {
				game.MoveTo(pos)
				game.DisplayCurrSquare()
			}
			game.MoveTo(old)
		case inputpkg.CtrlC, inputpkg.Escape:
			game.tb.WriteString(ansi.CurDownN(game.board.BottomRight().Y-game.pos.Y) + "\r")
			game.tb.Flush()
			common.Exit(0)
		}

		if game.board.AllUncovered() {
			game.MoveTo(NewPoint(0, 0))
			game.tb.WriteString(ansi.CurUp + "\r" + ansi.ClearLine + "SUCCESS\n\r")
			game.MoveTo(game.board.BottomRight())
			game.tb.WriteByte('\r')
			game.tb.Flush()
			break
		}
		if ml := game.board.MinesLeft(); ml != minesLeft {
			minesLeft = ml
			old := game.pos
			game.MoveTo(NewPoint(0, 0))
			game.tb.Print(ansi.CurUp+ansi.ClearLine+"\r"+"*: ", minesLeft, "\n\r")
			game.MoveTo(old)
		}

		game.tb.Flush()
	}
}

func (game *Game) NewBoard(width, height int) {
	game.board = NewBoard(width, height)
}

func (game *Game) RandomizeBoard(nmines int) {
	game.board.Randomize(nmines)
}

func (game *Game) MoveTo(pos Point) {
	if diff := pos.X - game.pos.X; diff > 0 {
		game.tb.WriteString(ansi.CurRightN(diff))
	} else if diff < 0 {
		game.tb.WriteString(ansi.CurLeftN(-diff))
	}
	if diff := pos.Y - game.pos.Y; diff > 0 {
		game.tb.WriteString(ansi.CurDownN(diff))
	} else if diff < 0 {
		game.tb.WriteString(ansi.CurUpN(-diff))
	}
	game.pos = pos
}

func (game *Game) DisplayBoard() {
	for _, row := range game.board.board {
		for _, sq := range row {
			//game.tb.WriteByte(sq.ToDisplayByte())
			//game.tb.WriteRune(sq.ToDisplayRune())
			game.tb.WriteString(sq.ToDisplayString())
		}
		game.tb.Write([]byte("\n\r"))
	}
	game.tb.WriteString(ansi.CurUpN(game.board.Height()) + "\r")
}

func (game *Game) DisplayFailedBoard() {
	old := game.pos
	game.MoveTo(NewPoint(0, 0))
	for h, row := range game.board.board {
		for w, sq := range row {
			if old.Y == h && old.X == w {
				game.tb.WriteString(ansi.BackRed)
			}
			//game.tb.WriteByte(sq.ToDisplayByte())
			//game.tb.WriteRune(sq.ToDisplayRune())
			game.tb.WriteString(sq.ToDisplayString())
			if old.Y == h && old.X == w {
				game.tb.WriteString(ansi.BackDefault)
			}
		}
		game.tb.Write([]byte("\n\r"))
	}
	game.tb.WriteString(ansi.CurUpN(game.board.Height()) + "\r")
	game.MoveTo(old)
}

func (game *Game) DisplayCurrSquare() {
	//game.tb.WriteRune(game.board.Get(game.pos).ToDisplayRune())
	game.tb.WriteString(game.board.Get(game.pos).ToDisplayString())
	game.tb.WriteString(ansi.CurLeft)
}

type Board struct {
	board               [][]Square
	mines, minesFlagged int
	uncovered           int
}

func NewBoard(width, height int) *Board {
	board := make([][]Square, height)
	for i := 0; i < height; i++ {
		board[i] = make([]Square, width)
	}
	return &Board{
		board: board,
	}
}

func (b *Board) Get(pos Point) Square {
	return b.board[pos.Y][pos.X]
}

func (b *Board) Set(pos Point, sq Square) {
	old := b.Get(pos)
	if old.IsFlagged() {
		if !sq.IsFlagged() {
			b.minesFlagged--
		}
	} else if sq.IsFlagged() {
		b.minesFlagged++
	}
	if !old.IsShowing() {
		if sq.IsShowing() {
			b.uncovered++
		}
	} else if !sq.IsShowing() {
		b.uncovered--
	}
	b.board[pos.Y][pos.X] = sq
}

func (b *Board) MinesLeft() int {
	return b.mines - b.minesFlagged
}

func (b *Board) AllUncovered() bool {
	return b.MinesLeft() == 0 && b.mines+b.uncovered == b.Area()
}

func (b *Board) Randomize(nmines int) {
	width, height := b.Width(), b.Height()
	area := b.Area()
	if nmines > area {
		//panic("more mines than squares")
		common.Fatal("more mines than squares")
	}
	nums := make([]int, area)
	for i := 0; i < area; i++ {
		b.board[i/width][i%width] = Empty
		nums[i] = i
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(area, func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
	for _, n := range nums[:nmines] {
		h, w := n/width, n%width
		b.board[h][w] = Mine
		// Upper
		if h != 0 {
			// Upper-left
			if w != 0 && b.board[h-1][w-1].IsNum() {
				b.board[h-1][w-1] += 1
			}
			// Upper-Mid
			if b.board[h-1][w].IsNum() {
				b.board[h-1][w] += 1
			}
			// Upper-right
			if w != width-1 && b.board[h-1][w+1].IsNum() {
				b.board[h-1][w+1] += 1
			}
		}

		// Mid-Left
		if w != 0 && b.board[h][w-1].IsNum() {
			b.board[h][w-1] += 1
		}
		// Mid-right
		if w != width-1 && b.board[h][w+1].IsNum() {
			b.board[h][w+1] += 1
		}

		// Lower
		if h != height-1 {
			// Lower-left
			if w != 0 && b.board[h+1][w-1].IsNum() {
				b.board[h+1][w-1] += 1
			}
			// Lower-Mid
			if b.board[h+1][w].IsNum() {
				b.board[h+1][w] += 1
			}
			// Lower-right
			if w != width-1 && b.board[h+1][w+1].IsNum() {
				b.board[h+1][w+1] += 1
			}
		}
	}
	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			b.board[h][w] = b.board[h][w].WithHidden()
		}
	}
	b.mines, b.minesFlagged, b.uncovered = nmines, 0, 0
}

// Returns the positions that were changed (now showing). If the length is
// zero, the a mine was uncovered
func (b *Board) TryUncover(pos Point) (res []Point) {
	sq := b.Get(pos)
	if sq.IsMine() {
		b.ShowAll()
	} else if !sq.IsShowing() {
		b.Set(pos, sq.WithShowing())
		res = append(res, pos)
		if sq.IsEmpty() {
			for y := -1; y <= 1; y++ {
				for x := -1; x <= 1; x++ {
					newPos := NewPoint(pos.X+x, pos.Y+y)
					if pos.Eq(newPos) || !newPos.Within(NewPoint(0, 0), b.BottomRight()) {
						continue
					}
					sq := b.Get(newPos)
					if !sq.IsShowing() {
						if sq.IsEmpty() {
							res = append(res, b.TryUncover(newPos)...)
						} else if sq.IsNum() {
							b.Set(newPos, sq.WithShowing())
							res = append(res, newPos)
						}
					}
				}
			}
		}
	} else if sq.IsNum() {
		// Chording (all bombs around square marked, try uncover safe squares)
		// TODO: make sure is correct (e.g., make sure all necessary squares are
		// uncovered)
		// TODO: Change which ones are shown as failed (marked red)?
		for y := -1; y <= 1; y++ {
			for x := -1; x <= 1; x++ {
				newPos := NewPoint(pos.X+x, pos.Y+y)
				if pos.Eq(newPos) || !newPos.Within(NewPoint(0, 0), b.BottomRight()) {
					continue
				}
				sq := b.Get(newPos)
				if !sq.IsShowing() && !sq.IsFlagged() {
					got := b.TryUncover(newPos)
					if len(got) == 0 {
						return nil
					}
					res = append(res, pos)
					res = append(res, got...)
				}
			}
		}
	}
	return
}

func (b *Board) ShowAll() {
	for h := 0; h < b.Height(); h++ {
		for w := 0; w < b.Width(); w++ {
			b.board[h][w] = b.board[h][w].WithShowing()
		}
	}
}

func (b *Board) String() string {
	s := ""
	width, height := b.Width(), b.Height()
	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			sq := b.board[h][w]
			if sq.IsNum() {
				sq += '0'
			}
			s += string(sq)
		}
		s += "\n\r"
	}
	return s
}

func (b *Board) Width() int {
	return len(b.board[0])
}

func (b *Board) Height() int {
	return len(b.board)
}

func (b *Board) Area() int {
	return b.Width() * b.Height()
}

func (b *Board) BottomRight() Point {
	return NewPoint(b.Width(), b.Height())
}

type Square byte

const (
	// All three indicate the number/mine isn't showing
	FlagHidden   Square = 0b1100_0000
	FlagFlagged  Square = 0b1000_0000
	FlagPossible Square = 0b0100_0000

	// Holds the count
	Empty Square = 0
	Mine  Square = '*'

	Possible byte = '?'
	Flagged  byte = 'F'
)

func (s Square) ToDisplayByte() byte {
	if s.IsNum() {
		switch n := byte(s.WithShowing()); n {
		case 0:
			return ' '
		default:
			return n + '0'
		}
	} else if s.IsFlagged() {
		return 'F'
	} else if s.IsPossible() {
		return '?'
	} else if s.IsMine() {
		return '*'
	}
	return ' '
}

func (s Square) ToDisplayRune() rune {
	if s.IsHidden() {
		return 0x25FC
	} else if s.IsFlagged() {
		return 'F'
	} else if s.IsPossible() {
		return '?'
	} else if s.IsMine() {
		return '*'
	}
	switch n := byte(s.WithShowing()); n {
	case 0:
		return ' '
	default:
		return rune(n + '0')
	}
}

func (s Square) ToDisplayString() string {
	if s.IsHidden() {
		return "\u25FC"
	} else if s.IsFlagged() {
		return ansi.SetBlinking + ansi.ForeRed + "F" + ansi.ForeDefault + ansi.ResetBlinking
	} else if s.IsPossible() {
		return "?"
	} else if s.IsMine() {
		return "*"
	}
	switch n := byte(s.WithShowing()); n {
	case 0:
		return " "
	case 1:
		return ansi.SetDim + ansi.ForeCyan + "1" + ansi.ForeDefault + ansi.ResetDim
	case 2:
		return ansi.ForeGreen + "2" + ansi.ForeDefault
	case 3:
		return ansi.ForeRed + "3" + ansi.ForeDefault
	case 4:
		return ansi.ForeBlue + "4" + ansi.ForeDefault
	case 5:
		return ansi.ForeMagenta + "5" + ansi.ForeDefault
	case 6:
		return ansi.ForeCyan + "6" + ansi.ForeDefault
	case 7:
		return ansi.ForeYellow + "7" + ansi.ForeDefault
	case 8:
		return string(n + '0')
	default:
		return string(n + '0')
	}
}

func (s Square) WithFlagged() Square {
	return s.WithShowing() | FlagFlagged
}

func (s Square) WithPossible() Square {
	return s.WithShowing() | FlagPossible
}

func (s Square) WithHidden() Square {
	return s | FlagHidden
}

func (s Square) WithShowing() Square {
	return s &^ FlagHidden
}

func (s Square) IsFlagged() bool {
	return s&FlagHidden == FlagFlagged
}

func (s Square) IsPossible() bool {
	return s&FlagHidden == FlagPossible
}

func (s Square) IsHidden() bool {
	return s&FlagHidden == FlagHidden
}

// If it is Flagged, Possible, or Hidden, then it is not showing.
func (s Square) IsShowing() bool {
	return s&FlagHidden == 0
}

func (s Square) IsNum() bool {
	s = s.WithShowing()
	return s >= 0 && s <= 8
}

func (s Square) IsEmpty() bool {
	return s.WithShowing() == Empty
}

func (s Square) IsMine() bool {
	return s.WithShowing() == Mine
}

type Point struct {
	X, Y int
}

func NewPoint(x, y int) Point {
	return Point{X: x, Y: y}
}

func (p Point) Eq(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

func (p Point) Within(topLeft, bottomRight Point) bool {
	return p.X >= topLeft.X && p.Y >= topLeft.Y &&
		p.X < bottomRight.X && p.Y < bottomRight.Y
}

func (p Point) WithinInclusive(topLeft, bottomRight Point) bool {
	return p.X >= topLeft.X && p.Y >= topLeft.Y &&
		p.X <= bottomRight.X && p.Y <= bottomRight.Y
}
