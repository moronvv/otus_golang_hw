package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordCnt struct {
	Word string
	Cnt  int
}

func countFrequency(words []string) map[string]int {
	freq := map[string]int{}
	for _, word := range words {
		if word != "" {
			freq[word]++
		}
	}

	return freq
}

func sortByFrequency(wordsFreq map[string]int) []string {
	freq := []wordCnt{}
	for word, cnt := range wordsFreq {
		freq = append(freq, wordCnt{word, cnt})
	}

	sort.Slice(freq, func(i, j int) bool {
		if freq[i].Cnt == freq[j].Cnt {
			return freq[i].Word < freq[j].Word
		}

		return freq[i].Cnt > freq[j].Cnt
	})

	sortedWords := []string{}
	for _, wc := range freq {
		sortedWords = append(sortedWords, wc.Word)
	}

	return sortedWords
}

func Top10(s string) []string {
	words := strings.Fields(s)
	freq := countFrequency(words)
	sorted := sortByFrequency(freq)

	limit := 10
	if len(sorted) < 10 {
		limit = len(sorted)
	}
	return sorted[:limit]
}
