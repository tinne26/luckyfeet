package tile

import "image"

import "github.com/tinne26/luckyfeet/src/game/context"
import "github.com/tinne26/luckyfeet/src/game/components/tile/tcsts"

var CollisionFuncs [tcsts.GeometryMaxSentinel]func(ctx *context.Context, tileOrient Orientation, tileRect, targetRect image.Rectangle) bool
var LandingFuncs [tcsts.GeometryMaxSentinel]func(ctx *context.Context, tileOrient Orientation, tileRect image.Rectangle, ox, fx, y int) bool
func init() {
	
	// --- collisions ---

	CollisionFuncs[tcsts.GeometryNone] = func(ctx *context.Context, tileOrient Orientation, tileRect, targetRect image.Rectangle) bool {
		return false
	}
	CollisionFuncs[tcsts.Geometry20x20] = func(ctx *context.Context, tileOrient Orientation, tileRect, targetRect image.Rectangle) bool {
		return targetRect.Overlaps(tileRect)
	}
	CollisionFuncs[tcsts.GeometryBL20x19] = func(ctx *context.Context, tileOrient Orientation, tileRect, targetRect image.Rectangle) bool {
		sr := image.Rect(0, 1, 20, 20)
		return targetRect.Overlaps(tileOrient.ApplyToTileRect(sr).Add(tileRect.Min))
	}
	CollisionFuncs[tcsts.GeometryTR19x19] = func(ctx *context.Context, tileOrient Orientation, tileRect, targetRect image.Rectangle) bool {
		sr := image.Rect(1, 0, 20, 19)
		return targetRect.Overlaps(tileOrient.ApplyToTileRect(sr).Add(tileRect.Min))
	}
	CollisionFuncs[tcsts.GeometryBL20x9] = func(ctx *context.Context, tileOrient Orientation, tileRect, targetRect image.Rectangle) bool {
		sr := image.Rect(0, 11, 20, 20)
		return targetRect.Overlaps(tileOrient.ApplyToTileRect(sr).Add(tileRect.Min))
	}
	CollisionFuncs[tcsts.GeometryBR19x9] = func(ctx *context.Context, tileOrient Orientation, tileRect, targetRect image.Rectangle) bool {
		sr := image.Rect(1, 11, 20, 20)
		return targetRect.Overlaps(tileOrient.ApplyToTileRect(sr).Add(tileRect.Min))
	}
	CollisionFuncs[tcsts.GeometryMT18x17] = func(ctx *context.Context, tileOrient Orientation, tileRect, targetRect image.Rectangle) bool {
		sr := image.Rect(1, 0, 19, 17)
		return targetRect.Overlaps(tileOrient.ApplyToTileRect(sr).Add(tileRect.Min))
	}
	CollisionFuncs[tcsts.GeometryBR19x20] = func(ctx *context.Context, tileOrient Orientation, tileRect, targetRect image.Rectangle) bool {
		sr := image.Rect(1, 0, 20, 20)
		return targetRect.Overlaps(tileOrient.ApplyToTileRect(sr).Add(tileRect.Min))
	}
	CollisionFuncs[tcsts.GeometryBR18x16] = func(ctx *context.Context, tileOrient Orientation, tileRect, targetRect image.Rectangle) bool {
		sr := image.Rect(2, 4, 20, 20)
		return targetRect.Overlaps(tileOrient.ApplyToTileRect(sr).Add(tileRect.Min))
	}
	CollisionFuncs[tcsts.GeometryBL17x16] = func(ctx *context.Context, tileOrient Orientation, tileRect, targetRect image.Rectangle) bool {
		sr := image.Rect(0, 4, 17, 20)
		return targetRect.Overlaps(tileOrient.ApplyToTileRect(sr).Add(tileRect.Min))
	}
	CollisionFuncs[tcsts.GeometryBL1_17x16] = func(ctx *context.Context, tileOrient Orientation, tileRect, targetRect image.Rectangle) bool {
		sr := image.Rect(1, 4, 18, 20)
		return targetRect.Overlaps(tileOrient.ApplyToTileRect(sr).Add(tileRect.Min))
	}
	CollisionFuncs[tcsts.GeometryMM4x4] = func(ctx *context.Context, tileOrient Orientation, tileRect, targetRect image.Rectangle) bool {
		return targetRect.Overlaps(image.Rect(8, 8, 12, 12).Add(tileRect.Min))
	}
	

	// --- landings ---

	noLandingFunc := func(ctx *context.Context, tileOrient Orientation, tileRect image.Rectangle, ox, fx, y int) bool {
		return false
	}
	for i := 0; i < tcsts.GeometryMaxSentinel; i++ {
		LandingFuncs[i] = noLandingFunc
	}
	// LandingFuncs[tcsts.Geometry20x20] = func(ctx *context.Context, tileOrient Orientation, tileRect image.Rectangle, ox, fx, y int) bool {
	// 	return y == tileRect.Min.Y && tileRect.Min.X <= fx && tileRect.Max.X >= ox
	// }
	// LandingFuncs[tcsts.GeometryBL20x19] = func(ctx *context.Context, tileOrient Orientation, tileRect image.Rectangle, ox, fx, y int) bool {
	// 	rect := tileOrient.ApplyToTileRect(image.Rect(0, 1, 20, 20)).Add(tileRect.Min)
	// 	return rect.Min.Y == y && rect.Min.X <= fx && rect.Max.X >= ox
	// }
	// LandingFuncs[tcsts.GeometryTR19x19] = func(ctx *context.Context, tileOrient Orientation, tileRect image.Rectangle, ox, fx, y int) bool {
	// 	rect := tileOrient.ApplyToTileRect(image.Rect(1, 0, 20, 19)).Add(tileRect.Min)
	// 	return rect.Min.Y == y && rect.Min.X <= fx && rect.Max.X >= ox
	// }
	LandingFuncs[tcsts.GeometryBL20x9] = func(ctx *context.Context, tileOrient Orientation, tileRect image.Rectangle, ox, fx, y int) bool {
		rect := tileOrient.ApplyToTileRect(image.Rect(0, 11, 20, 20)).Add(tileRect.Min)
		return rect.Min.Y == y && rect.Min.X <= fx && rect.Max.X >= ox
	}
	LandingFuncs[tcsts.GeometryBR19x9] = func(ctx *context.Context, tileOrient Orientation, tileRect image.Rectangle, ox, fx, y int) bool {
		rect := tileOrient.ApplyToTileRect(image.Rect(1, 11, 20, 20)).Add(tileRect.Min)
		return rect.Min.Y == y && rect.Min.X <= fx && rect.Max.X >= ox
	}
	LandingFuncs[tcsts.GeometryMT18x17] = func(ctx *context.Context, tileOrient Orientation, tileRect image.Rectangle, ox, fx, y int) bool {
		rect := tileOrient.ApplyToTileRect(image.Rect(1, 0, 19, 17)).Add(tileRect.Min)
		return rect.Min.Y == y && rect.Min.X <= fx && rect.Max.X >= ox
	}

	LandingFuncs[tcsts.GeometryBR18x16] = func(ctx *context.Context, tileOrient Orientation, tileRect image.Rectangle, ox, fx, y int) bool {
		rect := tileOrient.ApplyToTileRect(image.Rect(2, 4, 20, 20)).Add(tileRect.Min)
		return rect.Min.Y == y && rect.Min.X <= fx && rect.Max.X >= ox
	}
	LandingFuncs[tcsts.GeometryBL17x16] = func(ctx *context.Context, tileOrient Orientation, tileRect image.Rectangle, ox, fx, y int) bool {
		rect := tileOrient.ApplyToTileRect(image.Rect(0, 4, 17, 20)).Add(tileRect.Min)
		return rect.Min.Y == y && rect.Min.X <= fx && rect.Max.X >= ox
	}
	LandingFuncs[tcsts.GeometryBL1_17x16] = func(ctx *context.Context, tileOrient Orientation, tileRect image.Rectangle, ox, fx, y int) bool {
		rect := tileOrient.ApplyToTileRect(image.Rect(1, 4, 18, 20)).Add(tileRect.Min)
		return rect.Min.Y == y && rect.Min.X <= fx && rect.Max.X >= ox
	}
}
