package interfaces

import "github.com/hajimehoshi/ebiten/v2"

// Common interfaces shared across multiple files to escape
// excessive type coupling and dependency cycles.

type Background[Context any] interface {
	Update()
	DrawLogical(canvas *ebiten.Image, ctx Context)
}
