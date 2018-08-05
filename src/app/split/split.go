package split

import (
	"app/char"
	"app/model"
	"fmt"
	"math/rand"
	"time"
)

const (
	maxSinglePlainChars   = 160
	maxSingleUnicodeChars = 70
	maxMultiPlainChars    = 153
	maxMultiUnicodeChars  = 67
)

func Split(m string) []model.Split {
	datacoding := getDatacoding(m)
	length := getLength(m, datacoding)
	if datacoding == model.DatacodingPlain {
		if length > maxSinglePlainChars {
			splitted := splitImpl(m, maxMultiPlainChars)
			return attachMetadata(splitted, model.DatacodingPlain)
		}
		return []model.Split{{Message: m, UDH: "", Datacoding: model.DatacodingPlain}}
	}
	if length > maxSingleUnicodeChars {
		splitted := splitImpl(m, maxMultiUnicodeChars)
		return attachMetadata(splitted, model.DatacodingUnicode)
	}
	return []model.Split{{Message: m, UDH: "", Datacoding: model.DatacodingUnicode}}
}

func splitImpl(message string, limit int) []model.Split {
	totalLength := 0
	currentMessage := ""
	var splitted []model.Split

	for _, r := range message {
		totalLength += char.GetLength(r)
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
}

func attachMetadata(messages []model.Split, datacoding model.Datacoding) []model.Split {
	length := len(messages)
	reference := fmt.Sprintf("%02X", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(256))
	for i := range messages {
		udh := "050003" + reference + fmt.Sprintf("%02X", length) + fmt.Sprintf("%02X", i+1)
		messages[i].UDH = udh
		messages[i].Datacoding = datacoding
	}
	return messages
}

func getLength(body string, datacoding model.Datacoding) int {
	var length int
	switch datacoding {
	case model.DatacodingPlain:
		for _, r := range body {
			if char.IsBaseGSM7(r) {
				length++
			} else {
				length = length + 2
			}
		}
	case model.DatacodingUnicode:
		length = len([]rune(body))
	}
	return length
}

func getDatacoding(body string) model.Datacoding {
	for _, r := range body {
		if !char.IsBaseGSM7(r) {
			if !char.IsExtendedGSM7(r) {
				return model.DatacodingUnicode
			}
		}
	}
	return model.DatacodingPlain
}
