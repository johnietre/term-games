package cli

import (
	"math/rand"
	"time"
)

func Run() {
	board := NewBoard(9, 9)
	board.Randomize(10)
	println(board.String())
}

type Square byte

const (
	Empty    Square = 0
	Bomb     Square = '*'
	Possible Square = '?'
	Flag     Square = 'F'
)

func (s Square) IsNum() bool {
	return s >= 0 && s <= 8
}

func (s Square) IsBomb() bool {
	return s == Bomb
}

type Board struct {
	board [][]Square
}

func NewBoard(width, height int) Board {
	board := make([][]Square, height)
	for i := 0; i < height; i++ {
		board[i] = make([]Square, width)
	}
	return Board{
		board: board,
	}
}

func (b Board) Randomize(nbombs int) {
	width, height := b.Width(), b.Height()
	area := b.Area()
	if nbombs > area {
		panic("more bombs than squares")
	}
	nums := make([]int, area)
	for i := 0; i < area; i++ {
		b.board[i/height][i%width] = Empty
		nums[i] = i
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(area, func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
	for _, n := range nums[:nbombs] {
		h, w := n/height, n%width
		b.board[h][w] = Bomb
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
}

func (b Board) String() string {
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

func (b Board) Width() int {
	return len(b.board[0])
}

func (b Board) Height() int {
	return len(b.board)
}

func (b Board) Area() int {
	return b.Width() * b.Height()
}
