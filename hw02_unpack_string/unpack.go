package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

const (
	Nil  rune = '0'
	Ctrl rune = '\\'
)

// isCtrl управляющий слэш?
func isCtrl(char rune) bool {
	return Ctrl == char
}

// isDigit арабская цифра?
func isDigit(char rune) (int, bool) {
	d := char - Nil
	if d >= 0 && d < 10 {
		return int(d), true
	}
	return -1, false
}

func Unpack(input string) (string, error) {
	var prev, char rune
	var ctrl, bCtrl bool
	var bldr strings.Builder
	// перебор range по рунам - самый надежный способ
	for _, char = range input {
		digit, yes := isDigit(char)
		bCtrl = isCtrl(char)
		if ctrl { // был установлен слэш
			if !yes && !bCtrl {
				return "", ErrInvalidString
			}
			prev = char
			ctrl = false
			continue
		}
		if bCtrl { // устанавливаем слэш, сбрасываем текущий символ
			ctrl = true
			char = 0
		} else if yes { // записываем в буфер умноженный prev, сбрасываем текущий символ
			if prev <= 0 {
				return "", ErrInvalidString
			}
			bldr.WriteString(strings.Repeat(string(prev), digit))
			prev = 0
			char = 0
		}
		if prev > 0 {
			bldr.WriteRune(prev)
		}
		prev = char
	}
	// если в самом конце был слэш, то cчитаем это ошибкой (это не оговорено в ДЗ)
	if ctrl {
		return "", ErrInvalidString
	}
	if prev > 0 {
		bldr.WriteRune(prev)
	}
	return bldr.String(), nil
}
