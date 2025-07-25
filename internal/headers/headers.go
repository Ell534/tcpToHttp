package headers

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	index := bytes.Index(data, []byte(crlf))
	if index == -1 {
		return 0, false, nil
	}
	if index == 0 {
		return len(crlf), true, nil
	}
	headerLine := data[:index]
	firstColonIndex := bytes.IndexByte(headerLine, ':')
	if firstColonIndex == -1 {
		return 0, false, errors.New("invalid header line, no colon, :, separator found")
	}
	if firstColonIndex == 0 {
		return 0, false, errors.New("invalid header line, no field-name present")
	}
	if headerLine[firstColonIndex-1] == ' ' {
		return 0, false, errors.New("invalid header line format, whitespace found between field-name and colon")
	}
	fieldName := string(bytes.TrimSpace(headerLine[:firstColonIndex]))
	fieldValue := string(bytes.TrimSpace(headerLine[firstColonIndex+1:]))
	if !validChar([]byte(fieldName)) {
		return 0, false, fmt.Errorf("invalid header token found: %s", fieldName)
	}
	loweredFieldName := strings.ToLower(fieldName)
	h[loweredFieldName] = fieldValue
	bytesConsumed := index + len(crlf)
	return bytesConsumed, false, nil
}

var tokenChars = []byte{'!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

func validChar(data []byte) bool {
	for _, c := range data {
		if !isChar(c) {
			return false
		}
	}
	return true
}

func isChar(c byte) bool {
	if c >= 'A' && c <= 'Z' ||
		c >= 'a' && c <= 'z' ||
		c >= '0' && c <= '9' {
		return true
	}

	return slices.Contains(tokenChars, c)
}
