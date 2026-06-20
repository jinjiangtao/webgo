package utils

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
)

func Tokenize(text string) []string {
	if text == "" {
		return nil
	}

	text = strings.TrimSpace(text)
	var tokens []string
	seen := make(map[string]bool)

	re := regexp.MustCompile(`[\p{Han}]+|[a-zA-Z0-9]+`)
	matches := re.FindAllString(text, -1)

	for _, match := range matches {
		if len(match) == 0 {
			continue
		}

		if isChinese(match) {
			chineseTokens := tokenizeChinese(match)
			for _, t := range chineseTokens {
				if !seen[t] && len(t) > 0 {
					seen[t] = true
					tokens = append(tokens, t)
				}
			}
		} else {
			lower := strings.ToLower(match)
			if !seen[lower] {
				seen[lower] = true
				tokens = append(tokens, lower)
			}
			if len(lower) >= 3 {
				subTokens := generateSubstrings(lower)
				for _, st := range subTokens {
					if !seen[st] {
						seen[st] = true
						tokens = append(tokens, st)
					}
				}
			}
		}
	}

	tokens = append(tokens, strings.ToLower(text))

	return tokens
}

func isChinese(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

func tokenizeChinese(text string) []string {
	var tokens []string
	runes := []rune(text)
	n := len(runes)

	for length := 2; length <= min(n, 4); length++ {
		for i := 0; i <= n-length; i++ {
			token := string(runes[i : i+length])
			tokens = append(tokens, token)
		}
	}

	for i := 0; i < n; i++ {
		tokens = append(tokens, string(runes[i]))
	}

	return tokens
}

func generateSubstrings(s string) []string {
	var subs []string
	n := len(s)
	for length := 2; length < n; length++ {
		for i := 0; i <= n-length; i++ {
			subs = append(subs, s[i:i+length])
		}
	}
	return subs
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Highlight(text string, keywords []string, highlightTag string) string {
	if len(keywords) == 0 || text == "" {
		return text
	}

	if highlightTag == "" {
		highlightTag = "mark"
	}

	type matchPos struct {
		start int
		end   int
	}

	var matches []matchPos
	runes := []rune(text)
	lowerText := strings.ToLower(text)

	for _, kw := range keywords {
		kw = strings.TrimSpace(kw)
		if kw == "" {
			continue
		}
		kwLower := strings.ToLower(kw)
		kwRunes := []rune(kw)
		kwLen := len(kwRunes)

		start := 0
		for {
			idx := strings.Index(lowerText[start:], kwLower)
			if idx == -1 {
				break
			}

			actualStart := runeIndex(lowerText, start+idx)
			actualEnd := actualStart + kwLen

			overlap := false
			for _, m := range matches {
				if actualStart < m.end && actualEnd > m.start {
					overlap = true
					if actualStart < m.start {
						m.start = actualStart
					}
					if actualEnd > m.end {
						m.end = actualEnd
					}
					break
				}
			}

			if !overlap {
				matches = append(matches, matchPos{start: actualStart, end: actualEnd})
			}

			byteStart := start + idx + len(kwLower)
			if byteStart >= len(lowerText) {
				break
			}
			start = byteStart
		}
	}

	if len(matches) == 0 {
		return text
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].start < matches[j].start
	})

	var merged []matchPos
	for _, m := range matches {
		if len(merged) == 0 {
			merged = append(merged, m)
			continue
		}
		last := &merged[len(merged)-1]
		if m.start <= last.end {
			if m.end > last.end {
				last.end = m.end
			}
		} else {
			merged = append(merged, m)
		}
	}

	var result strings.Builder
	lastEnd := 0
	openTag := "<" + highlightTag + " class='search-highlight'>"
	closeTag := "</" + highlightTag + ">"

	for _, m := range merged {
		if m.start > lastEnd {
			result.WriteString(string(runes[lastEnd:m.start]))
		}
		result.WriteString(openTag)
		result.WriteString(string(runes[m.start:m.end]))
		result.WriteString(closeTag)
		lastEnd = m.end
	}

	if lastEnd < len(runes) {
		result.WriteString(string(runes[lastEnd:]))
	}

	return result.String()
}

func runeIndex(s string, byteIndex int) int {
	return len([]rune(s[:byteIndex]))
}

func CalculateScore(title, content, tags string, tokens []string, searchText string) float64 {
	var score float64
	searchText = strings.ToLower(strings.TrimSpace(searchText))
	lowerTitle := strings.ToLower(title)
	lowerContent := strings.ToLower(content)
	lowerTags := strings.ToLower(tags)

	if searchText != "" {
		if lowerTitle == searchText {
			score += 50
		}
		if strings.HasPrefix(lowerTitle, searchText) {
			score += 30
		}
		if strings.Contains(lowerTitle, searchText) {
			score += 20
		}
		if strings.Contains(lowerTags, searchText) {
			score += 15
		}
		if strings.Contains(lowerContent, searchText) {
			score += 10
		}
	}

	for _, token := range tokens {
		token = strings.ToLower(strings.TrimSpace(token))
		if token == "" {
			continue
		}

		tokenLen := len([]rune(token))
		weight := 1.0
		if tokenLen >= 3 {
			weight = 1.5
		}
		if tokenLen >= 5 {
			weight = 2.0
		}

		if strings.Contains(lowerTitle, token) {
			score += 3.0 * weight
		}
		if strings.Contains(lowerTags, token) {
			score += 2.0 * weight
		}
		if strings.Contains(lowerContent, token) {
			score += 1.0 * weight
		}
	}

	return score
}

func GenerateSuggestions(keywords []string, input string, limit int) []string {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}

	inputLower := strings.ToLower(input)
	type scoredItem struct {
		text  string
		score float64
	}

	var scored []scoredItem
	seen := make(map[string]bool)

	for _, kw := range keywords {
		kwLower := strings.ToLower(kw)
		var s float64

		if kwLower == inputLower {
			s += 100
		}
		if strings.HasPrefix(kwLower, inputLower) {
			s += 50
		}
		if strings.Contains(kwLower, inputLower) {
			s += 20
		}

		kwRunes := []rune(kw)
		inputRunes := []rune(input)
		if len(inputRunes) >= 2 && len(kwRunes) >= len(inputRunes) {
			matchCount := 0
			for _, ir := range inputRunes {
				for _, kr := range kwRunes {
					if ir == kr {
						matchCount++
						break
					}
				}
			}
			s += float64(matchCount) * 2
		}

		if s > 0 && !seen[kw] {
			seen[kw] = true
			scored = append(scored, scoredItem{text: kw, score: s})
		}
	}

	sort.Slice(scored, func(i, j int) bool {
		if scored[i].score != scored[j].score {
			return scored[i].score > scored[j].score
		}
		return len([]rune(scored[i].text)) < len([]rune(scored[j].text))
	})

	if limit <= 0 || limit > len(scored) {
		limit = len(scored)
	}

	result := make([]string, limit)
	for i := 0; i < limit; i++ {
		result[i] = scored[i].text
	}

	return result
}
