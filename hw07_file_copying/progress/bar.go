package progress

import (
	"fmt"
	"io"
	"strings"
)

type Bar struct {
	src    io.Reader
	readed int64
	limit  int64
}

func New(src io.Reader, limit int64) *Bar {
	return &Bar{src, 0, limit}
}
func (b *Bar) Read(p []byte) (int, error) {
	n, err := b.src.Read(p)
	if err != nil {
		b.displayBar(err)
		return n, err
	}
	b.readed += int64(n)
	b.displayBar(nil)

	return n, nil
}
func (b *Bar) displayBar(err error) {
	n := 20
	m := int(b.readed / b.limit * int64(n))
	l := 1
	if m > n-1 {
		l = 0
	}
	fmt.Print(
		"[%s%s%s] %d%%",
		strings.Repeat("=", m),
		strings.Repeat(">", l),
		strings.Repeat("-", n-m-l),
		int(b.readed/b.limit*100),
	)
	if err == io.EOF {
		fmt.Print("\n")
	}
}
