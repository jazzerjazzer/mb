package split

import (
	"app/model"
	"fmt"
)

type Datacoding string

const (
	datacodingPlain       = "plain"
	datacodingUnicode     = "unicode"
	maxSinglePlainChars   = 160
	maxSingleUnicodeChars = 70
	maxMultiPlainChars    = 153
	maxMultiUnicodeChars  = 67
)

func Split(m string) []model.Split {
	datacoding := getDatacoding(m)
	length := getCharLength(m, datacoding)
	if datacoding == datacodingPlain {
		if length > maxSinglePlainChars {
			return getUDH(splitImpl(m, maxMultiPlainChars), datacodingPlain)
		}
		return []model.Split{{Message: m, UDH: "", Datacoding: datacodingPlain}}
	}
	if length > maxSingleUnicodeChars {
		return getUDH(splitImpl(m, maxMultiPlainChars), datacodingUnicode)
	}
	return []model.Split{{Message: m, UDH: "", Datacoding: datacodingUnicode}}
}

func splitImpl(message string, limit int) []model.Split {
	totalLength := 0
	currentMessage := ""
	var splitted []model.Split

	for _, r := range message {
		totalLength += getLength(r)
		if totalLength == limit {
			currentMessage = currentMessage + string(r)
			splitted = append(splitted, model.Split{Message: currentMessage})
			totalLength = 0
			currentMessage = ""
			continue
		} else if totalLength > limit {
			splitted = append(splitted, model.Split{Message: currentMessage})
			totalLength = 0
			currentMessage = ""
		}
		currentMessage = currentMessage + string(r)
	}
	if currentMessage != "" {
		splitted = append(splitted, model.Split{Message: currentMessage})
	}

	return splitted
	// TODO: Return the messages with UDH
}

func getUDH(messages []model.Split, datacoding Datacoding) []model.Split {
	length := len(messages)
	reference := "CC"
	for i, _ := range messages {
		udh := "050003" + reference + fmt.Sprintf("%02X", length) + fmt.Sprintf("%02X", i+1)
		messages[i].UDH = udh
		messages[i].Datacoding = string(datacoding)
	}

	return messages
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