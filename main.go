package main

import (
	"log"
	"strings"

	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2input"
	ebiten3 "github.com/OpenDiablo2/OpenDiablo2/d2core/d2input/ebiten"

	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2render/sdl2"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2render/ebiten"

	"github.com/OpenDiablo2/OpenDiablo2/d2app"
	ebiten2 "github.com/OpenDiablo2/OpenDiablo2/d2core/d2audio/ebiten"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2config"
	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2term"
	"github.com/OpenDiablo2/OpenDiablo2/d2script"
)

// GitBranch is set by the CI build process to the name of the branch
//nolint:gochecknoglobals This is filled in by the build system
var GitBranch string

// GitCommit is set by the CI build process to the commit hash
//nolint:gochecknoglobals This is filled in by the build system
var GitCommit string

func main() {
	log.SetFlags(log.Lshortfile)
	log.Println("OpenDiablo2 - Open source Diablo 2 engine")

	var err error

	if err = d2config.Load(); err != nil {
		panic(err)
	}

	// Initialize our providers
	var renderer d2interface.Renderer
	var inputManager d2interface.InputManager

	switch strings.ToLower(d2config.Config.Backend) {
	case "ebiten":
		renderer, err = ebiten.CreateRenderer()
		inputManager = d2input.New(ebiten3.InputService{})
	case "sdl2":
		renderer, err = sdl2.CreateRenderer()
		inputManager = d2input.New(renderer.(d2interface.InputService))
	default:
		panic("unknown renderer specified")
	}

	if err != nil {
		panic(err)
	}

	audio, err := ebiten2.CreateAudio()
	if err != nil {
		panic(err)
	}

	term, err := d2term.New(inputManager)

	if err != nil {
		log.Fatal(err)
	}

	scriptEngine := d2script.CreateScriptEngine()

	app := d2app.Create(GitBranch, GitCommit, inputManager, term, scriptEngine, audio, renderer)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
