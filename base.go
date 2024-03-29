package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"unicode/utf8"
)

type wordsBase struct {
	items           []string
	charsFreq       map[rune]int
	itemFreqIndexes map[string]int
}

func loadBase(fileName string) (*wordsBase, error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	text := strings.ReplaceAll(strings.ToLower(string(b)), "ё", "е")
	words := strings.Split(text, "\n")
	for i := range words {
		words[i] = strings.TrimSpace(words[i])
		if wlen := utf8.RuneCountInString(words[i]); wlen != wordLen {
			return nil, fmt.Errorf("wrong word length: %q (%d)", words[i], wlen)
		}
	}

	base := &wordsBase{
		items:           words,
		charsFreq:       make(map[rune]int),
		itemFreqIndexes: make(map[string]int),
	}

	if len(base.items) == 0 {
		return nil, errors.New("empty base")
	}

	wordChars := newCharSet()
	for _, word := range base.items {
		wordChars.clear()
		for _, c := range word {
			if !wordChars.has(c) {
				base.charsFreq[c]++
				wordChars.add(c)
			}
		}
	}

	for _, word := range base.items {
		wordChars.clear()
		for _, c := range word {
			if !wordChars.has(c) {
				base.itemFreqIndexes[word] += base.charsFreq[c]
				wordChars.add(c)
			}
		}
	}

	return base, nil
}

func (b *wordsBase) count() int {
	return len(b.items)
}

func (b *wordsBase) hasWord(word string) bool {
	for _, w := range b.items {
		if w == word {
			return true
		}
	}
	return false
}
