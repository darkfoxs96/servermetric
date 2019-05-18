package tools

import "strings"

func IncludeStr(arr []string, search string) bool {
	for _, v := range arr {
		if v == search {
			return true
		}
	}

	return false
}

func ContainsOnly(text string, only ...string) bool {
	for _, item := range only {
		if strings.Contains(text, item) {
			return true
		}
	}

	return false
}
