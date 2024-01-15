package scene

import "github.com/hajimehoshi/ebiten/v2"

type Scene[Context any] interface {
	Update(context Context) (*Change, error)
	DrawLogical(canvas *ebiten.Image, foremost bool, context Context)
	DrawHiRes(canvas *ebiten.Image, foremost bool, context Context)
}
