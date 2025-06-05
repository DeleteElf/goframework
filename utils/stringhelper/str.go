package stringhelper

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"regexp"
	"strings"
)

// ConvertCamelToSnake 驼峰命名法转回下划线命名
func ConvertCamelToSnake(camel string, joinStr string) string {
	if len(joinStr) == 0 {
		joinStr = "_"
	}
	replaceStr := `${1}` + joinStr + `${2}`
	str := regexp.MustCompile(`([a-z])([A-Z])`).ReplaceAllString(camel, replaceStr)
	return strings.ToLower(str)
}

// ConvertCamelToSnakeWithDefault 驼峰命名法转回下划线命名
func ConvertCamelToSnakeWithDefault(camel string) string {
	return ConvertCamelToSnake(camel, "_")
}

// ConvertToCamel 字符串转成驼峰命名法，下划线后首字母大写，首个前缀只有1个字符则丢弃
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

// FormatStringByObject 格式化字符串，使用对象作为参数封装
func FormatStringByObject(format string, data interface{}) string {
	//时间格式化字符串 2006-01-02 15:04:05.000 Mon MST
	return fmt.Sprintf(format, data)
}

// FormatString 格式化字符串，使用参数列表封装
func FormatString(format string, args ...any) string {
	//时间格式化字符串 2006-01-02 15:04:05.000 Mon MST
	return fmt.Sprintf(format, args)
}
