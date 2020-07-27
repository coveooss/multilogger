package reutils

import "regexp"

// ProgressiveReplaceAll is similar to regexp.ReplaceAllStringFunc, but it allows the user to backtrack
// the searched string instead of simply going from match to match.
func ProgressiveReplaceAll(re *regexp.Regexp, source string, repl func(string, *int) string) (result string) {
	for {
		match := re.FindStringIndex(source)
		if match == nil {
			// There is no more match, we append the rest of the source and return
			result += source
			break
		}
		start, end := match[0], match[1]

		// We append the non matching part
		result += source[:start]
		pos := end - start
		result += repl(source[start:end], &pos)
		source = source[start+pos:]
	}
	return
}
