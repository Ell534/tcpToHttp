package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

const (
	Initialised = 1
	Done        = 2
	BufferSize  = 8
)

type Request struct {
	RequestLine RequestLine
	State       int
}

func (r *Request) parse(data []byte) (int, error) {
	if r.State == Initialised {
		requestLine, bytesCon, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if bytesCon == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.State = Done
		return bytesCon, nil
	}
	if r.State == Done {
		return 0, fmt.Errorf("error: trying to read data in a done state")
	}
	return 0, fmt.Errorf("error: unknown state")
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, BufferSize, BufferSize)
	readToIndex := 0

	request := &Request{
		State: Initialised,
	}

	for request.State != Done {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if err == io.EOF {
				request.State = Done
				break
			}
			return nil, fmt.Errorf("error reading stream: %w", err)
		}
		readToIndex += numBytesRead
		numBytesParsed, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, fmt.Errorf("error parsing buffer: %w", err)
		}
		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}

	return request, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	index := bytes.Index(data, []byte(crlf))
	if index == -1 {
		return nil, 0, nil
	}

	bytesConsumed := len(data[:index]) + len(crlf)

	requestLineInfo := string(data[:index])
	requestLine, err := requestLineFromString(requestLineInfo)
	if err != nil {
		return nil, 0, err
	}

	return requestLine, bytesConsumed, nil
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
