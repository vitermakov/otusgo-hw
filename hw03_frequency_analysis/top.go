package hw03frequencyanalysis

import "sort"

const (
	spaceDelimiter = ' '
	tabDelimiter   = '\t'
	nlDelimiter    = '\n'
	/* dashDelimiter  = '-'. */
)

// FreqWord слово с частотой.
type FreqWord struct {
	word  string
	count int
}

// Dictionary хранилище слов с частотностью.
type Dictionary struct {
	data  []FreqWord
	index map[string]int
}

// DictionaryNew конструктор - память под слайс и мапу должы быть выделены.
func DictionaryNew() *Dictionary {
	return &Dictionary{
		data:  make([]FreqWord, 0),
		index: make(map[string]int),
	}
}

// takeInto учесть одно свойство.
func (d *Dictionary) takeInto(word string) {
	// пропускаем пустые слова
	if len(word) == 0 {
		return
	}
	i, exists := d.index[word]
	if exists {
		d.data[i].count++
	} else {
		d.index[word] = len(d.data)
		d.data = append(d.data, FreqWord{
			word:  word,
			count: 1,
		})
	}
}

// Analyze анализ текста.
func (d *Dictionary) Analyze(text string, tokenizer *Tokenizer) {
	tokenizer.Proc(text, func(word string, dl rune) {
		if len(word) > 0 {
			d.takeInto(word)
		}
		// if dl == dashDelimiter {
		//	d.takeInto(string(dl))
		// }
	})
}

func (d Dictionary) GetTop(n int) []FreqWord {
	// частотности сортируем в обратном порядке, слова в прямом
	sort.Slice(d.data, func(i, j int) bool {
		if d.data[i].count == d.data[j].count {
			return d.data[i].word < d.data[j].word
		}
		return d.data[i].count > d.data[j].count
	})
	// выбираем последние n, если они есть
	if len(d.data) <= n {
		return d.data
	}
	return d.data[0:n]
}

// Tokenizer токенизатор.
type Tokenizer struct {
	dlChars []rune
}

func TokenizerNew(delims []rune) *Tokenizer {
	return &Tokenizer{
		dlChars: delims,
	}
}

// Proc при каждом новом токене вызываем callback.
func (t Tokenizer) Proc(text string, callback func(word string, dl rune)) {
	var prev int
	var word string
	for i, char := range text {
		if t.isDelimiter(char) {
			word = text[prev:i]
			callback(word, char)
			prev = i + 1
		}
	}
	word = text[prev:]
	callback(word, 0)
}

func (t Tokenizer) isDelimiter(c rune) bool {
	for _, v := range t.dlChars {
		if v == c {
			return true
		}
	}
	return false
}

func Top10(text string) []string {
	dict := DictionaryNew()
	// чтобы результаты совпадали с результатами теста необходимо
	// кроме пробелов как разделители использовать знаки табов и переносов строк
	tokenizer := TokenizerNew([]rune{
		spaceDelimiter, nlDelimiter, tabDelimiter, // dashDelimiter,
	})
	dict.Analyze(text, tokenizer)
	topFreq := dict.GetTop(10)

	result := make([]string, len(topFreq))
	for i, v := range topFreq {
		result[i] = v.word
	}
	return result
}
