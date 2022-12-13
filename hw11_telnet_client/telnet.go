package main

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TClient{address: address, timeout: timeout, in: in, out: out}
}

// TClient реализация TelnetClient.
type TClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer

	handle net.Conn
	// Исходя из предписанных тестов функции Send() и Receive() должны
	// возвращать управление после того, как из потока читается символ переноса строки.
	// Чтобы не делать эту логику самостоятельно для чтения in -> handle и handle -> out
	// посредников bufio.Scanner
	sendBuff *bufio.Scanner
	recvBuff *bufio.Scanner
}

func (tc *TClient) Connect() error {
	h, err := net.DialTimeout("tcp", tc.address, tc.timeout)
	if err != nil {
		return err
	}
	tc.handle = h

	tc.sendBuff = bufio.NewScanner(tc.in)
	// tc.sendBuff.Split(bufio.ScanLines)

	tc.recvBuff = bufio.NewScanner(tc.handle)
	tc.recvBuff.Split(bufio.ScanLines)

	return nil
}

func (tc *TClient) Close() error {
	return tc.handle.Close()
}

func (tc *TClient) copyLine(writer io.Writer, scanner *bufio.Scanner) error {
	// сканируем один раз до EOF или LN
	if !scanner.Scan() {
		return io.EOF
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	// так как данных может быть много используем функцию io.Copy.
	data := append(scanner.Bytes(), '\n')
	_, err := io.Copy(writer, bytes.NewReader(data))
	return err
}

func (tc *TClient) Send() error {
	return tc.copyLine(tc.handle, tc.sendBuff)
}

func (tc *TClient) Receive() error {
	return tc.copyLine(tc.out, tc.recvBuff)
}
