package main

import "embed"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/luckyfeet/src/game"
import "github.com/tinne26/luckyfeet/src/game/utils"
import "github.com/tinne26/luckyfeet/src/game/material/au"
import "github.com/tinne26/luckyfeet/src/plumbing/autoconfig"
import "github.com/tinne26/luckyfeet/src/plumbing/ico"

// WASM compilation on Windows:
// > $env:GOOS="js"; $env:GOARCH="wasm"; go build -o luckyfeet.wasm -trimpath main.go

//go:embed assets/*
var filesys embed.FS

func main() {
	// screenshots key (or --qshot for "Q")
	// os.Setenv("EBITENGINE_SCREENSHOT_KEY", "p")

	utils.SetMaxMultRawWindowSize(640, 360, 40)
	autoconfig.Apply()
	err := ico.LoadAndSetWindowIcons(filesys)
	if err != nil { panic(err) }
	au.Initialize()

	adapter, err := game.New(filesys)
	if err != nil { panic(err) }
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	ebiten.SetTPS(120) // had to hack this at the end due to poor jump/tic-tac key detection
	ebiten.SetWindowTitle("Lucky Feet") // \U0001F407
	err = ebiten.RunGame(adapter)
	if err != nil { panic(err) }
}
