package stammer

import "strings"

func FilterSkipList(words []string, skipList map[string]bool) []string {
	ret := make([]string, 0, len(words))
	for _, word := range words {
		if skipList[strings.ToLower(word)] {
			continue
		}

		ret = append(ret, word)
	}

	return ret
}
