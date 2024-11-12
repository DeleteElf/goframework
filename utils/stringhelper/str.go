package stringhelper

import (
	"regexp"
	"strings"
)

func ConvertCamelToSnake(camel string, joinStr string) string {
	if len(joinStr) == 0 {
		joinStr = "_"
	}
	replaceStr := `${1}` + joinStr + `${2}`
	str := regexp.MustCompile(`([a-z])([A-Z])`).ReplaceAllString(camel, replaceStr)
	return strings.ToLower(str)
}

func ConvertCamelToSnakeWithDefault(camel string) string {
	return ConvertCamelToSnake(camel, "_")
}
