package utils

import "math"
import "image"

import "github.com/hajimehoshi/ebiten/v2"

func ProjectNearest(logicalCanvas, canvas *ebiten.Image) *ebiten.Image {
	logicalBounds, canvasBounds := logicalCanvas.Bounds(), canvas.Bounds()
	logicalWidth, logicalHeight := logicalBounds.Dx(), logicalBounds.Dy()
	canvasWidth, canvasHeight := canvasBounds.Dx(), canvasBounds.Dy()
	
	// create options
	opts := ebiten.DrawImageOptions{}

	// trivial case: both screens have the same size
	if logicalWidth == canvasWidth && logicalHeight == canvasHeight {
		canvas.DrawImage(logicalCanvas, &opts)
		return canvas
	}

	// get aspect ratios
	logicalAspectRatio := float64(logicalWidth)/float64(logicalHeight)
	canvasAspectRatio  := float64(canvasWidth)/float64(canvasHeight)
	var scalingFactor float64
	var tx, ty int = canvasBounds.Min.X, canvasBounds.Min.Y

	// compare aspect ratios	
	if logicalAspectRatio == canvasAspectRatio {
		// simple case, aspect ratios match, only scaling is necessary
		scalingFactor = float64(canvasWidth)/float64(logicalWidth)
		opts.GeoM.Scale(scalingFactor, scalingFactor)
		opts.GeoM.Translate(float64(tx), float64(ty))
		canvas.DrawImage(logicalCanvas, &opts)
	} else {
		// aspect ratios don't match, must also apply translation
		if canvasAspectRatio < logicalAspectRatio {
			// (we have excess canvas height)
			adjustedCanvasHeight := int(float64(canvasWidth)/logicalAspectRatio)
			ty += (canvasHeight - adjustedCanvasHeight)/2
			canvasHeight = adjustedCanvasHeight
		} else { // canvasAspectRatio > logicalAspectRatio
			// (we have excess canvas width)
			adjustedCanvasWidth := int(float64(canvasHeight)*logicalAspectRatio)
			tx += (canvasWidth - adjustedCanvasWidth)/2
			canvasWidth = adjustedCanvasWidth
		}

		scalingFactor := float64(canvasWidth)/float64(logicalWidth)
		opts.Filter = ebiten.FilterLinear
		opts.GeoM.Scale(scalingFactor, scalingFactor)
		opts.GeoM.Translate(float64(tx), float64(ty))
		canvas.DrawImage(logicalCanvas, &opts)
	}

	// return the scaled, active canvas area
	rect := image.Rect(tx, ty, canvasWidth + tx, canvasHeight + ty)
	return canvas.SubImage(rect).(*ebiten.Image)
}

// Unless you are on macOS, of course.
func ProjectPixelPerfect(logicalCanvas, canvas *ebiten.Image) *ebiten.Image {
	logicalBounds, canvasBounds := logicalCanvas.Bounds(), canvas.Bounds()
	logicalWidth, logicalHeight := logicalBounds.Dx(), logicalBounds.Dy()
	canvasWidth , canvasHeight  := canvasBounds.Dx(), canvasBounds.Dy()
	
	// create options
	opts := ebiten.DrawImageOptions{}

	// trivial case: both screens have the same size
	if logicalWidth == canvasWidth && logicalHeight == canvasHeight {
		canvas.DrawImage(logicalCanvas, &opts)
		return canvas
	}

	// get zoom levels
	var tx, ty int = canvasBounds.Min.X, canvasBounds.Min.Y
	xZoom := float64(canvasWidth)/float64(logicalWidth)
	yZoom := float64(canvasHeight)/float64(logicalHeight)
	zoomLevel := math.Min(xZoom, yZoom)
	var outWidth, outHeight int
	if zoomLevel < 1.0 {
		// minification (using linear filtering)
		opts.Filter = ebiten.FilterLinear
		outWidth  = int(float64(logicalWidth)*zoomLevel)
		outHeight = int(float64(logicalHeight)*zoomLevel)
	} else {
		// integer scaling
		intZoomLevel := int(zoomLevel)
		outWidth  = logicalWidth*intZoomLevel
		outHeight = logicalHeight*intZoomLevel
		zoomLevel = float64(intZoomLevel)
	}

	// projection
	tx += (canvasWidth - outWidth) >> 1
	ty += (canvasHeight - outHeight) >> 1
	opts.GeoM.Scale(zoomLevel, zoomLevel)
	opts.GeoM.Translate(float64(tx), float64(ty))
	canvas.DrawImage(logicalCanvas, &opts)

	// return active subimage
	rect := image.Rect(tx, ty, outWidth + tx, outHeight + ty)
	return canvas.SubImage(rect).(*ebiten.Image)
}

func ProjectStretched(logicalCanvas, canvas *ebiten.Image) {
	logicalBounds, canvasBounds := logicalCanvas.Bounds(), canvas.Bounds()
	logicalWidth, logicalHeight := logicalBounds.Dx(), logicalBounds.Dy()
	canvasWidth, canvasHeight := canvasBounds.Dx(), canvasBounds.Dy()
	
	// trivial case: both screens have the same size
	opts := ebiten.DrawImageOptions{}
	if logicalWidth == canvasWidth && logicalHeight == canvasHeight {
		canvas.DrawImage(logicalCanvas, &opts)
		return
	}

	// general case: apply scalings
	horzScaling := float64(canvasWidth)/float64(logicalWidth)
	vertScaling := float64(canvasHeight)/float64(logicalHeight)
	opts.Filter = ebiten.FilterLinear
	opts.GeoM.Scale(horzScaling, vertScaling)
	canvas.DrawImage(logicalCanvas, &opts)
}
