package cmd

import (
	"log"
	"regexp"
	"sort"
	"unicode/utf8"
)

func Top10(text string) []string {
	if utf8.RuneCountInString(text) == ZERO {
		return nil
	}
	spl := regexp.MustCompile(REGEX).FindAllString(text, NEGATIVE)
	size := len(spl)
	tempoCalc := make(map[string]int32, size>>DIV)
	for i := ZERO; i < size; i++ {
		tempoCalc[spl[i]]++
	}
	words := make([]Word, 0, len(tempoCalc))
	for k, v := range tempoCalc {
		words = append(words, NewWord(k, v))
	}
	rsl, err := SortWord(words)
	if err != nil {
		log.Println(err)
		return nil
	}
	return rsl
}

func SortWord(words []Word) ([]string, error) {
	if len(words) == ZERO {
		return nil, ErrorEmpty
	}
	rsl := make([]string, ZERO, MAX)
	sort.Slice(words, func(i, j int) bool {
		if words[i].cnt == words[j].cnt {
			return words[i].val < words[j].val
		}
		return words[i].cnt > words[j].cnt
	})
	for i := 0; i < len(words) && i != MAX; i++ {
		rsl = append(rsl, words[i].val)
	}
	return rsl, nil
}
