package server

import (
	scommon "github.com/johnietre/term-games/slitherio/common"
	jtutils "github.com/johnietre/utils/go"
)

func Run() {
}

type App struct {
	games *jtutils.SyncMap[string, scommon.GameArea]
}
