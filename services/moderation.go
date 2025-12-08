package services

import "strings"

func IsFilenameStudyRelated(filename string, keywords []string) bool {
	filename = strings.ToLower(filename)

	for _, kw := range keywords {
		kw = strings.ToLower(kw)

		if kw != "" && strings.Contains(filename, kw) {
			return true
		}
	}

	return false
}
