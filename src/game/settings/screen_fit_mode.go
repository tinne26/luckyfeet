package settings

type ScreenFitMode uint8

const (
	ScreenFitPixelPerfect ScreenFitMode = 0
	ScreenFitProportional ScreenFitMode = 1
	ScreenFitStretch      ScreenFitMode = 2
)

func (self ScreenFitMode) NextMode() ScreenFitMode {
	next := (self + 1)
	if next <= ScreenFitStretch { return next }
	return ScreenFitPixelPerfect
}
