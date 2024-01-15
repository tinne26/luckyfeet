package tile

import "image"

type Orientation uint8
func (self Orientation) RotatedRight() Orientation {
	switch self & 0b0011 {
	case 0b0000: return (self & 0b0100) | 0b0001
	case 0b0001: return (self & 0b0100) | 0b0010
	case 0b0010: return (self & 0b0100) | 0b0011
	case 0b0011: return (self & 0b0100) | 0b0000
	default:
		panic("broken code")
	}
}

func (self Orientation) RotatedLeft() Orientation {
	switch self & 0b0011 {
	case 0b0000: return (self & 0b0100) | 0b0011
	case 0b0001: return (self & 0b0100) | 0b0000
	case 0b0010: return (self & 0b0100) | 0b0001
	case 0b0011: return (self & 0b0100) | 0b0010
	default:
		panic("broken code")
	}
}

func (self Orientation) Mirrored() Orientation {
	switch self & 0b0100 {
	case 0b0000: return (self & 0b0011) | 0b0100
	case 0b0100: return (self & 0b0011) | 0b0000
	default:
		panic("broken code")
	}
}

func (self Orientation) IsMirrored() bool {
	return self & 0x04 != 0
}

func (self Orientation) RotationStr() string {
	switch self & 0x03 {
	case 0b0000: return "0 deg."
	case 0b0001: return "90 deg."
	case 0b0010: return "180 deg."
	case 0b0011: return "270 deg."
	default:
		panic("broken code")
	}
}

func (self Orientation) HasRotation() bool {
	return self & 0x03 != 0
}

func (self Orientation) RotationDegs() int {
	switch self & 0x03 {
	case 0b0000: return 0
	case 0b0001: return 90
	case 0b0010: return 180
	case 0b0011: return 270
	default:
		panic("broken code")
	}
}

func (self Orientation) ApplyToTileRect(rect image.Rectangle) image.Rectangle {
	switch self & 0x03 {
	case 0b0000: // 0 degrees
		// nothing to modify
	case 0b0001: // 90 degrees
		rect = image.Rect(20 - rect.Max.Y, rect.Min.X, 20 - rect.Min.Y, rect.Max.X)
	case 0b0010: // 180 degrees
		rect = image.Rect(20 - rect.Max.X, 20 - rect.Max.Y, 20 - rect.Min.X, 20 - rect.Min.Y)
	case 0b0011: // 270 degrees
		rect = image.Rect(rect.Min.Y, 20 - rect.Max.X, rect.Max.Y, 20 - rect.Min.X)
	default:
		panic("broken code")
	}

	// mirroring
	if !self.IsMirrored() { return rect }
	return image.Rect(20 - rect.Max.X, rect.Min.Y, 20 - rect.Min.X, rect.Max.Y)
}
