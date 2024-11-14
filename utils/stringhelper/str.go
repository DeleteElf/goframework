package stringhelper

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func ConvertToCamel(src string) string {
	words := strings.Split(src, "_")
	var result string
	caser := cases.Title(language.Und, cases.NoLower)
	var hasFirst bool = false
	for i, word := range words {
		if i == 0 && len(word) == 1 { //首个前缀只有1个字符，则不进行转换
			continue
		}
		if len(word) == 0 { //如果出现连续__导致的无效长度字符，则不进行转换
			continue
		}
		if !hasFirst {
			result += strings.ToLower(word) //首单词不进行大写的支持
			hasFirst = true
		} else {
			result += caser.String(word)
		}
	}
	return result
}
