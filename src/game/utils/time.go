package utils

import "fmt"
import "math"

func FmtTicksToTimeStrCents(ticks int) string {
	secs := float64(ticks)/120.0
	secsw, secsf := math.Modf(secs)
	secsi := int(secsw)
	mins  := secsi/60.0
	secsi -= mins*60
	return fmt.Sprintf("%02d:%s", mins, paddedSecs(float64(secsi) + secsf))
}

func paddedSecs(secs float64) string {
	if secs < 10 {
		return fmt.Sprintf("0%.02f", secs)
	} else {
		return fmt.Sprintf("%.02f", secs)
	}
}

func FmtTicksToTimeStrSecs(ticks int) string {
	secs := ticks/120.0
	mins := secs/60.0
	secs -= mins*60
	return fmt.Sprintf("%02d:%02d", mins, secs)
}
