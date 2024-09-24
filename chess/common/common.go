package common

import (
	"fmt"

	"github.com/johnietre/term-games/common/ansi"
)

type Piece uint8

const (
	Empty Piece = 0b0000

	Pawn   Piece = 0b0001
	Knight Piece = 0b0010
	Bishop Piece = 0b0011
	Rook   Piece = 0b0100
	Queen  Piece = 0b0101
	King   Piece = 0b0110

	WhitePawn   Piece = 0b00000 | Pawn
	WhiteKnight Piece = 0b00000 | Knight
	WhiteBishop Piece = 0b00000 | Bishop
	WhiteRook   Piece = 0b00000 | Rook
	WhiteQueen  Piece = 0b00000 | Queen
	WhiteKing   Piece = 0b00000 | King

	BlackPawn   Piece = 0b10000 | Pawn
	BlackKnight Piece = 0b10000 | Knight
	BlackBishop Piece = 0b10000 | Bishop
	BlackRook   Piece = 0b10000 | Rook
	BlackQueen  Piece = 0b10000 | Queen
	BlackKing   Piece = 0b10000 | King
)

func (p Piece) Type() Piece {
	return p & 0b0111
}

func (p Piece) String() string {
	switch p.Type() {
	case Pawn:
		return "P"
	case Knight:
		return "N"
	case Bishop:
		return "B"
	case Rook:
		return "R"
	case Queen:
		return "Q"
	case King:
		return "K"
	default:
		return " "
	}
}

func (p Piece) TermString() string {
	if p.IsBlack() {
		return ansi.ForeBlack + p.String()
	}
	return ansi.ForeWhite + p.String()
}

func (p Piece) IsEmpty() bool {
	return p == Empty
}

func (p Piece) IsWhite() bool {
	return p&0b0000 != 0 && p != 0
}

func (p Piece) IsBlack() bool {
	return p&0b1000 != 0
}

func (p Piece) IsPawn() bool {
	return p&Pawn == Pawn
}

func (p Piece) IsKnight() bool {
	return p&Knight == Knight
}

func (p Piece) IsBishop() bool {
	return p&Bishop == Bishop
}

func (p Piece) IsRook() bool {
	return p&Rook == Rook
}

func (p Piece) IsQueen() bool {
	return p&Queen == Queen
}

func (p Piece) IsKing() bool {
	return p&King == King
}

type Rank uint8

const (
	RankA Rank = 0
	RankB Rank = 1
	RankC Rank = 2
	RankD Rank = 3
	RankE Rank = 4
	RankF Rank = 5
	RankG Rank = 6
	RankH Rank = 7
)

type File uint8

const (
	File1 File = 0
	File2 File = 1
	File3 File = 2
	File4 File = 3
	File5 File = 4
	File6 File = 5
	File7 File = 6
	File8 File = 7
)

type Move uint32

const (
	MoveCastle  Move = TODO
	MoveQCastle Move = TODO
)

type BoardBytes [32]byte

type Board struct {
	Board    BoardBytes
	prevMove Move
}

var (
	ErrEmpty       = fmt.Errorf("square empty")
	ErrInvalidMove = fmt.Errorf("invalid move")
)

func (b Board) SquareTermString(rank Rank, file File) string {
	piece := b.Get(rank, file)
	pieceStr := " "
	return pieceStr + TODO
}

func (b Board) NewMove(
	fromRank Rank, fromFile File,
	toRank Rank, toFile File,
) (Move, error) {
	piece := b.Get(fromRank, fromFile)
	if piece.IsEmpty() {
		return 0, ErrEmpty
	}
	switch piece.Type() {
	case Pawn:
	case Knight:
	case Bishop:
	case Rook:
	case Queen:
	case King:
	}
	return 0, nil
}

func (b Board) newMovePawn(
	fromRank Rank, fromFile File,
	toRank Rank, toFile File,
) (Move, error) {
	rdiff, fdiff := AbsDiff(fromRank, toRank), AbsDiff(fromFile, toFile)
	if rdiff == 0 {
		return 0, ErrInvalidMove
	}
	if rdiff == 1 {
	} else if rdiff == 2 {
	}
	if fdiff == 0 {
	} else if fdiff == 1 {
	} else {
		return 0, ErrInvalidMove
	}
}

func (b Board) newMoveKnight(
	fromRank Rank, fromFile File,
	toRank Rank, toFile File,
) (Move, error) {
}

func (b Board) Get(rank Rank, file File) Piece {
	r, f := byte(rank), byte(file)
	startBit := f*24 + r
	endBit := startBit + 2
	startIndex, endIndex := startBit/8, endBit/8
	if startIndex == endIndex {
		return Piece((b.Board[startIndex] >> (startBit % 8)) & 0x07)
	}
	bitEnd := (startIndex + 1) * 8
	by := (b.Board[startIndex] >> (startBit % 8)) &^ (1 << (bitEnd - startBit))
	by |= TODO

	return Piece(by)
}

func (b *Board) Set(rank Rank, file File, piece Piece) {
	r, f := byte(rank), byte(file)
	startBit := f*24 + r
	endBit := startBit + 2
	startIndex, endIndex := startBit/8, endBit/8
	if startIndex == endIndex {
		by := b.Board[startIndex]
		by = by & (0xFF & (byte(piece) << (startBit % 8)))
		b.Board[startIndex] = TODO
	} else {
		by := b.Board[startIndex]
		by = by & TODO
		b.Board[startIndex] = by
	}
}

func AbsDiff[T Rank | File](t1, t2 T) T {
	if t1 > t2 {
		return t1 - t2
	} else {
		return t2 - t1
	}
}
