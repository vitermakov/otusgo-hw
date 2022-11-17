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
	var ready bool
	n := 40
	m := int(float64(b.readed) / float64(b.limit) * float64(n))
	l := 1
	if m > n-1 {
		l = 0
	}
	msg := "Копирование"
	if b.readed == b.limit {
		ready = true
		msg = "Готово     "
	}
	fmt.Printf(
		"\r[%s%s%s] %d%% %s",
		strings.Repeat("=", m),
		strings.Repeat(">", l),
		strings.Repeat("-", n-m-l),
		int(float64(b.readed)/float64(b.limit)*100),
		msg,
	)
	if ready {
		fmt.Print("\n")
	}
}
