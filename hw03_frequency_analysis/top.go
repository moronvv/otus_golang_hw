package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var re = regexp.MustCompile(`([а-яА-я-]+)`)

// Clear words from non-letter symbols.
func sanitize(word string) string {
	switch word {
	case "", "-":
		return ""
	default:
		word = strings.ToLower(word)

		matches := re.FindStringSubmatch(word)
		if len(matches) > 1 {
			return matches[1]
		}

		return ""
	}
}

func countFrequency(words []string) map[string]int {
	freq := map[string]int{}
	for _, word := range words {
		word = sanitize(word)
		if word != "" {
			freq[word]++
		}
	}

	return freq
}

func sortByFrequency(wordsFreq map[string]int) []string {
	words := []string{}
	for word := range wordsFreq {
		words = append(words, word)
	}

	sort.Slice(words, func(i, j int) bool {
		iFreq, jFreq := wordsFreq[words[i]], wordsFreq[words[j]]

		// lexigraphic sort for words with same frequency
		if iFreq == jFreq {
			return words[i] < words[j]
		}
		return iFreq > jFreq
	})

	return words
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
