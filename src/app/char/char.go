package char

func IsBaseGSM7(r rune) bool {
	_, ok := baseGSM7[r]
	return ok
}

func IsExtendedGSM7(r rune) bool {
	_, ok := extendedGSM7[r]
	return ok
}

func IsUnicode(r rune) bool {
	return !IsBaseGSM7(r) && !IsExtendedGSM7(r)
}

func GetLength(r rune) int {
	if IsUnicode(r) {
		return 1
	} else if IsBaseGSM7(r) {
		return 1
	} else { // Extended
		return 2
	}
}
