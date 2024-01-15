package utils

import "github.com/hajimehoshi/ebiten/v2"

// Returns false if the window size is too big.
func SetRawWindowSize(width, height int, logicalMargin int) bool {
	scale := ebiten.DeviceScaleFactor()
	fsWidth, fsHeight := ebiten.ScreenSizeInFullscreen()
	if width + logicalMargin > fsWidth || height + logicalMargin > fsHeight {
		return false
	}
	scaledWidth  := int(float64(width )/scale)
	scaledHeight := int(float64(height)/scale)
	ebiten.SetWindowSize(scaledWidth, scaledHeight)
	return true
}

func SetMaxMultRawWindowSize(width, height int, logicalMargin int) {
	scale := ebiten.DeviceScaleFactor()
	fsWidth, fsHeight := ebiten.ScreenSizeInFullscreen()
	maxWidthMult  := (fsWidth  - logicalMargin)/width
	maxHeightMult := (fsHeight - logicalMargin)/height
	if maxWidthMult < maxHeightMult { maxHeightMult = maxWidthMult }
	if maxHeightMult < maxWidthMult { maxWidthMult = maxHeightMult }
	if maxWidthMult <= 0 || maxHeightMult <= 0 {
		maxWidthMult  = 1
		maxHeightMult = 1
	}

	width, height = width*maxWidthMult, height*maxHeightMult
	scaledWidth  := int(float64(width )/scale)
	scaledHeight := int(float64(height)/scale)
	ebiten.SetWindowSize(scaledWidth, scaledHeight)
}
