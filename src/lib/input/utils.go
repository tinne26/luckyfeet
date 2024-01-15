package input

func uint8toBinaryStr(value uint8) string {
	var bytes [10]byte = [10]byte{'0', 'b', '0', '0', '0', '0', '0', '0', '0', '0'}
	for i := 0; i < 8; i++ {
		if value & (1 << i) != 0 {
			bytes[9 - i] = '1'
		}
	}
	return string(bytes[:])
}
