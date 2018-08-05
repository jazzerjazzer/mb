package char

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsBaseGSM7(t *testing.T) {
	tests := []struct {
		description string
		givenRune   rune
		isBaseGSM7  bool
	}{
		{
			description: "empty base case",
		},
		{
			description: "unicode/ç",
			givenRune:   'ç',
			isBaseGSM7:  false,
		},
		{
			description: "unicode/Ş",
			givenRune:   'Ş',
			isBaseGSM7:  false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.description, func(t *testing.T) {
			require.Equal(t, testCase.isBaseGSM7, IsBaseGSM7(testCase.givenRune))
		})
	}
	for k := range baseGSM7 {
		t.Run(fmt.Sprintf("BaseGSM7Charset: %s", string(k)), func(t *testing.T) {
			require.True(t, IsBaseGSM7(k))
		})
	}
	for k := range extendedGSM7 {
		t.Run(fmt.Sprintf("ExtendedGSM7Charset: %s", string(k)), func(t *testing.T) {
			require.False(t, IsBaseGSM7(k))
		})
	}
}

func TestIsExtendedGSM7(t *testing.T) {
	tests := []struct {
		description string
		givenRune   rune
		isExtended  bool
	}{
		{
			description: "empty base case",
			isExtended:  false,
		},
		{
			description: "unicode",
			givenRune:   'ç',
			isExtended:  false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.description, func(t *testing.T) {
			require.Equal(t, testCase.isExtended, IsExtendedGSM7(testCase.givenRune))
		})
	}
	for k := range extendedGSM7 {
		t.Run(fmt.Sprintf("ExtendedGSM7Charset: %s", string(k)), func(t *testing.T) {
			require.True(t, IsExtendedGSM7(k))
		})
	}
	for k := range baseGSM7 {
		t.Run(fmt.Sprintf("BaseGSM7Charset: %s", string(k)), func(t *testing.T) {
			require.False(t, IsExtendedGSM7(k))
		})
	}
}

func TestIsUnicode(t *testing.T) {
	tests := []struct {
		description string
		givenRune   rune
		isUnicode   bool
	}{
		{
			description: "unicode/ç",
			givenRune:   'ç',
			isUnicode:   true,
		},
		{
			description: "unicode/Ş",
			givenRune:   'Ş',
			isUnicode:   true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.description, func(t *testing.T) {
			require.Equal(t, testCase.isUnicode, IsUnicode(testCase.givenRune))
		})
	}
	for k := range baseGSM7 {
		t.Run(fmt.Sprintf("BaseGSM7Charset: %s", string(k)), func(t *testing.T) {
			require.False(t, IsUnicode(k))
		})
	}
	for k := range extendedGSM7 {
		t.Run(fmt.Sprintf("ExtendedGSM7Charset: %s", string(k)), func(t *testing.T) {
			require.False(t, IsUnicode(k))
		})
	}
}

func TestGetLength(t *testing.T) {
	tests := []struct {
		description    string
		givenRune      rune
		expectedLength int
	}{
		{
			description:    "unicode/ç",
			givenRune:      'ç',
			expectedLength: 1,
		},
		{
			description:    "unicode/Ş",
			givenRune:      'Ş',
			expectedLength: 1,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.description, func(t *testing.T) {
			require.Equal(t, testCase.expectedLength, GetLength(testCase.givenRune))
		})
	}
	for k := range baseGSM7 {
		t.Run(fmt.Sprintf("BaseGSM7Charset: %s", string(k)), func(t *testing.T) {
			require.Equal(t, 1, GetLength(k))
		})
	}
	for k := range extendedGSM7 {
		t.Run(fmt.Sprintf("ExtendedGSM7Charset: %s", string(k)), func(t *testing.T) {
			require.Equal(t, 2, GetLength(k))
		})
	}
}
