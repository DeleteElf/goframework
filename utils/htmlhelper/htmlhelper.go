package htmlhelper

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
	"io"
	"os"
)

// ReadAll 读取文件的所有内容
func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(f)
}

// MarkdownToHtml 将Markdown格式的内容转化成html格式的内容
func MarkdownToHtml(md string) string {
	return string(MarkdownToHtmlByte([]byte(md)))
}

// MarkdownToHtmlByte 将Markdown格式的内容转化成html格式的内容
func MarkdownToHtmlByte(md []byte) []byte {
	bytes := blackfriday.Run(md)
	return bluemonday.UGCPolicy().SanitizeBytes(bytes)
}
