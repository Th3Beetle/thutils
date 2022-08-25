package thutils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

const headerDelim = "\r\n\r\n"

func ReadAll(conn *net.TCPConn) ([]byte, error) {
	errorMessage := "failed to read message: %s"
	reader := bufio.NewReader(conn)
	headers, err := readHeader(reader)

	if err != nil {
		return nil, fmt.Errorf(errorMessage, err)
	}

	headerString := string(headers)
	cl, err := extractContentLength(headerString)
	if err != nil {
		return nil, fmt.Errorf(errorMessage, err)
	}

	body := make([]byte, cl)
	io.ReadFull(reader, body)
	message := append(headers, body[:]...)
	return message, nil
}

func readHeader(reader *bufio.Reader) ([]byte, error) {
	var message []byte
	for {
		singleByte, err := reader.ReadByte()

		if err != nil {
			return nil, fmt.Errorf("failed to read header: %s", err)
		}

		message = append(message, singleByte)
		if bytes.HasSuffix(message, []byte(headerDelim)) {
			break
		}
	}
	return message, nil
}

func extractContentLength(headers string) (int, error) {
	var clValue int
	var err error

	contentLength := strings.Split(headers, "Content-Length: ")
	if len(contentLength) > 1 {
		valueString := strings.Split(contentLength[1], "\r\n")[0]
		clValue, err = strconv.Atoi(valueString)
		if err != nil {
			return 0, fmt.Errorf("failed to extract content length: %s", err)
		}
	}
	return clValue, nil
}
