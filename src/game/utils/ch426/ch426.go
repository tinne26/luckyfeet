package ch426 // so, I said "encode 42 bits in 6 chars" instead of "encode 7 bits per char"...

import "strings"
import "errors"

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789{}~#@+-_*[]()<>\\|/:;!?=&%$àèìòùáéíóúâêîôûäëïöüÀÈÌÒÙÁÉÍÓÚÂÊÎÔÛÄËÏÖÜ"
var lookupTable [256]uint8
var alphabetRunes []rune
func init() {
	alphabetRunes = make([]rune, 128)
	var runeCount int
	for _, codePoint := range alphabet {
		// fmt.Printf("#%d (%q) => '%v'\n", runeCount, codePoint, codePoint)
		alphabetRunes[runeCount] = codePoint
		lookupTable[uint8(codePoint)] = uint8(runeCount)
		runeCount += 1
	}
	if runeCount != 128 { panic("broken ch426 package") }
}

func Encode(data []byte) string {
	var builder strings.Builder
	
	roomBits := 7
	charValue := uint8(0)
	for _, nextByte := range data {
		charValue |= nextByte >> (8 - roomBits)

		builder.WriteRune(alphabetRunes[charValue])
		roomBits -= 1
		charValue = (nextByte << roomBits) & 0b0111_1111
		if roomBits == 0 {
			builder.WriteRune(alphabetRunes[charValue])
			roomBits = 7
			charValue = 0
		}
	}
	if roomBits != 7 {
		builder.WriteRune(alphabetRunes[charValue])
	}

	return builder.String()
}

func Decode(data string) ([]byte, error) {
	buff := make([]byte, 0, 512)

	usedBits := 0
	nextByteValue := uint8(0)
	for _, codePoint := range data {
		if codePoint > 255 {
			return nil, errors.New("character '" + string(codePoint) + "' is not valid for a ch426 encoding")
		}
		bits7 := lookupTable[uint8(codePoint)]
		if bits7 == 0 && codePoint != 'A' {
			return nil, errors.New("character '" + string(codePoint) + "' is not valid for a ch426 encoding")
		}
		
		// push the 7 bits into the buffer
		nextByteValue |= (bits7 << 1) >> usedBits
		if usedBits >= 1 {
			buff = append(buff, nextByteValue)
			nextByteValue = bits7 << (9 - usedBits)
			usedBits -= 1 // +7 and -8
		} else {
			usedBits += 7
		}
	}

	// notice: even if some bits are marked as "used", if they don't
	// constitute a complete byte, then they were only added as padding

	return buff, nil
}
