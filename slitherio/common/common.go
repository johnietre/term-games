package common

type Game struct {
	GameArea GameArea
	MaxLen   int
}

type GameArea struct {
	Area [][]rune
}

func (ga GameArea) Dims() GameDims {
	height, width := len(ga.Area), 0
	if height != 0 {
		width = len(ga.Area[0])
	}
	return GameDims{
		TotalWidth:  width,
		TotalHeight: height,
	}
}

type GameDims struct {
	TotalWidth     int
	TotalHeight    int
	ViewableWidth  int
	ViewableHeight int
}

func (gd GameDims) TotalArea() int {
	return gd.TotalHeight * gd.TotalWidth
}

func (gd GameDims) ViewableArea() int {
	return gd.ViewableHeight * gd.ViewableWidth
}

type Player struct {
	Id     rune
	Length int
}
