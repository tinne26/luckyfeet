package autoconfig

import "github.com/hajimehoshi/ebiten/v2"

import "os"

func DetectResizable() {
	for _, arg := range os.Args {
		if arg == "--resizable" {
			ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
			return
		}
	}
}

func DetectWindowed() {
	for _, arg := range os.Args {
		if arg == "--windowed" { return }
	}
	ebiten.SetFullscreen(true)
}

func PreferOpenGL() {
	// allow directX if passed as program flag
	for _, arg := range os.Args {
		if arg == "--directX" || arg == "--directx" { return }
	}

	// set openGL as the graphics backend otherwise
	err := os.Setenv("EBITENGINE_GRAPHICS_LIBRARY", "opengl")
	if err != nil { panic(err) }
}

func DetectScreenshotKey() {
	for _, arg := range os.Args {
		switch arg {
		case "--qshot": 
			err := os.Setenv("EBITENGINE_SCREENSHOT_KEY", "q")
			if err != nil { panic(err) }
			return
		}
	}
}
