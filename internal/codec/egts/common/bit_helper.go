package common

func Bit(flags byte, pos uint) bool {
	return (flags>>pos)&1 == 1
}

func Bits(flags byte, pos, width uint) byte {
	return (flags >> pos) & ((1 << width) - 1)
}

func BoolToBit(val bool) byte {
	if val {
		return 1
	}

	return 0
}

func SetBit(flags *byte, pos uint, val bool) {
	if val {
		*flags |= 1 << pos
	}
}

func SetBits(flags *byte, pos, width uint, value byte) {
	mask := byte((1<<width)-1) << pos
	*flags &^= mask
	*flags |= (value << pos) & mask
}
