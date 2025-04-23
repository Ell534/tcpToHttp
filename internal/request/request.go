package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("error reading request: %s\n", err.Error())
	}
	requestStr := string(request)
	requestLineStruct, err := parseRequestLine(requestStr)
	if err != nil {
		return nil, fmt.Errorf("error when parsing request line: %s\n", err.Error())
	}
	requestStruct := Request{
		*requestLineStruct,
	}

	return &requestStruct, nil
}

func parseRequestLine(request string) (*RequestLine, error) {
	splitRequest := strings.Split(request, "\r\n")
	requestLine := splitRequest[0]
	splitReqLine := strings.Split(requestLine, " ")

	if len(splitReqLine) < 3 {
		return nil, fmt.Errorf("request line has invalid number of parts, it has %d parts\n", len(splitReqLine))
	}

	method := splitReqLine[0]
	if !isUpperCase(method) {
		return nil, fmt.Errorf("method is not uppercase")
	}

	requestTarget := splitReqLine[1]

	versionParts := strings.Split(splitReqLine[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s\n", request)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognised HTTP-version: %s", httpPart)
	}

	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognised HTTP-version: %s", version)
	}

	result := RequestLine{
		HttpVersion:   versionParts[1],
		RequestTarget: requestTarget,
		Method:        method,
	}
	return &result, nil
}

func isUpperCase(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsUpper(r) || !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
