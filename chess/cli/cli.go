package cli

import (
	"github.com/johnietre/term-games/common"
	"github.com/johnietre/term-games/common/ansi"
)

func Run() {
	defer common.GlobalDeferrer.Run()

	restore, _, err := common.MakeTermRaw()
	if err != nil {
		common.Fatal("error making term raw: ", err)
	}
	common.GlobalDefer(restore)

	printBoard(true)
	println("\r")
	printBoard(false)
	println("\r")
}

/*
func printBoard(asWhite bool) {
  tb := common.NewTermBuffer(nil)

  const firstRow = ansi.BackCyan+" "+ansi.BackBlue+" "+ansi.BackCyan+" "+ansi.BackBlue+" "+
    ansi.BackCyan+" "+ansi.BackBlue+" "+ansi.BackCyan+" "+ansi.BackBlue+" "+ansi.BackDefault
  const nextRow = ansi.BackBlue+" "+ansi.BackCyan+" "+ansi.BackBlue+" "+ansi.BackCyan+" "+
    ansi.BackBlue+" "+ansi.BackCyan+" "+ansi.BackBlue+" "+ansi.BackCyan+" "+ansi.BackDefault
  var nums [8]int
  if asWhite {
    nums = [8]int{8,7,6,5,4,3,2,1}
  } else {
    nums = [8]int{1,2,3,4,5,6,7,8}
  }
  for i := 0; i < 8; i++ {
    if i % 2 == 0 {
      tb.Printf("%d%s", nums[i], firstRow+"\n\r")
    } else {
      tb.Printf("%d%s", nums[i], nextRow+"\n\r")
    }
  }
  if asWhite {
    tb.Write([]byte(" ABCDEFGH"))
  } else {
    tb.Write([]byte(" HGFEDCBA"))
  }

  tb.Flush()
}

func printStartingPieces(asWhite bool) {
  const pawns = "PPPPPPPP"
  const piecesWhite = "RNBQKBNR"
  const piecesBlack = "RNBKQBNR"

  tb := common.NewTermBuffer(nil)
  if asWhite {
    tb.WriteString(ansi.ForeBlack+pawns+"\n\r"+piecesWhite+"\n\n\n\n\n\r"+ansi.ForeWhite+pawns+piecesWhite)
  } else {
    tb.WriteString(ansi.ForeWhite+pawns+"\n\r"+piecesBlack+"\n\n\n\n\n\r"+ansi.ForeBlack+pawns+piecesBlack)
  }
  tb.WriteString(ansi.ForeDefault)

  tb.Flush()
}
*/

func printBoard(asWhite bool) {
	tb := common.NewTermBuffer(nil)

	firstRow := []string{ansi.BackCyan, ansi.BackBlue, ansi.BackCyan, ansi.BackBlue,
		ansi.BackCyan, ansi.BackBlue, ansi.BackCyan, ansi.BackBlue, ansi.BackDefault}
	nextRow := []string{ansi.BackBlue, ansi.BackCyan, ansi.BackBlue, ansi.BackCyan,
		ansi.BackBlue, ansi.BackCyan, ansi.BackBlue, ansi.BackCyan, ansi.BackDefault}

	var pieces []string
	if asWhite {
		pieces = []string{"R", "N", "B", "Q", "K", "B", "N", "R"}
	} else {
		pieces = []string{"R", "N", "B", "K", "Q", "B", "N", "R"}
	}

	var nums [8]int
	if asWhite {
		nums = [8]int{8, 7, 6, 5, 4, 3, 2, 1}
	} else {
		nums = [8]int{1, 2, 3, 4, 5, 6, 7, 8}
	}
	for f := 0; f < 8; f++ {
		n := nums[f]
		tb.Print(ansi.ForeDefault, ansi.BackDefault, n)
		for r := 0; r < 8; r++ {
			piece := " "
			if f == 1 {
				if asWhite {
					piece = ansi.ForeBlack
				} else {
					piece = ansi.ForeWhite
				}
				piece += "P"
			} else if f == 6 {
				if asWhite {
					piece = ansi.ForeWhite
				} else {
					piece = ansi.ForeBlack
				}
				piece += "P"
			} else if f == 0 {
				if asWhite {
					piece = ansi.ForeBlack
				} else {
					piece = ansi.ForeWhite
				}
				piece += pieces[r]
			} else if f == 7 {
				if asWhite {
					piece = ansi.ForeWhite
				} else {
					piece = ansi.ForeBlack
				}
				piece += pieces[r]
			}
			if f%2 == 0 {
				tb.Printf("%s", firstRow[r]+piece)
			} else {
				tb.Printf("%s", nextRow[r]+piece)
			}
		}
		tb.Write([]byte("\n\r" + ansi.ForeDefault + ansi.BackDefault))
	}
	if asWhite {
		tb.Write([]byte(" ABCDEFGH"))
	} else {
		tb.Write([]byte(" HGFEDCBA"))
	}

	tb.Flush()
}
