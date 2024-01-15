package input

import "errors"
import "unsafe"

type GamepadGUID struct {
	hi uint64
	lo uint64
}

func StringToGamepadGUID(guidStr string) (GamepadGUID, error) {
	var guid GamepadGUID
	if len(guidStr) != 32 {
		return GamepadGUID{}, errors.New("gamepad GUID string must have 32 characters")
	}
	
	for i := 0; i < 16; i += 2 {
		guid.hi <<= 4
		guid.hi |= uint64(hex2byte(guidStr[i + 0]))
		guid.hi <<= 4
		guid.hi |= uint64(hex2byte(guidStr[i + 1]))
	}
	for i := 16; i < 32; i += 2 {
		guid.lo <<= 4
		guid.lo |= uint64(hex2byte(guidStr[i + 0]))
		guid.lo <<= 4
		guid.lo |= uint64(hex2byte(guidStr[i + 1]))
	}

	return guid, nil
}

func BytesToGamepadGUID(buffer []byte) (GamepadGUID, error) {
	var guid GamepadGUID
	if len(buffer) != 16 {
		return guid, errors.New("guid buffer must have exactly 16 bytes")
	}
	guid.hi = uint64Load(buffer[0 : 8])
	guid.lo = uint64Load(buffer[8 : 16])
	return guid, nil
}

func (self GamepadGUID) Equal(other GamepadGUID) bool {
	return self.hi == other.hi && self.lo == other.lo
}

func (self GamepadGUID) ToBytes() []byte {
	bytes := make([]byte, 16)
	uint64Write(self.hi, bytes[0 : 8])
	uint64Write(self.lo, bytes[8 : 16])
	return bytes
}

func (self GamepadGUID) WriteToBuffer(buffer []byte) {
	if len(buffer) < 16 { panic("guid buffer must have at least 16 bytes") }
	uint64Write(self.hi, buffer[0 : 8])
	uint64Write(self.lo, buffer[8 : 16])
}

// Shorter encoding, uses 0-9, a-z, and then A-Z for
// compressing sequences of zeroes.
func (self GamepadGUID) ToZeroName() string {
	if self.hi == 0 && self.lo == 0 {
		panic("can't write name for 0 GamepadGUID")
	}
	buffer := make([]byte, 32)
	index1, zeroAcc := writeZeroCompressedUint64(self.hi, 0, buffer)
	index2,       _ := writeZeroCompressedUint64(self.lo, zeroAcc, buffer[index1 : ])
	return unsafe.String(&buffer[0], index1 + index2)
}

// Uses lowercase for hex values. For uppercase, see ToStringUpper().
func (self GamepadGUID) ToString() string {
	buffer := make([]byte, 32)
	uint64WriteLower(self.hi, buffer[ 0 : 16])
	uint64WriteLower(self.lo, buffer[16 : 32])
	return unsafe.String(&buffer[0], 32)
}

func (self GamepadGUID) ToStringUpper() string {
	buffer := make([]byte, 32)
	uint64WriteUpper(self.hi, buffer[ 0 : 16])
	uint64WriteUpper(self.lo, buffer[16 : 32])
	return unsafe.String(&buffer[0], 32)
}

// Unsafe internal function. Pre: len(buffer) = 16
func uint64WriteLower(x uint64, buffer []byte) {
	for i := 16; i > 0; i -= 2 {
		b := byte(x)
		buffer[i - 1] = toHexAsciiLower(b & 0b1111)
		buffer[i - 2] = toHexAsciiLower(b >> 4)
		x >>= 8
	}
}

// Unsafe internal function. Pre: len(buffer) = 16
func uint64WriteUpper(x uint64, buffer []byte) {
	for i := 16; i > 0; i -= 2 {
		b := byte(x)
		buffer[i - 1] = toHexAsciiUpper(b & 0b1111)
		buffer[i - 2] = toHexAsciiUpper(b >> 4)
		x >>= 8
	}
}

// Unsafe internal function. Pre: len(buffer) = 8
func uint64Write(x uint64, buffer []byte) {
	buffer[0] = byte(x >> 56)
	buffer[1] = byte(x >> 48)
	buffer[2] = byte(x >> 40)
	buffer[3] = byte(x >> 32)
	buffer[4] = byte(x >> 24)
	buffer[5] = byte(x >> 16)
	buffer[6] = byte(x >> 8)
	buffer[7] = byte(x)
}

// Unsafe internal function. Pre: len(buffer) = 8
func uint64Load(buffer []byte) uint64 {
	return (uint64(buffer[0]) << 56) | (uint64(buffer[1]) << 48) | 
		(uint64(buffer[2]) << 40) | (uint64(buffer[3]) << 32) | (uint64(buffer[4]) << 24) |
		(uint64(buffer[5]) << 16) | (uint64(buffer[6]) <<  8) | (uint64(buffer[7]) <<  0)
}

func toHexAsciiLower(value byte) byte {
	if value > 9 { return 'a' - 10 + value }
	return '0' + value
}

func toHexAsciiUpper(value byte) byte {
	if value > 9 { return 'A' - 10 + value }
	return '0' + value
}

func hex2byte(hex byte) byte {
	if hex >= '0' && hex <= '9' { return hex - '0' }
	if hex >= 'a' && hex <= 'z' { return hex - 'a' + 10 }
	if hex >= 'A' && hex <= 'Z' { return hex - 'A' + 10 }
	panic("invalid hex value '" + string(hex) + "'")
}

func writeZeroCompressedUint64(x uint64, zeroAcc uint8, buffer []byte) (int, uint8) {
	if zeroAcc > 27 { panic("zeroAcc > 27") }
	if zeroAcc < 0 { panic("zeroAcc < 0") }
	
	bufferIndex := 0
	for bitIndex := 0; bitIndex < 64; bitIndex += 4 {
		hex := byte(x >> (60 - bitIndex)) & 0b1111
		if hex == 0 { // zero accumulation case
			zeroAcc += 1
			if zeroAcc == 27 { // too many zeros, flush
				buffer[bufferIndex] = 'Z'
				zeroAcc = 0
				bufferIndex += 1
			}
		} else { // non-zero
			// flush zeros
			if zeroAcc > 0 {
				if zeroAcc == 1 {
					buffer[bufferIndex] = '0'
				} else {
					buffer[bufferIndex] = ('A' - 2) + zeroAcc
				}
				zeroAcc = 0
				bufferIndex += 1
			}

			// write actual hex
			buffer[bufferIndex] = toHexAsciiLower(hex)
			bufferIndex += 1
		}
	}

	// no need to flush remaining zeroes, they are implicit
	return bufferIndex, zeroAcc
}
