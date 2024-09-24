package cli

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/johnietre/term-games/common"
	"github.com/johnietre/term-games/common/ansi"
	inputpkg "github.com/johnietre/term-games/common/input"
	utils "github.com/johnietre/utils/go"
)

var (
	cactusIcon rune = '|'
	duckIcon   rune = '-'

	dinoIcon rune = '*'

	dinoHead rune = 0x25B5
	dinoBody rune = 0x2588
	dinoFeet rune = 0x21D3

	dinoHorHead rune = 0x25B6
	dinoHorBody rune = 0x2586
	dinoHorFeet rune = 0x21D0
)

func Run() {
	defer common.GlobalDeferrer.Run()

	width, _, err := common.GetTermSize()
	if err != nil {
		common.Fatal("error getting term dims: ", err)
	}

	restore, _, err := common.MakeTermRaw()
	if err != nil {
		common.Fatal("error making term raw: ", err)
	}
	common.GlobalDefer(restore)

	game := newGame(width)
	game.run()
}

type Game struct {
	player        Object
	ground        []Object
	viewableWidth int

	speed         time.Duration
	distance      int
	newObjProb    float32
	consecNewObjs int

	tb *common.TermBuffer

	rng *rand.Rand

	currInput atomic.Uint64
}

const (
	screenHeight = 20
)

func newGame(width int) *Game {
	return &Game{
		player: Object{
			Bottom: 0,
			Top:    3,
		},
		ground: make([]Object, width+1),
		// Correct initial speed; want to travel 100 in ~10s
		speed:         time.Millisecond * 100,
		distance:      0,
		newObjProb:    1.0,
		consecNewObjs: 0,
		viewableWidth: width,
		tb:            common.NewTermBuffer(nil),
		rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (game *Game) run() {
	go game.listenInput()
	game.init()
	game.render()
	prevInputNum := uint32(0)
	numSpaces := 0
	horizontalCount := -1
GameLoop:
	for ; true; game.distance++ {
		// When to render
		game.renderDistance()
		game.movePlayer()
		game.renderWithNewObj(game.newObj())
		if game.distance%100 == 0 {
			// Adjust
			game.speed -= 50 * game.speed / 1000
		}
		game.flush()
		next := time.Now().Add(game.speed)
		for i := 0; !time.Now().After(next); i++ {
			uin := game.currInput.Load()
			num := uint32(uin >> 32)
			if game.player.Horizontal {
				if horizontalCount >= 0 {
					if horizontalCount != 10000 {
						horizontalCount++
					} else {
						print("done")
						horizontalCount = -1
						game.player.Horizontal = false
					}
				}
			}
			if num == prevInputNum {
				continue
			}
			prevInputNum = num
			//input := game.getInput()
			input := inputpkg.Input(uin)
			if input == inputpkg.Input(' ') {
				numSpaces++
			} else {
				if numSpaces != 0 {
					fmt.Fprintln(outFile, "RESETING:", numSpaces)
				}
				numSpaces = 0
			}
			if input == inputpkg.ArrowDown {
				game.player.Horizontal = true
				horizontalCount = 0
				continue
			}
			switch input {
			case inputpkg.CtrlC:
				break GameLoop
			}
			if time.Now().After(next) {
				break
			}
		}
	}
	game.tb.Print("\n\n")
	game.tb.Flush()
}

func (game *Game) getInput() inputpkg.Input {
	buf, n, err := common.ReadStdinBytes()
	if err != nil {
		common.Fatal("error reading from stdin: ", err)
	} else if n == 0 {
		return inputpkg.Unknown
	}
	return inputpkg.FromBytes(buf[:])
}

func (game *Game) listenInput() {
	exiting := false
	for i := uint32(1); true; i++ {
		buf, n, err := common.ReadStdinBytes()
		if err != nil {
			common.Fatal("error reading from stdin: ", err)
		} else if n == 0 {
			continue
		}
		input := inputpkg.FromBytes(buf[:])
		if input == inputpkg.CtrlC {
			if exiting {
				common.Exit(0)
			}
			exiting = true
		}
		uin := (uint64(i) << 32) | uint64(input)
		game.currInput.Store(uin)
	}
}

func (game *Game) init() {
	game.tb.WriteByte('\r')
	for i := 0; i < screenHeight; i++ {
		game.tb.WriteByte('\n')
	}
	for i := 0; i < game.viewableWidth; i++ {
		game.tb.WriteByte('=')
	}
	game.tb.WriteString(ansi.CurUp)

	/*
	  game.ground[3] = Object{
	    Top: 3,
	    Bottom: 0,
	  }
	  game.ground[7] = Object{
	    Top: 1,
	    Bottom: 0,
	  }
	  game.ground[8] = Object{
	    Top: 2,
	    Bottom: 0,
	  }
	  game.ground[9] = Object{
	    Top: 3,
	    Bottom: 0,
	  }
	*/
}

func (game *Game) movePlayer() {
}

func (game *Game) renderDistance() {
	game.tb.WriteString(ansi.CurUpN(screenHeight - 1))
	game.tb.WriteByte('\r')
	game.tb.Print("Distance: ", game.distance)
	game.tb.WriteString(ansi.CurDownN(screenHeight - 1))
}

func (game *Game) render() {
	//
}

func (game *Game) renderPlayer() {
	// TODO: clear when changing between vert and hor
	game.tb.WriteString(ansi.CurRight)
	game.player.renderPlayerTo(game.tb)
	game.tb.WriteString(ansi.CurLeft)
}

func (game *Game) renderWithNewObj(newObj Object) {
	game.tb.WriteByte('\r')
	for i := 0; i < screenHeight-2; i++ {
		game.tb.WriteString(ansi.CurUp)
		game.tb.WriteString(ansi.ClearLine)
	}
	game.tb.WriteString(ansi.CurDownN(screenHeight - 2))
	game.tb.WriteString(ansi.ClearLine)

	game.renderPlayer()

	i, width := 0, len(game.ground)
	for i = 1; i < game.viewableWidth; i++ {
		obj := game.ground[i]
		obj.renderTo(game.tb)
		game.tb.WriteString(ansi.CurRight)
		game.ground[i-1] = obj
	}
	for ; i < width; i++ {
		game.ground[i-1] = game.ground[i]
	}
	game.ground[width-1] = newObj

	game.tb.WriteByte('\r')
}

func (game *Game) flush() {
	if _, err := game.tb.Flush(); err != nil {
		common.Fatal("error writing to stdout: ", err)
	}
}

func (game *Game) newObj() Object {
	if game.rng.Float32() > game.newObjProb {
		game.newObjProb *= 1 + 0.1
		// TODO
		game.consecNewObjs--
		return Object{}
	} else if game.consecNewObjs == 3 {
		// TODO: newObjProb
		// TODO: do better
		game.consecNewObjs *= -1
		return Object{}
	} else if game.consecNewObjs < 0 {
		// TODO: newObjProb
		// TODO: do better
		game.consecNewObjs++
		return Object{}
	}

	game.newObjProb *= 1 - 0.1
	game.consecNewObjs++

	top := game.rng.Intn(5)
	return Object{
		Bottom: 0,
		Top:    top,
	}
}

var outFile, _ = utils.OpenAppend("out.txt")
var _ = fmt.Print

// An object in the game (e.g., cactus or bird)
type Object struct {
	// 0 means ground
	Bottom int
	// Exclusive
	Top int
	// When horizontal, the bottom is the left and top is right
	Horizontal bool
}

func (o Object) Height() int {
	return o.Top - o.Bottom
}

func (o Object) renderTo(w writer) {
	height := o.Height()
	if height == 0 {
		return
	}
	icon := cactusIcon
	if o.Bottom != 0 {
		icon = duckIcon
		w.WriteString(ansi.CurUpN(o.Bottom))
	}
	assert(o.Top != 0)
	for h := 0; h < height; h++ {
		w.WriteRune(icon)
		w.WriteString(ansi.CurLeft)
		w.WriteString(ansi.CurUp)
	}
	w.WriteString(ansi.CurDownN(height))
}

func (o Object) renderPlayerTo(w writer) {
	height := o.Height()
	if height == 0 {
		return
	}
	if o.Horizontal {
		for h := 0; h < height; h++ {
			//w.WriteRune(dinoIcon)
			switch h {
			case 0:
				w.WriteRune(dinoHorFeet)
			case 1:
				w.WriteRune(dinoHorBody)
			case 2:
				w.WriteRune(dinoHorHead)
			}
		}
		w.WriteString(ansi.CurLeftN(height))
	} else {
		if o.Bottom != 0 {
			w.WriteString(ansi.CurUpN(o.Bottom))
		}
		for h := 0; h < height; h++ {
			//w.WriteRune(dinoIcon)
			switch h {
			case 0:
				w.WriteRune(dinoFeet)
			case 1:
				w.WriteRune(dinoBody)
			case 2:
				w.WriteRune(dinoHead)
			}
			w.WriteString(ansi.CurLeft)
			w.WriteString(ansi.CurUp)
		}
		w.WriteString(ansi.CurDownN(height))
	}
}

type writer interface {
	Write([]byte) (int, error)
	WriteString(string) (int, error)
	WriteByte(byte) error
	WriteRune(rune) (int, error)
}

func assert(b bool, msg ...string) {
	if !b {
		if len(msg) != 0 {
			panic(msg[0])
		}
		panic("assertion failed")
	}
}
