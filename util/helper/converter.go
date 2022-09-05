package util

import (
	"strconv"
	"strings"
)

func GenerateIDs(baseValue string, count uint64) []string {
	ids := make([]string, count)
	for i := uint64(0); i < count; i++ {
		if v, err := strconv.ParseUint(baseValue, 10, 64); err == nil {
			ids[i] = strconv.FormatUint(v+i, 10)
		}
	}
	return ids
}

func ValueExtractor(data string, token string) string {
	indexToken := strings.Index(data, token)
	if indexToken < 0 {
		return ""
	}
	indexToken += len(token)
	indexEnd := strings.Index(data[indexToken:], "&")
	if indexToken < 0 {
		indexEnd = strings.Index(data[indexToken:], "*")
	}
	indexEnd += indexToken

	if indexToken < 0 || indexEnd < 0 {
		return ""
	}

	if indexEnd < indexToken {
		return ""
	}

	return data[indexToken:indexEnd]
}
