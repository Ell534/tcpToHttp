package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	request, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	requestLine, err := parseRequestLine(request)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	index := bytes.Index(data, []byte(crlf))
	if index == -1 {
		return nil, fmt.Errorf("could not find a crlf in the request line")
	}

	requestLineInfo := string(data[:index])
	requestLine, err := requestLineFromString(requestLineInfo)
	if err != nil {
		return nil, err
	}

	return requestLine, nil
}

func requestLineFromString(reqStr string) (*RequestLine, error) {
	reqParts := strings.Split(reqStr, " ")
	if len(reqParts) != 3 {
		return nil, fmt.Errorf("request line is not formatted correctly, more than 3 parts: %s", reqStr)
	}

	method := reqParts[0]
	for _, char := range method {
		if char < 'A' || char > 'Z' {
			return nil, fmt.Errorf("the method is invalid: %s", method)
		}
	}

	reqTarget := reqParts[1]

	httpVersion := reqParts[2]
	versionParts := strings.Split(httpVersion, "/")
	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognised HTTP version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognised HTTP version: %s", version)
	}

	return &RequestLine{
		HttpVersion:   version,
		RequestTarget: reqTarget,
		Method:        method,
	}, nil
}
