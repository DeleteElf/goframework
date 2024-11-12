package stringhelper

import (
	"regexp"
	"strings"
)

func ConvertCamelToSnake(camel string, joinStr string) string {
	var snake []string
	for _, word := range regexp.MustCompile(`([a-z])([A-Z])`).ReplaceAllString(camel, `$1_$2`) {
		snake = append(snake, strings.ToLower(string(word)))
	}
	if len(joinStr) == 0 {
		joinStr = "_"
	}
	return strings.Join(snake, joinStr)
}

func ConvertCamelToSnakeWithDefault(camel string) string {
	return ConvertCamelToSnake(camel, "_")
}
