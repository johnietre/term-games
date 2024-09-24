package cli

import (
	"log"

	"github.com/johnietre/term-games/common"
	jtutils "github.com/johnietre/utils/go"
)

var (
	termWidth  int
	termHeight int

	globalDeferrer = common.NewSyncDeferredFunc(jtutils.NewT(true))
	termBuf        = common.NewTermBuffer(nil)
)

func Run() {
	log.SetFlags(0)
	defer globalDeferrer.Run()

	var err error
	termWidth, termHeight, err = common.GetTermSize()
	if err != nil {
		Fatal("error getting term dimensions: ", err)
	}

	restore, _, err := common.MakeTermRaw()
	if err != nil {
		Fatal("error making term raw: ", err)
	}
	globalDeferrer.Add(restore)

	termBuf.PrintfFlush(
		"Width: %d, height: %d, Area: %d\r\n",
		termWidth, termHeight, termWidth*termWidth,
	)
}

type App struct {
	termBuf *common.TermBuffer
}

func (app *App) run() {
}

func runGame() {
}

type GameArea struct {
	width  int
	height int
}

type Player struct {
}

func Fatal(args ...any) {
	globalDeferrer.Run()
	log.Fatal(args...)
}

func Fatalf(format string, args ...any) {
	globalDeferrer.Run()
	log.Fatalf(format, args...)
}

func Fatalln(args ...any) {
	globalDeferrer.Run()
	log.Fatalln(args...)
}
