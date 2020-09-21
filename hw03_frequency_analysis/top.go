package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"regexp"
	"sort"
	"strings"
)

func Top10(s string) []string {
	words := splitToWords(s)
	frequency := calcFrequency(words)
	uniqWords := getUniqWords(frequency)
	sortUniqWordsByFrequency(uniqWords, frequency)

	return take(uniqWords, 10)
}

var reSlice = regexp.MustCompile("[\\s!?.,:;`'\"(){}\\[\\]]+")

func splitToWords(s string) []string {
	return reSlice.Split(s, -1)
}

func calcFrequency(words []string) map[string]int {
	result := make(map[string]int)
	for _, word := range words {
		if word != "" && word != "-" {
			result[strings.ToLower(word)]++
		}
	}
	return result
}

func getUniqWords(frequency map[string]int) []string {
	result := make([]string, 0, len(frequency))
	for word := range frequency {
		result = append(result, word)
	}
	return result
}

func sortUniqWordsByFrequency(words []string, frequency map[string]int) {
	sort.Slice(words, func(i, j int) bool {
		return frequency[words[i]] > frequency[words[j]]
	})
}

func take(words []string, max int) []string {
	count := max
	if len(words) < max {
		count = len(words)
	}

	result := make([]string, count)
	copy(result, words)

	return result
}
