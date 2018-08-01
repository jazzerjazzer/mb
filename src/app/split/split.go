package split

type Datacoding string

const (
	datacodingPlain       = "plain"
	datacodingUnicode     = "unicode"
	maxSinglePlainChars   = 160
	maxSingleUnicodeChars = 70
	maxMultiPlainChars    = 153
	maxMultiUnicodeChars  = 67
)

func Split(m string) []string {
	datacoding := getDatacoding(m)
	length := getCharLength(m, datacoding)
	if datacoding == datacodingPlain {
		if length > maxSinglePlainChars {
			return splitImpl(m, maxMultiPlainChars)
		}
		return []string{m}
	}
	if length > maxSingleUnicodeChars {
		return splitImpl(m, maxMultiUnicodeChars)
	}
	return []string{m}
}

func splitImpl(message string, limit int) []string {
	totalLength := 0
	currentMessage := ""
	var splitted []string

	for _, r := range message {
		totalLength += getLength(r)
		if totalLength == limit {
			currentMessage = currentMessage + string(r)
			splitted = append(splitted, currentMessage)
			totalLength = 0
			currentMessage = ""
			continue
		} else if totalLength > limit {
			splitted = append(splitted, currentMessage)
			totalLength = 0
			currentMessage = ""
		}
		currentMessage = currentMessage + string(r)
	}
	if currentMessage != "" {
		splitted = append(splitted, currentMessage)
	}
	return splitted
	// TODO: Return the messages with UDH
}

func getCharLength(body string, datacoding Datacoding) int {
	var length int
	switch datacoding {
	case datacodingPlain:
		for _, r := range body {
			if isBaseGSM7(r) {
				length++
			} else {
				length = length + 2
			}
		}
	case datacodingUnicode:
		length = len([]rune(body))
	}
	return length
}

func getDatacoding(body string) Datacoding {
	for _, r := range body {
		if !isBaseGSM7(r) {
			if !isExtendedGSM7(r) {
				return datacodingUnicode
			}
		}
	}
	return datacodingPlain
}
