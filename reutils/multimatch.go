package reutils

import (
	"regexp"
	"sync"
)

// MultiMatch returns a map of matching elements from a list of regular expressions (returning the first matching element)
func MultiMatch(s string, expressions ...*regexp.Regexp) (map[string]string, int) {
	for i, re := range expressions {
		if matches := re.FindStringSubmatch(s); len(matches) != 0 {
			results := make(map[string]string, len(matches))
			results[""] = matches[0]
			for i, key := range re.SubexpNames() {
				if key != "" {
					results[key] = matches[i]
				}
			}
			return results, i
		}
	}
	return nil, -1
}

// NewRegexGroup cache compiled regex to avoid multiple interpretation of the same regex
func NewRegexGroup(key string, definitions ...string) (result []*regexp.Regexp, err error) {
	result = make([]*regexp.Regexp, len(definitions))
	for i := range definitions {
		regex, err := regexp.Compile(definitions[i])
		if err != nil {
			return nil, err
		}
		result[i] = regex
	}

	cachedRegex.Store(key, result)
	return
}

// GetRegexGroup tries to retreive a regular expression group previously created by NewRegexGroup.
func GetRegexGroup(key string) (result []*regexp.Regexp) {
	if cached, ok := cachedRegex.Load(key); ok {
		result = cached.([]*regexp.Regexp)
	}
	return
}

// DeleteRegexGroup tries to retreive a regular expression group previously created by NewRegexGroup.
func DeleteRegexGroup(key string) (result []*regexp.Regexp) {
	result = GetRegexGroup(key)
	cachedRegex.Delete(key)
	return
}

var cachedRegex sync.Map
